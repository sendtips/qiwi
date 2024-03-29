package qiwi

import (
	"fmt"
)

// apiLink holds QIWI API domain part of URL as string.
const apiLink string = "https://api.qiwi.com/partner" // no trailing slash

// PaymentType holds type of payment.
type PaymentType string

const (
	// CardPayment for card payments.
	CardPayment PaymentType = "CARD"
	// TokenPayment for shadowed card numbers payment.
	TokenPayment PaymentType = "TOKEN"
	// ApplePayPayment ApplePay payment.
	ApplePayPayment PaymentType = "APPLE_PAY_TOKEN"
	// GooglePayPayment GooglePay payment.
	GooglePayPayment PaymentType = "GOOGLE_PAY_TOKEN"
	// SBPPayment pay via SBP.
	SBPPayment PaymentType = "SBP"
)

// Payment main data structure, holds requests and responses on that requests from RSP.
type Payment struct {
	token         string         `json:"-"` // Authtorisation token
	apiLink       string         `json:"-"` // APILink sets payment gateway domain, no trailing slash
	siteid        string         // Same as SiteID, but used for requests, as SiteID sets in responses
	payid         string         // Same as BillID or PaymentID, but used only for requests, BillID sets in responses
	PublicKey     string         `json:"-"`                       // Merchant identification key	String	+
	SiteID        string         `json:"siteId,omitempty"`        // RSP site identifier
	BillID        string         `json:"billId,omitempty"`        // Unique invoice identifier in merchant's system. It must be generated on your side with any means. It could be any sequence of digits and letters. Also you might use underscore _ and dash -. If not used, for each URL opening a new invoice is created. String(200)	-
	PaymentID     string         `json:"paymentId,omitempty"`     // Payment operation unique identifier in RSP's system
	CamptureID    string         `json:"captureId,omitempty"`     // Capture operation unique identifier in RSP's system
	RefundID      string         `json:"refundId,omitempty"`      // Refund operation unique identifier in RSP's system
	Amount        Amount         `json:"amount,omitempty"`        // Amount of customer order rounded down to 2 digits (always in rubles)
	PaymentMethod *PaymentMethod `json:"paymentMethod,omitempty"` // Payment method
	Customer      *Customer      `json:"customer,omitempty"`      // Information about the customer
	Creation      *Time          `json:"creationDateTime,omitempty"`
	NotifyDate    *Time          `json:"createddatetime,omitempty"` // Time used in Notify
	Expiration    *Time          `json:"expirationDateTime,omitempty"`
	Comment       string         `json:"comment,omitempty"`     // Comment to the invoice
	CallbackURL   string         `json:"callbackUrl,omitempty"` // Callback URL used to receive notification
	SuccessURL    string         `json:"successUrl,omitempty"`  // URL for redirect from the QIWI form in case of successful payment. URL should be within the merchant's site.
	PayURL        string         `json:"payUrl,omitempty"`      // Payment page on QIWI site
	Req           *Requirements  `json:"requirements,omitempty"`
	CustomField   *CustomField   `json:"customFields,omitempty"`
	Flags         []string       `json:"flags,omitempty"`
	Status        *Status        `json:"status,omitempty"`
	Splits        []*Split       `json:"splits,omitempty"` // https://developer.qiwi.com/en/payments/#payments_split

	Error
}

// PaymentMethod  holds payment type, card or applepay. googlepay data.
type PaymentMethod struct {
	Type PaymentType `json:"type,omitempty"` // Payment method type
	// "CARD" — payment card
	// "TOKEN" — card payment token
	// "APPLE_PAY_TOKEN" — encrypted Apple Pay payment token
	// "GOOGLE_PAY_TOKEN" — encrypted Google Pay payment token

	PAN string `json:"pan,omitempty"` // optional string(19) Card string //Card number. For type=CARD only

	ExpiryDate string `json:"expiryDate,omitempty"`
	// optional
	// string(5)
	// Card expiry date (MM/YY). For type=CARD only

	CVV string `json:"cvv2,omitempty"`
	// optional
	// string(4)
	// Card CVV2/CVC2. For type=CARD only

	Name string `json:"holderName,omitempty"`
	// optional
	// string(26)
	// Customer card holder (Latin letters). For type=CARD only

	ApplePayToken *PKPaymentToken `json:"paymentData,omitempty"`
	// optional
	// string
	// Payment token string. For type=TOKEN, APPLE_PAY_TOKEN, GOOGLE_PAY_TOKEN only
	Token string `json:"paymentToken,omitempty"`

	T3DS *T3DS `json:"external3dSec,omitempty"`
	// optional
	// object
	// Payment data from Apple Pay or Google Pay.
}

// T3DS 3D-Secure.
type T3DS struct {
	Type string `json:"type"`
	// require
	// string
	// Payment data type: APPLE_PAY or GOOGLE_PAY.

	OnlinePaymentCrypto string `json:"onlinePaymentCryptogram,omitempty"`
	// optional
	// string
	// Contents of "onlinePaymentCryptogram" field from decrypted Apple payment token. For type=APPLE_PAY only.

	Cryptogram string `json:"cryptogram,omitempty"`
	// optional
	// string
	// Contents of "cryptogram" from decrypted Google payment token. For type=GOOGLE_PAY only.

	ECIIndicator string `json:"eciIndicator,omitempty"`
	// optional
	// string(2)
	// ECI indicator. It should be sent if it is received in Apple (Google) payment token. Otherwise, do not send this parameter.
}

// CustomField in this API used only
// to set callback url.
type CustomField struct {
	CallbackURL string `json:"invoice_callback_url,omitempty"`
}

type Requirements struct {
	SBP *SBPData `json:"sbp,omitempty"`
}

// Customer user related data.
type Customer struct {
	Account string `json:"account,omitempty"`
	Email   string `json:"email,omitempty"`
	Phone   string `json:"phone,omitempty"`
}

// StatusCode operation status reflects its current state.
type StatusCode string

const (
	StatusCreated   StatusCode = "CREATED"   // For invoices only one status is used
	StatusWait      StatusCode = "WAITING"   // Awaiting for 3DS authentication API responses)
	StatusCompleted StatusCode = "COMPLETED" // Request for authentication is successfully processed API responsess
	StatusOK        StatusCode = "SUCCESS"   // Request for authentication is successfully processed Notifications
	StatusFail      StatusCode = "DECLINE"   // Request for payment confirmation is rejected Notifications, API responses
)

func (sc *StatusCode) UnmarshalJSON(b []byte) (err error) {
	s := string(b[1 : len(b)-1])
	if s == "DECLINED" || s == "" { // There are two responses from QIWI API
		s = "DECLINE" // with the same meaning DECLINE and DECLINED.
	}
	*sc = StatusCode(s)
	return err
}

// StatusError API errors describe a reason for rejection of the operation.
type StatusError string

const (
	StatusInvalidState                StatusError = "INVALID_STATE"                  // Incorrect transaction status
	StatusInvalidAmount               StatusError = "INVALID_AMOUNT"                 // Incorrect payment amount
	StatusDeclinedMPI                 StatusError = "DECLINED_BY_MPI"                // Rejected by MPI
	StatusDeclinedFraud               StatusError = "DECLINED_BY_FRAUD"              // Rejected by fraud monitoring
	StatusGatawayIntegrationError     StatusError = "GATEWAY_INTEGRATION_ERROR"      // Acquirer integration error
	StatusGatawayTechnicalError       StatusError = "GATEWAY_TECHNICAL_ERROR"        // Technical error on acquirer side
	StatusAcquiringMPITechError       StatusError = "ACQUIRING_MPI_TECH_ERROR"       // Technical error on 3DS authentication
	StatusAcquiringGatawayTechError   StatusError = "ACQUIRING_GATEWAY_TECH_ERROR"   // Technical error
	StatusAcquiringAcquirerTechError  StatusError = "ACQUIRING_ACQUIRER_ERROR"       // Technical error
	StatusAcquiringAuthTechnicalError StatusError = "ACQUIRING_AUTH_TECHNICAL_ERROR" // Error on funds authentication
	StatusAcquiringIssuerNotAvailable StatusError = "ACQUIRING_ISSUER_NOT_AVAILABLE" // Issuer error. Issuer is not available at the moment
	StatusAcquiringSuspectedFraud     StatusError = "ACQUIRING_SUSPECTED_FRAUD"      // Issuer error. Fraud suspicion
	StatusAcquiringLimitExceeded      StatusError = "ACQUIRING_LIMIT_EXCEEDED"       // Issuer error. Some limit exceeded
	StatusAcquiringNotPermitted       StatusError = "ACQUIRING_NOT_PERMITTED"        // Issuer error. Operation not allowed
	StatusAcquiringIncorrectCVV       StatusError = "ACQUIRING_INCORRECT_CVV"        // Issuer error. Incorrect CVV
	StatusAcquiringExpiredCard        StatusError = "ACQUIRING_EXPIRED_CARD"         // Issuer error. Incorrect card expiration date
	StatusAcquiringInvalidCard        StatusError = "ACQUIRING_INVALID_CARD"         // Issuer error. Verify card data
	StatusAcquiringInsufficientFunds  StatusError = "ACQUIRING_INSUFFICIENT_FUNDS"   // Issuer error. Not enough funds
	StatusAcquiringUnknown            StatusError = "ACQUIRING_UNKNOWN"              // Unknown error
	StatusBillAlreadyPaid             StatusError = "BILL_ALREADY_PAID"              // Invoice is already paid
	StatusPayinProcessingError        StatusError = "PAYIN_PROCESSING_ERROR"         // Payment processing error
)

// Status of request.
type Status struct {
	Value        StatusCode  `json:"value,omitempty"`
	Date         *Time       `json:"changedDateTime,omitempty"`
	Reason       StatusError `json:"reason,omitempty"`
	ReasonNotify StatusError `json:"reasonCode,omitempty"`
	Message      string      `json:"reasonMessage,omitempty"`
	ErrCode      string      `json:"errorCode,omitempty"`
}

// Error holds error reply from a carrier.
type Error struct {
	Service     string `json:"serviceName,omitempty"` // Service name produced the error
	ErrCode     string `json:"errorCode,omitempty"`   // Error code
	Description string `json:"description,omitempty"` // Error description for RSP
	ErrMessage  string `json:"userMessage,omitempty"` // Error description for Customer
	ErrDate     *Time  `json:"dateTime,omitempty"`    // Error date and time
	TraceID     string `json:"traceId,omitempty"`     // Error Log unique ID
}

// New create card payment session.
func New(billID, siteID, token, endpoint string) *Payment {
	if endpoint == "" {
		endpoint = apiLink
	}
	return &Payment{siteid: siteID, payid: billID, BillID: billID, apiLink: endpoint, token: token, PaymentMethod: &PaymentMethod{}}
}

// checkErrors checks if errors is presented in reply.
// Note, there are no checks for Notify type business-logic errors
// you must handle it on your own.
func (p *Payment) checkErrors(err error) error {
	if err == nil {
		if p.ErrCode != "" {
			err = fmt.Errorf("[QIWI] Error in response %w: %s (%s)", ErrReplyWithError, p.Description, p.ErrCode)
		} else if p.Status != nil {
			if p.Status.Reason != "" {
				err = fmt.Errorf("[QIWI] Error in response %w: %s (%s)", ErrReplyWithError, p.Status.Value, p.Status.Reason)
			}
		}
	}

	return err
}
