package finalto

type FixApiRequest struct {
	Event  string      `json:"event"  config:"event"   default:""` // 请求时间类型
	Params interface{} `json:"params,omitempty" config:"params"`   // 具体请求参数
}
