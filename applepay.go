package qiwi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

// PKPaymentToken is ApplePay payment token structure
// see https://developer.apple.com/library/archive/documentation/PassKit/Reference/PaymentTokenJSON/PaymentTokenJSON.html
type PKPaymentToken struct {
	Version   string    `json:"version"`
	Data      string    `json:"data"`
	Header    *APHeader `json:"header"`
	Signature string    `json:"signature"`
}

// APHeader holds Header of PKPaymentToken.
type APHeader struct {
	AppData       string `json:"applicationData,omitempty"`    // optional, HEX-string
	Key           string `json:"wrappedKey,omitempty"`         // used only for RSA_v1
	PubKey        string `json:"ephemeralPublicKey,omitempty"` // used only for EC_v1
	PubKeyHash    string `json:"publicKeyHash"`
	TransactionID string `json:"transactionId"`
}

// ApplePay executes payment via ApplePay
// pass context, amount and ApplePay token string.
func (p *Payment) ApplePay(ctx context.Context, amount int, token string) (err error) {
	// Decode token from base64
	data, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return fmt.Errorf("[QIWI] %w: %s", ErrBadJSON, err)
	}

	// p.PaymentMethod = &PaymentMethod{}
	p.PaymentMethod.ApplePayToken = &PKPaymentToken{Header: &APHeader{}}

	// Parse JSON data+-
	err = json.Unmarshal(data, p.PaymentMethod.ApplePayToken)
	if err != nil {
		return fmt.Errorf("[QIWI] %w: %s", ErrBadJSON, err)
	}

	p.PaymentMethod.Type = ApplePayPayment
	p.Amount = NewAmountInRubles(amount)

	// Make request link
	requestLink := fmt.Sprintf("/payin/v1/sites/%s/payments/%s", p.siteid, p.payid)

	// Send request
	err = proceedRequest(ctx, http.MethodPut, requestLink, p)

	return p.checkErrors(err)
}
