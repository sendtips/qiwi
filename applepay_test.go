package qiwi

import (
	"context"
	"encoding/base64"
	"fmt"
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
	apple_token_structure := []byte(`{
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

	payload := base64.StdEncoding.EncodeToString(apple_token_structure)

	// HTTP MOCK
	serv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
