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
	"time"
)

func TestCardRequest(t *testing.T) {
	// Expected reply from QIWI
	// HTTP/1.1 200 OK
	// Content-Type: application/json
	reply := `
    {
        "siteId": "test-01",
        "billId": "gg",
        "amount": {
            "currency": "RUB",
            "value": 42.24
        },
        "status": {
            "value": "WAITING",
            "changedDateTime": "2019-08-28T16:26:36.835+03:00"
        },
        "customFields": {},
        "comment": "Spasibo",
        "creationDateTime": "2019-08-28T16:26:36.835+03:00",
        "expirationDateTime": "2019-09-13T14:30:00+03:00",
        "payUrl": "https://oplata.qiwi.com/form/?invoice_uid=78d60ca9-7c99-481f-8e51-0100c9012087"
    }`

	errReply := `
	{
		  "serviceName" : "payin-core",
		  "errorCode" : "validation.wrongmethod",
		  "description" : "Wrong method",
		  "userMessage" : "Validation error",
		  "dateTime" : "2018-11-13T16:49:59.166+03:00",
		  "traceId" : "fd0e2a08c63ace83"
		}
	`

	serv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			var buf bytes.Buffer
			var payload Payment

			_, _ = io.Copy(&buf, r.Body)

			_ = json.Unmarshal(buf.Bytes(), &payload)
			isSale := false
			for _, flag := range payload.Flags {
				if flag == "SALE" {
					isSale = true
					break
				}
			}
			if !isSale {
				fmt.Println(w, errReply)
				return
			}

			fmt.Fprintln(w, reply)
		default:
			fmt.Fprint(w, errReply)
		}
	}))

	serv.Start()
	defer serv.Close()

	// Route request to mocked http server
	pay := New("billId", "SiteID", "TOKEN", serv.URL)
	err := pay.CardRequest(context.TODO(), 100)
	if err != nil {
		t.Errorf("CardRequest error: %s", err)
	}

	// dont need to pass this
	// if pay.PaymentMethod.Type != CardPayment {
	// 	t.Errorf("Wrong payment type %s", pay.PaymentMethod.Type)
	// }

	if pay.PayURL != "https://oplata.qiwi.com/form/?invoice_uid=78d60ca9-7c99-481f-8e51-0100c9012087" {
		t.Error("PayURL not received")
	}
}

func TestLocationTime(t *testing.T) {
	pay := New("billID", "siteID", "token", "")
	_ = pay.CardRequest(context.TODO(), 100)

	// Moscow time
	msktz, _ := time.LoadLocation("Europe/Moscow")

	msktime := time.Now().In(msktz)

	if pay.Expiration.Before(msktime) {
		t.Error("Bad expiration time", pay.Expiration, msktime)
	}
}

func TestPublicLink(t *testing.T) {
	const wantLink = `https://oplata.qiwi.com/create?publicKey=pubKey&comment=some%20comment&amount=3.01`

	payLink := PublicLink("pubKey", "some comment", 301)

	if payLink != wantLink {
		t.Errorf("Wrong payment link: %s, want %s", payLink, wantLink)
	}
}
