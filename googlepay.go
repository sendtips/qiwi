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
	_, err := w.Write(token)
	if err != nil {
		panic(err)
	}
	w.Close()

	return b.Bytes()
}

// GooglePay method executes Google Pay payment
func (p *Payment) GooglePay(ctx context.Context, amount int, token []byte) error {
	var err error

	p.PaymentMethod.Type = GooglePayPayment
	p.PaymentMethod.GooglePaymentToken = base64.StdEncoding.EncodeToString(zlibzompress(token))
	p.Amount = NewAmountInRubles(amount)

	// Make request link
	requestLink := fmt.Sprintf("/payin/v1/sites/%s/payments/%s", p.SiteID, p.PaymentID)

	err = proceedRequest(ctx, "PUT", requestLink, p)

	return p.checkErrors(err)

}
