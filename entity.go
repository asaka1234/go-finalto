package finalto

// 返回的结构化
type RequestForPositionAck struct {
	PosMaintRptID      string `json:"pos_maint_rpt_id"  config:"pos_maint_rpt_id"   default:""`
	PosReqID           string `json:"pos_req_id,omitempty" config:"pos_req_id"`
	TotalNumPosReports int    `json:"total_num_pos_reports,omitempty"  config:"total_num_pos_reports"   default:""`
	PosReqResult       int    `json:"pos_req_result"  config:"pos_req_result"   default:""`
	PosReqStatus       int    `json:"pos_req_status"  config:"pos_req_status"   default:""`
	Account            string `json:"account"  config:"account"   default:""`
	AccountType        int    `json:"account_type"  config:"account_type"   default:""`
	Symbol             string `json:"symbol,omitempty"  config:"symbol"   default:""`
	Text               string `json:"text,omitempty"  config:"text"   default:""`
}

type PositionReport struct {
	PosMaintRptID        string `json:"pos_maint_rpt_id"`
	PosReqID             string `json:"pos_req_id,omitempty"`
	PosReqResult         int    `json:"pos_req_result"`
	ClearingBusinessDate string `json:"clearing_business_date"`
	Account              string `json:"account"`
	AccountType          int    `json:"account_type"`
	Symbol               string `json:"symbol,omitempty"`
	ContractMultiplier   int    `json:"contract_multiplier,omitempty"`
	SettlPrice           string `json:"settl_price"`
	SettlPriceType       int    `json:"settl_price_type"`
	PriorSetlPrice       string `json:"prior_setl_price"`
	LongPos              int    `json:"long_pos,omitempty"`
	ShortPos             int    `json:"short_pos,omitempty"`
}
