package qiwi

// AddCallbackURLToBill adds custom callback URL to Bill.
// Note, bills used in card payments.
func AddCallbackURLToBill(p *Payment, url string) {
	p.CustomField = &CustomField{CallbackURL: url}
}

// AddCallbackURLToPayment adds custom callback URL to Payment.
func AddCallbackURLToPayment(p *Payment, url string) {
	p.CallbackURL = url
}
