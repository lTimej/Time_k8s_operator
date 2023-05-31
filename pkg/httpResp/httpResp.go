package httpResp

import (
	code2 "Time_k8s_operator/pkg/code"
	"net/http"
	"sync"
	"github.com/gin-gonic/gin"
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

func NewResponseOk(status int, code uint32, data interface{}) *Response {
	response := pool.Get().(*Response)
	response.HttpStatus = status
	response.Result.Data = data
	response.Result.Code = code
	response.Result.Msg = code2.GetMessage(code)
	return response
}

func NewResponseNotOk(status int, code uint32, data interface{}, msg string) *Response {
	response := pool.Get().(*Response)
	response.HttpStatus = status
	response.Result.Data = data
	response.Result.Code = code
	response.Result.Msg = msg
	return response
}

func ResponseOk(code uint32, data interface{}) *Response {
	return NewResponseOk(http.StatusOK, code, data)
}

func RepsonseNotOk(msg string) *Response {
	return NewResponseNotOk(http.StatusOK, 30, nil, msg)
}

func PutResponse(res *Response) {
	if res != nil {
		res.Result.Data = nil
		pool.Put(res)
	}
}

func HttpResp(c *gin.Context,code uint32,data interface{}){
	r := NewResponseOk(http.StatusOK, code, data)
	c.JSON(r.HttpStatus, &r.Result)
}