package qiwi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// ApplePayToken carry payment token
type ApplePayToken struct {
	//	Type        string        `json:"type"`
	PaymentData ApplePayTokenData `json:"paymentData"`
}

// ApplePayTokenData encrypted token data structure
type ApplePayTokenData struct {
	Version   string   `json:"version"`
	Data      string   `json:"data"`
	Header    APHeader `json:"header"`
	Signature string   `json:"signature"`
}

// APHeader internal ApplePayTokenData structure
type APHeader struct {
	PubKey        string `json:"ephemeralPublicKey"`
	PubKeyHash    string `json:"publicKeyHash"`
	TransactionID string `json:"transactionId"`
}

// decodeBase64 return base64 decoded string
// used internally to "unpack" applepay token
func decodeBase64(enc string) ([]byte, error) {
	dec, err := base64.StdEncoding.DecodeString(enc)
	return dec, err
}

// ApplePay executes payment via ApplePay
// Pass context, amount and ApplePay token string
func (p *Payment) ApplePay(ctx context.Context, amount int, token string) (err error) {

	// Decode token from base64
	data, err := decodeBase64(token)
	if err != nil {
		return fmt.Errorf("[QIWI] %w: %s", ErrBadJSON, err)
	}

	// Parse JSON data+-
	err = json.Unmarshal(data, &p.PaymentMethod.ApplePayToken)
	if err != nil {
		return fmt.Errorf("[QIWI] %w: %s", ErrBadJSON, err)
	}

	p.PaymentMethod.Type = ApplePayPayment
	p.Amount = NewAmountInRubles(amount)

	// Make request link
	requestLink := fmt.Sprintf("/payin/v1/sites/%s/payments/%s", p.SiteID, p.PaymentID)

	// Send request
	err = proceedRequest(ctx, "PUT", requestLink, p)

	return p.checkErrors(err)

}
