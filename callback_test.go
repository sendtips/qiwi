package qiwi

import (
	"testing"
)

func TestAddCallbackURLToBill(t *testing.T) {
	const sampleURL = "https://example.com"

	pay := New("billId", "SiteID", "TOKEN", "")

	if pay.CustomField != nil {
		t.Error("CallbackURL set by default at Bill step")
	}

	AddCallbackURLToBill(pay, sampleURL)

	if pay.CustomField.CallbackURL != sampleURL {
		t.Errorf("CallbackURL wrong value at Bill step %s, shoud be %s", pay.CustomField.CallbackURL, sampleURL)
	}
}

func TestAddCallbackURLToPayment(t *testing.T) {
	const sampleURL = "https://example.com"

	pay := New("billId", "SiteID", "TOKEN", "")

	if pay.CallbackURL != "" {
		t.Error("CallbackURL set by default at Bill step")
	}

	AddCallbackURLToPayment(pay, sampleURL)

	if pay.CallbackURL != sampleURL {
		t.Errorf("CallbackURL wrong value at Bill step %s, shoud be %s", pay.CallbackURL, sampleURL)
	}
}
