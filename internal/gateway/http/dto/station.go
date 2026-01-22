package dto

// StationSuggestHTTPReq 站点联想/搜索 HTTP 请求参数（通过 URL Query 传参）。
//
// 对应接口：GET /api/v1/station/suggest
type StationSuggestHTTPReq struct {
	Keyword string `query:"keyword,required"` // 关键字（站点名/拼音/编码片段）
	Limit   int    `query:"limit"`            // 返回条数（建议 1~50；不传由后端设默认）
}

// StationSuggestHTTPResp 站点联想/搜索 HTTP 响应。
type StationSuggestHTTPResp struct {
	Code  int32    `json:"code"`  // 业务状态码
	Msg   string   `json:"msg"`   // 状态说明/错误信息
	Items []string `json:"items"` // 候选站点列表（通常为站点名或“名(码)”的展示串）
}
