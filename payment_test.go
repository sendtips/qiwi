package qiwi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusCodeUnmarshalJSON(t *testing.T) {
	tests := []struct {
		input string
		want  StatusCode
	}{
		{`{"value": "CREATED"}`, StatusCreated},
		{`{"value": "WAITING"}`, StatusWait},
		{`{"value": "COMPLETED"}`, StatusCompleted},
		{`{"value": "SUCCESS"}`, StatusOK},
		{`{"value": "DECLINE"}`, StatusFail},
		{`{"value": "DECLINED"}`, StatusFail},
		{`{"value": ""}`, StatusFail},
	}

	for _, test := range tests {
		var res Status
		_ = json.Unmarshal([]byte(test.input), &res)

		if res.Value != test.want {
			t.Errorf("Incorrect status value have: %v, want: %s", res.Value, test.want)
		}
	}
}

func TestBusinessLogicErrorHandling(t *testing.T) {
	// ErrReply
	serv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
				   "status": {"value": "DECLINED", "reason": "test reason"}
				}`)
	}))
	serv.Start()
	defer serv.Close()

	// Route request to mocked http server
	pay := New("billId", "SiteID", "TOKEN", serv.URL)

	err := proceedRequest(context.TODO(), "PUT", "/Init", pay)

	if !errors.Is(err, ErrReplyWithError) {
		t.Errorf("Remote error not parsed: %s", err)
	}

	if pay.Status.Value != StatusFail {
		t.Error("Wrong error value")
	}
}
