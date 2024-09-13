package apiutil

// Response 标准返回体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type EmptyResponse struct{}
