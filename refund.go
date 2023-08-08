package qiwi

type RefundSplits struct {
	Split
	Commission *RefundSplitCommission `json:"splitCommissions,omitempty"`
}

// RefundSplitCommission contains commission information
type RefundSplitCommission struct {
	Amount  *Amount `json:"merchantCms,omitempty"`
	UserCms string  `json:"userCms,omitempty"`
}
