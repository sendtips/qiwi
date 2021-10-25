package qiwi

import (
	"testing"
	"encoding/json"
)

func TestSignature(t *testing.T) {

	// headers := `
	// POST /qiwi-notify.php HTTP/1.1
	// Accept: application/json
	// Content-Type: application/json
	// Signature: j4wnfnzd***v5mv2w=
	// Host: server.ru
	// `

	var notify Notify

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

	const key string = "TOKEN"
	const hash string = "426917662ee15d568a5cddc14620cee02c604364185ac3f3221ff33d1d2fa49f"
	sig := NewSignature(key, hash)

	err := json.Unmarshal(payload, &notify)
	if err != nil {
		t.Errorf("Signature payload JSON unmarshaling error: %s", err)
	}

	if !sig.Verify(notify) {
		t.Errorf("Wrong signature")
	}

}
