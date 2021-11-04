package qiwi

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/base64"
	"fmt"
)

// zlibzompress via zlib token data for googlepay
func zlibzompress(token []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, _ = w.Write(token)
	w.Close()

	return b.Bytes()
}

// GooglePay method executes Google Pay payment
func (p *Payment) GooglePay(ctx context.Context, amount int, token []byte) error {
	var err error

	p.PaymentMethod = &PaymentMethod{}
	p.PaymentMethod.Type = GooglePayPayment
	p.PaymentMethod.Token = base64.StdEncoding.EncodeToString(zlibzompress(token))
	p.Amount = NewAmountInRubles(amount)

	// Make request link
	requestLink := fmt.Sprintf("/payin/v1/sites/%s/payments/%s", p.SiteID, p.BillID)

	err = proceedRequest(ctx, "PUT", requestLink, p)

	return p.checkErrors(err)

}
