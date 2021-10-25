package qiwi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Signature type holds keys to calculate signature
type Signature struct {
	Key, Message []byte
	Hash         string
	Valid        bool
}

// NewSignature return new signature
func NewSignature(key, hash string) *Signature {
	return &Signature{Key: []byte(key), Hash: hash}
}

// sign calculates checksum
func (s *Signature) sign() bool {
	mac := hmac.New(sha256.New, s.Key)
	mac.Write(s.Message)
	expectedMAC := mac.Sum(nil)

	if hex.EncodeToString(expectedMAC) == s.Hash {
		s.Valid = true
		return true
	}

	return false
}

// Verify HMAC-SHA256 signature hash used in Notify type
func (s *Signature) Verify(p Notify) bool {

	switch p.Type {
	case PaymentNotify:
		// payment.paymentId|payment.createdDateTime|payment.amount.value
		s.Message = []byte(fmt.Sprintf("%s|%s|%s", p.Payment.PaymentID, p.Payment.NotifyDate, p.Payment.Amount.Value))

	case CaptureNotify:
		// capture.captureId|capture.createdDateTime|capture.amount.value
		s.Message = []byte(fmt.Sprintf("%s|%s|%s", p.Capture.CamptureID, p.Capture.NotifyDate, p.Capture.Amount.Value))

	case RefundNotify:
		// refund.refundId|refund.createdDateTime|refund.amount.value
		s.Message = []byte(fmt.Sprintf("%s|%s|%s", p.Refund.RefundID, p.Refund.NotifyDate, p.Refund.Amount.Value))

	case CheckCardNotify:
		// checkPaymentMethod.requestUid|checkPaymentMethod.checkOperationDate
		s.Message = []byte(fmt.Sprintf("%s|%s", p.CheckCard.RequestID, p.CheckCard.CheckDate))

	default:
		// fallback to PaymentType
		s.Message = []byte(fmt.Sprintf("%s|%s|%s", p.Payment.PaymentID, p.Payment.NotifyDate, p.Payment.Amount.Value))
	}

	return s.sign()
}
