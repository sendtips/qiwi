package qiwi

import (
	"fmt"
	"time"
)

// PaymentType holds type of payment
type PaymentType string

const (
	// CardPayment for card payments
	CardPayment = "CARD"
	// TokenPayment for shadowed card numbers payment
	TokenPayment = "TOKEN"
	// ApplePayPayment ApplePay payment
	ApplePayPayment = "APPLE_PAY_TOKEN"
	// GooglePayPayment GooglePay payment
	GooglePayPayment = "GOOGLE_PAY_TOKEN"
)

// Payment main data structure, holds requests and responses on that requests from RSP
type Payment struct {
	token         string        `json:"-"`                   // Authtorisation token
	apiLink       string        `json:"-"`                   // APILink sets payment gateway domain, no trailing slash
	PublicKey     string        `json:"publicKey,omitempty"` // Merchant identification key	String	+
	SiteID        string        `json:"siteId,omitempty"`
	BillID        string        `json:"billId,omitempty"`        // Unique invoice identifier in merchant's system. It must be generated on your side with any means. It could be any sequence of digits and letters. Also you might use underscore _ and dash -. If not used, for each URL opening a new invoice is created. String(200)	-
	PaymentID     string        `json:"paymentId,omitempty"`     // Payment operation unique identifier in RSP's system
	CamptureID    string        `json:"captureId,omitempty"`     // Capture operation unique identifier in RSP's system
	RefundID      string        `json:"refundId,omitempty"`      // Refund operation unique identifier in RSP's system
	Amount        Amount        `json:"amount,omitempty"`        // Amount of customer order rounded down to 2 digits (always in rubles)
	PaymentMethod PaymentMethod `json:"paymentMethod,omitempty"` // Payment method
	Customer      Customer      `json:"customer,omitempty"`      // Information about the customer
	Creation      QIWITime      `json:"creationDateTime,omitempty"`
	NotifyDate    QIWITime      `json:"createddatetime,omitempty"` // Time used in Notify
	Expiration    QIWITime      `json:"expirationDateTime,omitempty"`
	Comment       string        `json:"comment,omitempty"`    // Comment to the invoice
	SuccessURL    string        `json:"successUrl,omitempty"` // URL for redirect from the QIWI form in case of successful payment. URL should be within the merchant's site.
	PayURL        string        `json:"payUrl,omitempty"`     // Payment page on QIWI site
	// extras[cf1]	Extra field to add any information to invoice data	URL-encoded string
	// extras[cf2]	Extra field to add any information to invoice data	URL-encoded string
	// extras[cf3]	Extra field to add any information to invoice data	URL-encoded string
	// extras[cf4]	Extra field to add any information to invoice data	URL-encoded string
	// extras[cf5]	Extra field to add any information to invoice data	URL-encoded string
	Reply

	QIWIError
}

// type PaymentMethod struct {
// 	//CardMethod
// 	ApplePayMethod
// 	//GooglePayMethod
// }

// PaymentMethod  holds payment type, card or applepay. googlepay data
type PaymentMethod struct {
	Type PaymentType `json:"type"` // Payment method type
	// "CARD" — payment card
	// "TOKEN" — card payment token
	// "APPLE_PAY_TOKEN" — encrypted Apple Pay payment token
	// "GOOGLE_PAY_TOKEN" — encrypted Google Pay payment token
	PAN string `json:"pan,omitempty"` // optional string(19) Card string //Card number. For type=CARD only

	ExpiryDate string `json:"expiryDate,omitempty"`
	//optional
	//string(5)
	//Card expiry date (MM/YY). For type=CARD only

	CVV string `json:"cvv2,omitempty"`
	//optional
	//string(4)
	//Card CVV2/CVC2. For type=CARD only

	Name string `json:"holderName,omitempty"`
	// optional
	// string(26)
	//Customer card holder (Latin letters). For type=CARD only

	Token ApplePayToken `json:"paymentData,omitempty"` // TODO OR paymentData
	//optional
	//string
	//Payment token string. For type=TOKEN, APPLE_PAY_TOKEN, GOOGLE_PAY_TOKEN only
	GooglePaymentToken string `json:"paymentToken,omitempty"`

	T3DS T3DS `json:"external3dSec,omitempty"`
	//optional
	//object
	//Payment data from Apple Pay or Google Pay.

}

// T3DS 3D-Secure
type T3DS struct {
	Type string `json:"type"`
	//require
	//string
	//Payment data type: APPLE_PAY or GOOGLE_PAY.

	OnlinePaymentCrypto string `json:"onlinePaymentCryptogram,omitempty"`
	//optional
	//string
	// Contents of "onlinePaymentCryptogram" field from decrypted Apple payment token. For type=APPLE_PAY only.

	Cryptogram string `json:"cryptogram,omitempty"`
	//optional
	//string
	// Contents of "cryptogram" from decrypted Google payment token. For type=GOOGLE_PAY only.

	ECIIndicator string `json:"eciIndicator,omitempty"`
	//optional
	//string(2)
	//ECI indicator. It should be sent if it is received in Apple (Google) payment token. Otherwise, do not send this parameter.
}

// Customer user related data
type Customer struct {
	Account string `json:"account,omitempty"`
	Email   string `json:"email,omitempty"`
	Phone   string `json:"phone,omitempty"`
}

// Reply from RSP
type Reply struct {
	Status Status `json:"status,omitempty"`
}

// Status of request
type Status struct {
	Value  string `json:"value,omitempty"`
	Date   string `json:"changedDateTime,omitempty"`
	Reason string `json:"reason,omitempty"`
}

// QIWIError holds error reply from a carrier
type QIWIError struct {
	Service     string    `json:"serviceName"` // Service name produced the error
	ErrCode     string    `json:"errorCode"`   // Error code
	Description string    `json:"description"` // Error description for RSP
	ErrMessage  string    `json:"userMessage"` // Error description for Customer
	ErrDate     time.Time `json:"dateTime"`    // Error date and time
	TraceID     string    `json:"traceId"`     // Error Log unique ID
}

// New create card payment session
func New(billId, siteid, token, apiLink string) *Payment {
	if apiLink == "" {
		apiLink = "https://api.qiwi.com/partner/v1/sites" // no trailing slash
	}
	return &Payment{SiteID: siteid, BillID: billId, apiLink: apiLink, token: token}
}

// checkErrors checks if errors is presented in reply
func (p *Payment) checkErrors(err error) error {

	if err == nil {
		if p.ErrCode != "" {
			err = fmt.Errorf("[QIWI] RSP Response %w: %s (%s)", ErrReplyWithError, p.Description, p.ErrCode)
		}
	}

	return err

}
