package qiwi

import (
	"encoding/json"
)

type NotifyType string

const (
	PaymentNotify   NotifyType = "PAYMENT"
	CaptureNotify   NotifyType = "CAPTURE"
	RefundNotify    NotifyType = "REFUND"
	CheckCardNotify NotifyType = "CHECK_CARD"
)

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
		return notify, err
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
