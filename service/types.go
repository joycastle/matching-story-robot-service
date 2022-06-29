package service

type MatchingRPCResponse struct {
	Code   int    `json:"code"`
	Errmsg string `json:"errmsg"`
	Data   string `json:"data"`
}
