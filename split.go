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

// NewSplit without optional fields, see NewSplitExtra
func NewSplit(merchid string, a Amount) Split {
	return NewSplitExtra(merchid, a, "", "")
}

// NewSplitExtra extends to optional fields
func NewSplitExtra(merchid string, a Amount, orderid, c string) Split {
	s := Split{Type: MerchantDetails, SiteUID: merchid, Amount: a}

	if orderid != "" {
		s.OrderID = orderid
	}

	if c != "" {
		s.Comment = c
	}

	return s
}

// AddSplit adds split data to payment information
func (p *Payment) AddSplit(s Split) *Payment {
	p.Splits = append(p.Splits, &s)
	return p
}
