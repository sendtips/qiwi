package qiwi

import (
	"errors"
	"testing"
)

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

	tests := []struct {
		payload []byte
		want    *Notify
		err     error
		sig     string
	}{
		{
			payload,
			&Notify{Type: PaymentNotify, Payment: Payment{Amount: NewAmountInRubles(221124)}},
			nil, "426917662ee15d568a5cddc14620cee02c604364185ac3f3221ff33d1d2fa49f"},

		{
			[]byte(`{{{bad json}`),
			&Notify{},
			ErrBadJSON, ""},
	}

	for _, test := range tests {

		notify, err := NewNotify(key, test.sig, test.payload)

		if !errors.Is(err, test.err) {
			t.Error("Error occurred: ", err, test.err)
		}

		if notify.Type != test.want.Type {
			t.Error("Incorrect type")
		}

		if notify.Payment.Amount.Value != test.want.Payment.Amount.Value {
			t.Error("Amount is wrong")
		}

	}

}
