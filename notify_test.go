package qiwi

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotifyType(t *testing.T) {
	tests := []struct {
		input string
		want  NotifyType
	}{
		{"PAYMENT", PaymentNotify},
		{"CAPTURE", CaptureNotify},
		{"REFUND", RefundNotify},
		{"CHECK_CARD", CheckCardNotify},
	}

	for _, test := range tests {
		if NotifyType(test.input) != test.want {
			t.Errorf("Bad NotifyType %s", test.input)
		}

		if test.input != string(test.want) {
			t.Errorf("Bad string method result for %s", test.input)
		}
	}
}

func TestHook(t *testing.T) {
	const key = "TOKEN"

	payload := []byte(`
        {
           "payment":{
              "paymentid":"4504751",
              "tokendata":{
                 "paymenttoken":"4cc975be-483f-8d29-2b7de3e60c2f",
                 "expireddate":"2021-12-31T00:00:00+03:00"
              },
              "type":"payment",
              "createddatetime":"2019-10-08T11:31:37+03:00",
              "status":{
                 "value":"success",
                 "changeddatetime":"2019-10-08T11:31:37+03:00"
              },
              "amount":{
                 "value":2211.24,
                 "currency":"RUB"
              },
              "paymentMethod":{
                 "type":"CARD",
                 "maskedPan":"220024/*/*/*/*/*/*5036",
                 "rrn":null,
                 "authCode":null,
                 "type":"CARD"
              },
              "paymentCardInfo": {
                 "issuingCountry": "810",
                 "issuingBank": "QiwiBank",
                 "paymentSystem": "VISA",
                 "fundingSource": "CREDIT",
                 "paymentSystemProduct": "P|Visa Gold"
              },
              "customer":{
                 "ip":"79.142.20.248",
                 "account":"token32",
                 "phone":"0"
              },
              "billId":"testing122",
              "customFields":{},
              "flags":[
                 "SALE"
              ]
           },
           "type":"PAYMENT",
           "version":"1"
        }
        `)

	payloadSplits := []byte(`{
		"refund": {
			"refundId": "42f5ca91-965e-4cd0-bb30-3b64d9284048",
			"type": "REFUND",
			"createdDateTime": "2021-02-05T11:31:40+03:00",
			"status": {
				"value": "SUCCESS",
				"changedDateTime": "2021-02-05T11:31:40+03:00"
			},
			"amount": {
				"value": 3,
				"currency": "RUB"
			},
			"paymentMethod": {
				"type": "TOKEN",
				"paymentToken": "1620161e-d498-431b-b006-c52bb78c6bf2",
				"maskedPan": null,
				"cardHolder": null,
				"cardExpireDate": null
			},
			"customer": {
				"email": "glmgmmxr@qiwi123.com",
				"account": "sbderxuftsrt",
				"phone": "13387571067",
				"country": "yccsnnfjgthu",
				"city": "sqdvseezbpzo",
				"region": "shbvyjgspjvu"
			},
			"gatewayData": {
				"type": "ACQUIRING",
				"authCode": "181218",
				"rrn": "123"
			},
			"billId": "autogenerated-19cf2596-62a8-47f2-8721-b8791e9598d0",
			"flags": [
				"REVERSAL"
			],
			"refundSplits": [
				{
					"type": "MERCHANT_DETAILS",
					"siteUid": "Obuc-00",
					"splitAmount": {
						"value": 2,
						"currency": "RUB"
					},
					"splitCommissions": {
						"merchantCms": {
							"value": 0,
							"currency": "RUB"
						},
						"userCms": null
					},
					"orderId": "dressesforwhite",
					"comment": "Покупка 1"
				},
				{
					"type": "MERCHANT_DETAILS",
					"siteUid": "Obuc-01",
					"splitAmount": {
						"value": 1,
						"currency": "RUB"
					},
					"splitCommissions": {
						"merchantCms": {
							"value": 0.02,
							"currency": "RUB"
						},
						"userCms": null
					},
					"orderId": "shoesforvalya",
					"comment": "Покупка 2"
				}
			]
		},
		"type": "REFUND",
		"version": "1"
	}`)

	tests := []struct {
		payload []byte
		want    Notify
		err     error
		sig     string
	}{
		{
			payload,
			Notify{Type: PaymentNotify, Payment: Payment{Amount: NewAmountInRubles(221124)}},
			nil, "3c67f9a691e34e1a9e74e05927f3901186cc838cc81de2a3519c78b9612cf49e",
		},

		{
			[]byte(`{{{bad json}`),
			Notify{},
			ErrBadJSON, "",
		},

		{
			payload,
			Notify{Type: PaymentNotify, Payment: Payment{Amount: NewAmountInRubles(221124)}},
			ErrBadSignature, "BADSIGN",
		},

		{
			payloadSplits,
			Notify{Type: RefundNotify, Refund: Payment{Amount: NewAmountInRubles(300)}},
			nil, "681d7b40abafb5dcf958f653887303bd081839fa963bad6149055e90a2b02b16",
		},
	}

	for _, test := range tests {

		notify, err := NewNotify(key, test.sig, test.payload)

		if !errors.Is(err, test.err) {
			t.Error("Error occurred: ", err, test.err, notify)
		}

		if notify.Type != test.want.Type {
			t.Error("Incorrect type")
		}

		if notify.Payment.Amount.Value != test.want.Payment.Amount.Value {
			t.Error("Amount is wrong", notify.Payment.Amount.Value, test.want.Payment.Amount.Value)
		}

		if notify.Type == RefundNotify {
			if notify.Refund.RefundSplits[1].Commission.Amount.Value != 0.02 {
				t.Error("Split refuns amount is wrong")
			}
		}

	}
}

func TestHttpRequestHook(t *testing.T) {
	var payload bytes.Buffer

	_, _ = payload.WriteString(`{
                   "payment":{
                      "paymentid":"4504751",
                      "tokendata":{
                         "paymenttoken":"4cc975be-483f-8d29-2b7de3e60c2f",
                         "expireddate":"2021-12-31T00:00:00+03:00"
                      },
                      "type":"payment",
                      "createddatetime":"2019-10-08T11:31:37+03:00",
                      "status":{
                         "value":"success",
                         "changeddatetime":"2019-10-08T11:31:37+03:00"
                      },
                      "amount":{
                         "value":2211.24,
                         "currency":"RUB"
                      },
                      "paymentMethod":{
                         "type":"CARD",
                         "maskedPan":"220024/*/*/*/*/*/*5036",
                         "rrn":null,
                         "authCode":null,
                         "type":"CARD"
                      },
                      "paymentCardInfo": {
                         "issuingCountry": "810",
                         "issuingBank": "QiwiBank",
                         "paymentSystem": "VISA",
                         "fundingSource": "CREDIT",
                         "paymentSystemProduct": "P|Visa Gold"
                      },
                      "customer":{
                         "ip":"79.142.20.248",
                         "account":"token32",
                         "phone":"0"
                      },
                      "billId":"testing122",
                      "customFields":{},
                      "flags":[
                         "SALE"
                      ]
                   },
                   "type":"PAYMENT",
                   "version":"1"
                }`)

	tests := []struct {
		token, signature string
		payload          bytes.Buffer
		err              error
	}{
		{"TOKEN", "3c67f9a691e34e1a9e74e05927f3901186cc838cc81de2a3519c78b9612cf49e", payload, nil},
		{"TOKEN", "3c67f9a691e34e1a9e74e05927f3901186cc838cc81de2a3519c78b9612cf49e", genBigBody(), ErrBadJSON},
		{"TOKEN", "", payload, ErrBadSignature},
	}

	for _, test := range tests {
		handler := func(w http.ResponseWriter, r *http.Request) {
			r.Header.Add("Content-Type", "application/json")
			r.Header.Add("Signature", test.signature)

			notify, err := NotifyParseHTTPRequest(test.token, w, r)

			if !errors.Is(err, test.err) {
				t.Errorf("Error: %s", err)
			}

			if test.err == nil {
				if notify.Type != PaymentNotify {
					t.Error("Incorrect type")
				}

				if notify.Payment.Amount.Currency != rub {
					t.Error("Currency is wrong")
				}
			}

			fmt.Fprint(w, "{}")
		}

		req := httptest.NewRequest("POST", "/qiwinotify", &test.payload)
		rec := httptest.NewRecorder()

		handler(rec, req)
	}
}

// genBigBody used to emulate endless body request attack.
func genBigBody() bytes.Buffer {
	var buf bytes.Buffer

	tooBigLength := maxBodyBytes + int64(1)
	buf.WriteString(`{"payment": "`)
	for {
		if len(buf.Bytes()) >= int(tooBigLength) {
			buf.WriteString(`"}`)
			return buf
		}
		num, err := rand.Int(rand.Reader, big.NewInt(int64(10)))
		if err != nil {
			return buf
		}
		n := num.Int64()
		// Put only ints
		if n > 0 && n < 9 {
			buf.WriteString(fmt.Sprint(n))
		}
	}
}
