package qiwi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// expirationTime set session expiration time for card payment.
const expirationTime time.Duration = 5 * time.Minute

// T3DSStatus 3D Secure status.
type T3DSStatus string

const (
	Ok3DS   T3DSStatus = "PASSED"     // Ok3DS 3-D Secure passed
	Fail3DS T3DSStatus = "NOT_PASSED" // Fail3DS 3-D Secure not passed)
	None3DS T3DSStatus = "WITHOUT"    // None3DS 3-D Secure not required
)

// Card holds data for card payment.
type Card struct {
	CheckDate  Time       `json:"checkOperationDate"`      // System date of the operation
	RequestID  string     `json:"requestUid"`              // Card verification operation unique identifier String(200)
	Status     Status     `json:"status"`                  // Card verification status	String
	Valid      bool       `json:"isValidCard"`             // Logical flag means card is valid for purchases Bool
	T3DS       T3DSStatus `json:"threeDsStatus"`           // Information on customer authentication status.
	CardMethod CardMethod `json:"paymentMethod,omitempty"` // Payment method data for card
	CardInfo   CardInfo   `json:"cardInfo,omitempty"`      // Card information
	CardToken  CardToken  `json:"createdToken,omitempty"`  // Payment token data
}

// CardMethod carry card data.
type CardMethod struct {
	Type       PaymentType `json:"type"`           // Payment method type
	Payment    string      `json:"maskedPan"`      // Masked card PAN
	Expiration string      `json:"cardExpireDate"` // Card expiration date (MM/YY)
	Name       string      `json:"cardHolder"`     // Cardholder name
}

// CardInfo additional information about card.
type CardInfo struct {
	Country       string `json:"issuingCountry"`       // Issuer country code	String(3)
	Bank          string `json:"issuingBank"`          // Issuer name	String
	PaymentSystem string `json:"paymentSystem"`        // Card's payment system	String
	CardType      string `json:"fundingSource"`        // Card's type (debit/credit/..)	String
	Product       string `json:"paymentSystemProduct"` // Card's category	String
}

// CardToken shadowed card.
type CardToken struct {
	Token          string `json:"token"`       // Card payment token	String
	Name           string `json:"name"`        // Masked card PAN for which payment token issued	String
	ExpirationDate Time   `json:"expiredDate"` // Payment token expiration date. ISO-8601
	Account        string `json:"account"`     // Customer account for which payment token issued	String
}

// CardRequest request payment session on RSP site.
func (p *Payment) CardRequest(ctx context.Context, amount int) error {
	p.Amount = NewAmountInRubles(amount)
	p.Expiration = NowInMoscow().Add(expirationTime)
	p.Flags.Flags = []string{"SALE"} // one-step payment

	p.PaymentMethod = nil // need to make query

	// Make request link
	requestLink := fmt.Sprintf("/payin/v1/sites/%s/bills/%s", p.siteid, p.payid)

	return proceedRequest(ctx, http.MethodPut, requestLink, p)
}

// PublicLink an easy way to integrate with the QIWI payment form.
// When the form is opened, the customer is automatically billed for the order.
// Parameters of the invoice are sent unprotected in the link.
// A payment form with a choice of payment methods is shown to the customer.
func PublicLink(pubKey, comment string, amount int) string {
	// link qiwi payment link structure
	const link string = `https://oplata.qiwi.com/create?publicKey=%s&comment=%s&amount=%s`

	am := NewAmountInRubles(amount)

	return fmt.Sprintf(link, url.PathEscape(pubKey), url.PathEscape(comment), am.Value)
}
