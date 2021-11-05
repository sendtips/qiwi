package qiwi

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
)

// zlibzompress via zlib token data for googlepay.
func zlibzompress(token string) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, _ = w.Write([]byte(token))
	w.Close()

	return b.Bytes()
}

// GooglePay method executes Google Pay payment.
func (p *Payment) GooglePay(ctx context.Context, amount int, token string) error {
	var err error

	// p.PaymentMethod = &PaymentMethod{}
	p.PaymentMethod.Type = GooglePayPayment
	p.PaymentMethod.Token = base64.StdEncoding.EncodeToString(zlibzompress(token))
	p.Amount = NewAmountInRubles(amount)

	// Make request link
	requestLink := fmt.Sprintf("/payin/v1/sites/%s/payments/%s", p.siteid, p.payid)

	err = proceedRequest(ctx, http.MethodPut, requestLink, p)

	return p.checkErrors(err)
}
