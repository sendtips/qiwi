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
		{"TOKEN", "426917662ee15d568a5cddc14620cee02c604364185ac3f3221ff33d1d2fa49f", payloadPayment},
		{"TOKEN", "0b0b2140e0614ded05feefdb448fb77d29c4616591f036f2a7491e7799935eda", payloadCapture},
		{"TOKEN", "4c445af0c065bfc531866cd21f2b93d67f62820cef74341efed506da717a4424", payloadRefund},
		{"TOKEN", "747d49a7bb98c75249659822e3d43912225ac9a75bdddf7bc967af7439d87cfd", payloadCheckCard},
		{"TOKEN", "0a7d5b657b61bfac18f350d57408b79fca88eb0fd6db644dc9e4c436c3f0d056", payloadBADType},
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
