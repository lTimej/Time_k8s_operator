package httpResp

import (
	code2 "Time_k8s_operator/pkg/code"
	"net/http"
	"sync"
)

type respResult struct {
	Data interface{} `json:"data"`
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
}

type Response struct {
	HttpStatus int
	Result     respResult
}

var pool = sync.Pool{
	New: func() interface{} {
		return &Response{}
	},
}

func NewResponse(status int, code uint32, data interface{}) *Response {
	response := pool.Get().(*Response)
	response.HttpStatus = status
	response.Result.Data = data
	response.Result.Code = code
	response.Result.Msg = code2.GetMessage(code)
	return response
}

func ResponseOk(code uint32, data interface{}) *Response {
	return NewResponse(http.StatusOK, code, data)
}

func PutResponse(res *Response) {
	if res != nil {
		res.Result.Data = nil
		pool.Put(res)
	}
}
