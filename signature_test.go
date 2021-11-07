package qiwi

import (
	"encoding/json"
	"testing"
)

func TestSignature(t *testing.T) {
	// TEST PAYLOADS
	payloadPayment := []byte(`
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

	payloadCapture := []byte(`
        {
           "capture":{
              "captureId":"4504758",
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
           "type":"CAPTURE",
           "version":"1"
        }
        `)

	payloadRefund := []byte(`
    {
       "refund":{
          "refundId":"4504759",
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
       "type":"REFUND",
       "version":"1"
    }
    `)

	payloadCheckCard := []byte(`
    {
       "checkPaymentMethod":{
          "requestUid":"4504751",
          "isValidCard": true,
          "checkOperationDate":"2019-10-08T11:31:37+03:00",
          "status":{
             "value":"success",
             "changeddatetime":"2019-10-08T11:31:37+03:00"
          }
       },
       "type":"CHECK_CARD",
       "version":"1"
    }
    `)

	payloadBADType := []byte(`
    {
       "payment":{
          "paymentid":"4504750",
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
       "type":"UNKNOWN_TYPE",
       "version":"1"
    }
    `)

	// TESTS ITSELF
	tests := []struct {
		key, hash string
		payload   []byte
	}{
		{"TOKEN", "3c67f9a691e34e1a9e74e05927f3901186cc838cc81de2a3519c78b9612cf49e", payloadPayment},
		{"TOKEN", "5167b2bd0b07957e11686552ff40e3df688478f660c871f4257becdebd3acca5", payloadCapture},
		{"TOKEN", "f3db008ae056f18d3066f0cc8b3f1725af70da8e8a0924759f57d9db6412d659", payloadRefund},
		{"TOKEN", "43b3f5544546c5a3dcdb0a5ae35d60689578665394a0f6cbac01c78e6876f03f", payloadCheckCard},
		{"TOKEN", "5235526ab9b4fc6dccf5762efd79359470bf4dff673fc2d2f34c907e6d41c0c5", payloadBADType},
	}

	for _, test := range tests {
		var notify Notify

		sig := NewSignature(test.key, test.hash)

		err := json.Unmarshal(test.payload, &notify)
		if err != nil {
			t.Errorf("Signature payload JSON unmarshaling error: %s", err)
		}

		if !sig.Verify(&notify) {
			t.Errorf("Wrong signature")
		}
	}
}

func TestBadSignature(t *testing.T) {
	var notify Notify

	payload := []byte(`
        {
           "checkPaymentMethod":{
              "requestUid":"4504751",
              "isValidCard": true,
              "checkOperationDate":"2019-10-08T11:31:37+03:00",
              "status":{
                 "value":"success",
                 "changeddatetime":"2019-10-08T11:31:37+03:00"
              }
           },
           "type":"CHECK_CARD",
           "version":"1"
        }
    `)

	sig := NewSignature("TOKEN", "BADSIGN")

	err := json.Unmarshal(payload, &notify)
	if err != nil {
		t.Errorf("Signature payload JSON unmarshaling error: %s", err)
	}

	if sig.Verify(&notify) {
		t.Errorf("Bad signature shoud not pass")
	}
}
