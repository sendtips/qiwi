package qiwi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// maxBodyBytes limit inbound payload to 64kb
const maxBodyBytes = int64(65536)

// NotifyType from RSP
type NotifyType string

const (
	// PaymentNotify for payment notification
	PaymentNotify NotifyType = "PAYMENT"
	// CaptureNotify for capture (2 stage) payments
	CaptureNotify NotifyType = "CAPTURE"
	// RefundNotify for refunds
	RefundNotify NotifyType = "REFUND"
	// CheckCardNotify for card check requests
	CheckCardNotify NotifyType = "CHECK_CARD"
)

// Notify holds incoming data from RSP
type Notify struct {
	Type      NotifyType `json:"type"` // Notification type
	Payment   Payment    `json:"payment,omitempty"`
	Capture   Payment    `json:"capture,omitempty"`
	Refund    Payment    `json:"refund,omitempty"`
	CheckCard Card       `json:"checkPaymentMethod,omitempty"`
	Version   string     `json:"version"` // Notification version
}

// NewNotify returns Notify data from bytes
func NewNotify(signkey, sign string, payload []byte) (Notify, error) {
	var notify Notify
	var err error

	err = json.Unmarshal(payload, &notify)
	if err != nil {
		return notify, fmt.Errorf("[QIWI] Notify: %w (%s)", ErrBadJSON, err)
	}

	// Check signature
	if sign != "" {
		sig := NewSignature(signkey, sign)
		if !sig.Verify(notify) {
			err = ErrBadSignature
		}
	}

	return notify, err
}

// NotifyParseHTTPRequest parses http request and returns Notify
// And protects against a malicious client streaming
// us an endless request body
func NotifyParseHTTPRequest(signkey, sign string, w http.ResponseWriter, r *http.Request) (Notify, error) {
	var payload bytes.Buffer

	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	_, err := io.Copy(&payload, r.Body)
	if err != nil {
		return Notify{}, fmt.Errorf("[QIWI] Notify payload http parser: %w: %s", ErrBadJSON, err)
	}

	return NewNotify(signkey, sign, payload.Bytes())

}
