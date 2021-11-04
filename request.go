package qiwi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// requestTimeout sets timeout for context
// A duration string is a possibly signed sequence of decimal numbers.
const requestTimeout time.Duration = 15 * time.Second

var (
	// ErrBadJSON error throws when JSON marshal/unmarshal problem occurs.
	ErrBadJSON = errors.New("bad reply payload")
	// ErrBadStatusReply is bad gateway status code.
	ErrBadStatusReply = errors.New("bad status reply")
	// ErrReplyWithError business-logic error.
	ErrReplyWithError = errors.New("error in reply")
	// ErrBadSignature wrong signature error.
	ErrBadSignature = errors.New("wrong signature")
)

// global client to reuse existing connections.
var client http.Client

func init() {
	client = http.Client{}
}

// newRequest creates new http request to RSP.
func newRequest(ctx context.Context, method, link, apiToken string, payload []byte) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, link, bytes.NewBuffer(payload))
	req.Header.Set("User-Agent", userAgent+"/"+Version)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiToken)

	return req, err
}

// proceedRequest deal with data prep and preceedRequest
// handle response and pack all data back to our structure.
func proceedRequest(ctx context.Context, method, path string, p *Payment) error {
	var err error
	var payload []byte
	var result bytes.Buffer
	var req *http.Request
	var resp *http.Response

	payload, err = json.Marshal(p)
	if err != nil {
		return fmt.Errorf("[QIWI] %w: %s", ErrBadJSON, err)
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	link := p.apiLink + path
	req, err = newRequest(ctxTimeout, method, link, p.token, payload)
	if err != nil {
		return err
	}

	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusNoContent {
		return fmt.Errorf("[QIWI] %w: Resp http status code is #%d", ErrBadStatusReply, resp.StatusCode)
	}

	_, err = io.Copy(&result, resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(result.Bytes(), &p)
	if err != nil {
		return fmt.Errorf("[QIWI] %w: %s", ErrBadJSON, err)
	}

	return p.checkErrors(err)
}
