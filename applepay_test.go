package qiwi

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBase64Decode(t *testing.T) {
	enc := "YWJjMTIzIT8kKiYoKSctPUB+"
	expected := "abc123!?$*&()'-=@~"

	res, err := decodeBase64(enc)

	if string(res) != expected {
		t.Errorf("Wrong function return: %s, expected: %s", enc, expected)
	}

	if err != nil {
		t.Errorf("Error occurred: %s", err)
	}
}

func TestApplePay(t *testing.T) {

	// "paymentMethod": {
	// 	  "type": "APPLE_PAY_TOKEN",
	appleTokenStructure := []byte(`{
	   "paymentData":{
		  "version":"EC_v1",
		  "data":"IaD7LKDbJsOrGTlNGkKUC95Y+4an2YuN0swstaCaoovlj8dbgf16FmO5j4AX80L0xsRQYKLUpgUHbGoYF26PbraIdZUDtPtja4HdqDOXGESQGsAKCcRIyujAJpFv95+5xkNldDKK2WTe28lHUDTA9bykIgrvasYaN9VWnS92i2CZPpsI7yu13Kk3PrUceuA3Fb6wFgJ0l7HXL1RGhrA7V5JKReo/EwikMsK8AfChK7pvWaB51SsMvbMJF28JnincfVX39vYHdzEwpjSPngNiszGqZGeLdqNE3ngkoEK1AW2ymbYkIoy9KFdXayekELR6hQWnL4MCutLesLjKhyTN26fxBamPHzAf/IczAdWBDq2P/59jheIGrnK30slJJcr1Bocb8rqojyaVZIY+Xk24Nc6dvSdJhfDDyhX56pn5YtWOxWuVOT0tZSJvxBN/HeIuYcNG6R9u7CHpcelsi4I8O+1gruKKZQHweERG2DyCmoUO9zlajOSm",
		  "header":{
			 "ephemeralPublicKey":"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEzLx7FJhw1Z1PmxOLMTQBs1LgKEUT6 hcxLlT8wGrzwyY8tKeG+VPSjryVkTFYECrj+5r28PJWtmvn/unYhbYDaQ==",
			 "publicKeyHash":"OrWgjRGkqEWjdkRdUrXfiLGD0he/zpEu512FJWrGYFo=",
			 "transactionId":"1234567890ABCDEF"
		  },
		  "signature":"ZmFrZSBzaWduYXR1cmU="
	   }
	}`)
	// }

	payload := base64.StdEncoding.EncodeToString(appleTokenStructure)

	// HTTP MOCK
	serv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "PUT" {
			fmt.Fprintln(w, `{
				  "serviceName" : "payin-core",
				  "errorCode" : "validation.wrongmethod",
				  "description" : "bad http method",
				  "userMessage" : "Validation error",
				  "dateTime" : "2018-11-13T16:49:59.166+03:00",
				  "traceId" : "fd0e2a08c63ace83"
				}`)
			return
		}

		fmt.Fprintln(w, "{}")
	}))
	serv.Start()
	defer serv.Close()

	// Route request to mocked http server
	pay := New("billId", "SiteID", "TOKEN", serv.URL)
	amount := 500 // 5.00RUB
	err := pay.ApplePay(context.TODO(), amount, payload)
	if err != nil {
		t.Errorf("ApplePay method error: %s", err)
	}

}

func TestApplePayBadToken(t *testing.T) {

	appleBADToken := []byte(`{
	  {{bad json
	}`)

	payload := base64.StdEncoding.EncodeToString(appleBADToken)

	// HTTP MOCK
	serv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		var pl Payment
		var err error

		_, err = io.Copy(&buf, r.Body)
		if err != nil {
			err = fmt.Errorf("[QIWI] %w: %s", ErrBadJSON, err)
			fmt.Fprintln(w, `{
				  "serviceName" : "payin-core",
				  "errorCode" : "validation.copyerr",
				  "description" : "`+err.Error()+`",
				  "userMessage" : "Validation error",
				  "dateTime" : "2018-11-13T16:49:59.166+03:00",
				  "traceId" : "fd0e2a08c63ace83"
				}`)
			return
		}

		err = json.Unmarshal(buf.Bytes(), &pl)
		if err != nil {
			err = fmt.Errorf("[QIWI] %w: %s", ErrBadJSON, err)
			fmt.Fprintln(w, `{
				  "serviceName" : "payin-core",
				  "errorCode" : "validation.json",
				  "description" : "`+err.Error()+`",
				  "userMessage" : "Validation error",
				  "dateTime" : "2018-11-13T16:49:59.166+03:00",
				  "traceId" : "fd0e2a08c63ace83"
				}`)
			return
		}

		fmt.Fprintln(w, "{}")
	}))
	serv.Start()
	defer serv.Close()

	// Route request to mocked http server
	pay := New("billId", "SiteID", "TOKEN", serv.URL)
	amount := 500 // 5.00RUB
	err := pay.ApplePay(context.TODO(), amount, payload)
	if !errors.Is(err, ErrBadJSON) {
		t.Errorf("ApplePay bad token json wrong error: %s", err)
	}

}
