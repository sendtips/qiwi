package qiwi

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Examplezlibzompress() {
	// eJzzSM3JyVcozy/KSQEAGKsEPQ== -- compressed and base64 encoded

	// if string(res) != string(result) {
	// 	t.Errorf("Wrong result of compress function %+q", res)
	// }

	data := []byte("Hello world")
	res := zlibzompress(data)
	fmt.Println([]byte(res))
	// Output: [120 156 242 72 205 201 201 87 40 207 47 202 73 1 4 0 0 255 255 24 171 4 61]
}

func TestGooglePay(t *testing.T) {

	// google_token_structure := []byte(`
	// "paymentMethod": {
	// 	"type": "GOOGLE_PAY_TOKEN",
	// 	"paymentToken": "eJxVUtuK2zAQfd+vCHGShS9mS0hb8YChjabx"
	// 	}
	// `)
	type GooglePayToken struct {
		ProtoVer string `json:"protocolVersion"`
	}

	googlePayToken := []byte(`{
		  "protocolVersion":"ECv2",
		  "signature":"MEQCIH6Q4OwQ0jAceFEkGF0JID6sJNXxOEi4r+mA7biRxqBQAiAondqoUpU/bdsrAOpZIsrHQS9nwiiNwOrr24RyPeHA0Q\u003d\u003d",
		  "intermediateSigningKey":{
			"signedKey": "{\"keyExpiration\":\"1542323393147\",\"keyValue\":\"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE/1+3HBVSbdv+j7NaArdgMyoSAM43yRydzqdg1TxodSzA96Dj4Mc1EiKroxxunavVIvdxGnJeFViTzFvzFRxyCw\\u003d\\u003d\"}",
			"signatures": ["MEYCIQCO2EIi48s8VTH+ilMEpoXLFfkxAwHjfPSCVED/QDSHmQIhALLJmrUlNAY8hDQRV/y1iKZGsWpeNmIP+z+tCQHQxP0v"]
		  },
		  "signedMessage":"{\"tag\":\"jpGz1F1Bcoi/fCNxI9n7Qrsw7i7KHrGtTf3NrRclt+U\\u003d\",\"ephemeralPublicKey\":\"BJatyFvFPPD21l8/uLP46Ta1hsKHndf8Z+tAgk+DEPQgYTkhHy19cF3h/bXs0tWTmZtnNm+vlVrKbRU9K8+7cZs\\u003d\",\"encryptedMessage\":\"mKOoXwi8OavZ\"}"
		}`)

	// HTTP MOCK
	serv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		var buf, b bytes.Buffer
		var p Payment
		var tk GooglePayToken

		io.Copy(&buf, r.Body)

		err = json.Unmarshal(buf.Bytes(), &p)
		if err != nil {
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

		//
		// enc := base64.NewDecoder(base64.StdEncoding, )
		// enc.Read(buf.Bytes())
		dec, err := base64.StdEncoding.DecodeString(p.PaymentMethod.GooglePaymentToken)
		if err != nil {
			fmt.Fprintln(w, `{
		  "serviceName" : "payin-core",
		  "errorCode" : "validation.base64",
		  "description" : "`+err.Error()+`",
		  "userMessage" : "Validation error",
		  "dateTime" : "2018-11-13T16:49:59.166+03:00",
		  "traceId" : "fd0e2a08c63ace83"
		}`)
			return
		}

		br := bytes.NewReader(dec)
		data, err := zlib.NewReader(br)
		if err != nil {
			fmt.Fprintln(w, `{
				  "serviceName" : "payin-core",
				  "errorCode" : "validation.zlib",
				  "description" : "`+err.Error()+`",
				  "userMessage" : "Validation error",
				  "dateTime" : "2018-11-13T16:49:59.166+03:00",
				  "traceId" : "fd0e2a08c63ace83"
				}`)
			return
		}
		defer data.Close()
		io.Copy(&b, data)

		err = json.Unmarshal(b.Bytes(), &tk)
		if err != nil {
			fmt.Fprintln(w, `{
				  "serviceName" : "payin-core",
				  "errorCode" : "validation.tokenjson",
				  "description" : "`+err.Error()+`",
				  "userMessage" : "Validation error",
				  "dateTime" : "2018-11-13T16:49:59.166+03:00",
				  "traceId" : "fd0e2a08c63ace83"
				}`)
			return
		}

		if tk.ProtoVer != "ECv2" {
			fmt.Fprintln(w, `{
				  "serviceName" : "payin-core",
				  "errorCode" : "validation.tokenversionpayload",
				  "description" : "`+err.Error()+`",
				  "userMessage" : "Validation error",
				  "dateTime" : "2018-11-13T16:49:59.166+03:00",
				  "traceId" : "fd0e2a08c63ace83"
				}`)
			return
		}

		fmt.Fprintln(w, `{
			"siteId": "23044",
			"billId": "893794793973",
			"amount": {
			  "value": 100.00,
			  "currency": "RUB"
			},
			"status": {
			  "value": "CREATED",
			  "changedDateTime": "2018-03-05T11:27:41"
			},
			"comment": "Text comment",
			"customFields": {
			  "cf1": "Some data",
			  "FROM_MERCHANT_CONTRACT_ID": "1234"
			},
			"creationDateTime": "2018-03-05T11:27:41",
			"expirationDateTime": "2018-04-13T14:30:00",
			"payUrl": "https://oplata.qiwi.com/form/?invoice_uid=d875277b-6f0f-445d-8a83-f62c7c07be77"
		}`)
	}))
	serv.Start()
	defer serv.Close()

	// Route request to mocked http server
	pay := New("billId", "SiteID", "TOKEN", serv.URL)
	amount := 500 // 5.00RUB
	err := pay.GooglePay(context.TODO(), amount, googlePayToken)
	if err != nil {
		t.Errorf("GooglePay method error: %s", err)
	}

}
