package qiwi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func OffTestSplits(t *testing.T) {
	payload := []byte(`
        {
            "billId": "eqwptt",
            "invoiceUid": "44b2ef2a-edc6-4aed-87d3-01cf37ed2732",
            "amount": {
                "currency": "RUB",
                "value": "3.00"
            },
            "expirationDateTime": "2021-12-31T23:59:59+03:00",
            "status": {
                "value": "CREATED",
                "changedDateTime": "2021-02-05T10:21:17+03:00"
            },
            "comment": "My comment",
            "flags": [
                "TEST"
            ],
            "payUrl": "https://oplata.qiwi.com/form?invoiceUid=44b2ef2a-edc6-4aed-87d3-01cf37ed2732",
            "splits": [
                {
                    "type": "MERCHANT_DETAILS",
                    "siteUid": "Obuc-00",
                    "splitAmount": {
                        "currency": "RUB",
                        "value": "2.00"
                    },
                    "orderId": "dressesforwhite",
                    "comment": "Some purchase 1"
                },
                {
                    "type": "MERCHANT_DETAILS",
                    "siteUid": "Obuc-01",
                    "splitAmount": {
                        "currency": "RUB",
                        "value": "1.00"
                    },
                    "orderId": "shoesforvalya",
                    "comment": "Some purchase 2"
                }
            ]
        }
    `)

	// HTTP MOCK
	serv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		var pl *Payment
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

		err = json.Unmarshal(buf.Bytes(), pl)
		if err != nil {
			err = fmt.Errorf("[QIWI] %w: %s", ErrBadJSON, err)
			fmt.Fprintln(w, `{
                  "serviceName" : "payin-core",
                  "errorCode" : "validation.json",
                  "description" : "`+err.Error()+`",
                  "userMessage" : "Unmarshaling error",
                  "dateTime" : "2018-11-13T16:49:59.166+03:00",
                  "traceId" : "fd0e2a08c63ace83"
                }`)
			return
		}

		if buf.String() == string(payload) {
			err = fmt.Errorf("[QIWI] %w: %s", ErrBadJSON, err)
			fmt.Fprintln(w, `{
                  "serviceName" : "payin-core",
                  "errorCode" : "validation.json",
                  "description" : "`+err.Error()+`",
                  "userMessage" : "Bad json error",
                  "dateTime" : "2018-11-13T16:49:59.166+03:00",
                  "traceId" : "fd0e2a08c63ace83"
                }`)
			return
		}

		fmt.Fprintln(w, "{}")
	}))
	serv.Start()
	defer serv.Close()

	split := New("eqwptt", "SITEID", "TOKEN", serv.URL).AddSplit(NewSplit("MERCHID", NewAmountInRubles(2)))
	split = split.AddSplit(NewSplit("ANOTHERMERCHID", NewAmountInRubles(1)))
	err := split.CardRequest(context.TODO(), 300)
	if err != nil {
		t.Errorf("Splits method error: %s", err)
	}
}
