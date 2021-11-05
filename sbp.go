package qiwi

import (
	"context"
	"fmt"
	"net/http"
)

// SBP contains Image and URL for SBP payment.
type SBPData struct {
	QRCodeID string   `json:"qrcId,omitempty"`
	Image    SBPImage `json:"image,omitempty"`
	PayURL   string   `json:"payload,omitempty"`
}

// SBPImage holds base64 encoded image of SBP QRCode and its mimetype.
type SBPImage struct {
	MimeType string `json:"mediaType,omitempty"`
	Picture  string `json:"content,omitempty"`
}

// SBP returns SBP payment session data.
// Comment appears in bank payment.
func (p *Payment) SBP(ctx context.Context, amount int, comment string) error {
	p.PaymentMethod.Type = SBPPayment
	p.Amount = NewAmountInRubles(amount)
	p.Comment = comment

	requestLink := fmt.Sprintf("/payin/v1/sites/%s/payments/%s", p.siteid, p.payid)

	return proceedRequest(ctx, http.MethodPut, requestLink, p)
}
