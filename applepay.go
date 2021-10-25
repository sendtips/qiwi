package qiwi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// type ApplePayMethod struct {
// 	Type string `json:"type"` // Payment method type
// 	// "CARD" — payment card
// 	// "TOKEN" — card payment token
// 	// "APPLE_PAY_TOKEN" — encrypted Apple Pay payment token
// 	// "GOOGLE_PAY_TOKEN" — encrypted Google Pay payment token
// 	PAN string `json:"pan,omitempty"` // optional string(19) Card string //Card number. For type=CARD only
//
// 	ExpiryDate string `json:"expiryDate,omitempty"`
// 	//optional
// 	//string(5)
// 	//Card expiry date (MM/YY). For type=CARD only
//
// 	CVV string `json:"cvv2,omitempty"`
// 	//optional
// 	//string(4)
// 	//Card CVV2/CVC2. For type=CARD only
//
// 	Name string `json:"holderName,omitempty"`
// 	// optional
// 	// string(26)
// 	//Customer card holder (Latin letters). For type=CARD only
//
// 	Token string `json:"paymentToken,omitempty"`
// 	//optional
// 	//string
// 	//Payment token string. For type=TOKEN, APPLE_PAY_TOKEN, GOOGLE_PAY_TOKEN only
//
// 	T3DS T3DS `json:"external3dSec,omitempty"`
// 	//optional
// 	//object
// 	//Payment data from Apple Pay or Google Pay.
//
// }
//
// type T3DS struct {
// 	Type string `json:"type"`
// 	//require
// 	//string
// 	//Payment data type: APPLE_PAY or GOOGLE_PAY.
//
// 	OnlinePaymentCrypto string `json:"onlinePaymentCryptogram,omitempty"`
// 	//optional
// 	//string
// 	// Contents of "onlinePaymentCryptogram" field from decrypted Apple payment token. For type=APPLE_PAY only.
//
// 	Cryptogram string `json:"cryptogram,omitempty"`
// 	//optional
// 	//string
// 	// Contents of "cryptogram" from decrypted Google payment token. For type=GOOGLE_PAY only.
//
// 	ECIIndicator string `json:"eciIndicator,omitempty"`
// 	//optional
// 	//string(2)
// 	//ECI indicator. It should be sent if it is received in Apple (Google) payment token. Otherwise, do not send this parameter.
// }

type ApplePayToken struct {
	//	Type        string        `json:"type"`
	PaymentData ApplePayTokenData `json:"paymentData"`
}

type ApplePayTokenData struct {
	Version   string   `json:"version"`
	Data      string   `json:"data"`
	Header    APHeader `json:"header"`
	Signature string   `json:"signature"`
}

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

func (p *Payment) ApplePay(ctx context.Context, amount int, token string) (err error) {

	// Decode token from base64
	data, err := decodeBase64(token)
	if err != nil {
		return fmt.Errorf("[QIWI] %w: %s", ErrBadJSON, err)
	}

	// Parse JSON data+-
	err = json.Unmarshal(data, &p.PaymentMethod.Token)
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
