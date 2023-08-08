package qiwi

// RefundSplits contains split refund.
type RefundSplits struct {
	Split
	Commission *RefundSplitCommission `json:"splitCommissions,omitempty"`
}

// RefundSplitCommission contains commission information.
type RefundSplitCommission struct {
	Amount  *Amount `json:"merchantCms,omitempty"`
	UserCms string  `json:"userCms,omitempty"`
}
