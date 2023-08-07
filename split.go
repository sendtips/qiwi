package qiwi

type SplitDataType string

const MerchantDetails SplitDataType = "MERCHANT_DETAILS"

type Split struct {
	Type    SplitDataType `json:"type"`              // String	Data type. Always MERCHANT_DETAILS (merchant details)
	SiteUID string        `json:"siteUid"`           // String	Merchant ID
	Amount  Amount        `json:"splitAmount"`       // Object	Merchant reimbursement data
	OrderID string        `json:"orderId,omitempty"` // String	Order number (optional)
	Comment string        `json:"comment,omitempty"` // string comment	String	Comment for the order (optional)
}

// AddSplit without optional fields, see SplitExtra.
func (p *Payment) Split(a Amount, merchid string) *Payment {
	return p.SplitExtra(a, merchid, "", "")
}

// SplitExtra extends to optional fields.
func (p *Payment) SplitExtra(a Amount, merchid, orderid, c string) *Payment {
	s := &Split{Type: MerchantDetails, SiteUID: merchid, Amount: a}

	if orderid != "" {
		s.OrderID = orderid
	}

	if c != "" {
		s.Comment = c
	}

	p.addSplit(s)

	return p
}

// addSplit adds split data to payment information.
func (p *Payment) addSplit(s *Split) {
	p.Splits = append(p.Splits, s)
}
