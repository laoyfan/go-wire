package constant

import "time"

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type ErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
type LogLayout struct {
	Time      time.Time `json:"time"`                // 请求的时间
	Status    int       `json:"status,omitempty"`    // HTTP响应状态码
	Method    string    `json:"method,omitempty"`    // HTTP方法
	Path      string    `json:"path,omitempty"`      // 请求的路径
	Query     string    `json:"query,omitempty"`     // 请求的Query参数
	IP        string    `json:"IP,omitempty"`        // 客户端IP
	UserAgent string    `json:"userAgent,omitempty"` // 用户代理
	Error     string    `json:"error,omitempty"`     // 错误信息
	Cost      float64   `json:"cost,omitempty"`      // 请求耗时
	Source    string    `json:"source,omitempty"`    // 请求来源
}
