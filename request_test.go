package qiwi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	trans := New("billId", "SiteID", "TOKEN", "")
	dummy := &Payment{}

	if reflect.TypeOf(trans) != reflect.TypeOf(dummy) {
		t.Errorf("New() wrong return %T must be %T type", reflect.TypeOf(trans), reflect.TypeOf(dummy))
	}
}

func TestNewRequest(t *testing.T) {
	req, _ := newRequest(context.TODO(), "POST", "http://example.com/Init", "TOKEN", nil)
	dummy := &http.Request{}

	if reflect.TypeOf(req) != reflect.TypeOf(dummy) {
		t.Errorf("newRequest() wrong return %T must be %T type", reflect.TypeOf(req), reflect.TypeOf(dummy))
	}
}

func TestRequestHeaders(t *testing.T) {
	link, _ := url.Parse("http://example.com/Init")
	req, _ := newRequest(context.TODO(), "POST", link.String(), "TOKEN", nil)

	if req.URL.Hostname() != link.Hostname() {
		t.Error("Wrong hostname")
	}

	if req.Header.Get("Content-Type") != "application/json" {
		t.Error("Wrong content-type")
	}

	if req.Header.Get("Accept") != "application/json" {
		t.Error("Wrong Accept")
	}

	if req.Header.Get("Bearer") != "TOKEN" {
		t.Error("No authorization token")
	}
}

func TestProceedBadJSONRequest(t *testing.T) {

	// ErrBadJSON
	serv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "{{{{bad json")
	}))
	//serv_badjson.Listener = listner
	serv.Start()
	defer serv.Close()

	// Route request to mocked http server
	payload := New("billId", "SiteID", "TOKEN", serv.URL)

	err := proceedRequest(context.Background(), "POST", "/Init", payload)

	if !errors.Is(err, ErrBadJSON) {
		t.Errorf("Wrong error for bad JSON return: %s", err)
	}
}
func TestProceedBadHTTPStatusRequest(t *testing.T) {
	// ErrBadStatusReply
	serv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Something wrong", http.StatusInternalServerError)
	}))
	serv.Start()
	defer serv.Close()

	// Route request to mocked http server
	payload := New("billId", "SiteID", "TOKEN", serv.URL)

	err := proceedRequest(context.TODO(), "POST", "/Init", payload)

	if !errors.Is(err, ErrBadStatusReply) {
		t.Errorf("Wrong error for error HTTP error code response: %s", err)
	}
}
func TestProceeRequestWithError(t *testing.T) {
	// ErrReply
	serv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
			   "serviceName":"payin-core",
			   "errorCode":"validation.error",
			   "description":"Validation error",
			   "userMessage":"Validation error",
			   "dateTime":"2018-11-13T16:49:59.166+03:00",
			   "traceId":"fd0e2a08c63ace83"
			}`)
	}))
	serv.Start()
	defer serv.Close()

	// Route request to mocked http server
	payload := New("billId", "SiteID", "TOKEN", serv.URL)

	err := proceedRequest(context.TODO(), "PUT", "/Init", payload)

	if !errors.Is(err, ErrReplyWithError) {
		t.Errorf("Remote error not parsed: %s", err)
	}

	if payload.Description != "Validation error" {
		t.Error("Wrong RSP error message")
	}

	if payload.ErrCode != "validation.error" {
		t.Error("Wrong RSP error code")
	}
}

func TestProceedRequest(t *testing.T) {
	// GoodRequest
	serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reply := `{}`
		fmt.Fprintln(w, reply)
	}))
	defer serv.Close()

	// Route request to mocked http server
	payload := New("billId", "SiteID", "TOKEN", serv.URL)

	err := proceedRequest(context.TODO(), "POST", "/Init", payload)

	if err != nil {
		t.Errorf("Error shoud be empty: %s", err)
	}
}
