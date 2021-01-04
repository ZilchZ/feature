package httputil

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type ResponseJSON struct {
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"success"`
	Object  interface{} `json:"data"`
}

func (r *ResponseJSON) String() string {
	return fmt.Sprintf("code:%d message:%s data:%v", r.Code, r.Message, r.Object)
}

func NormalResponseDetail(c *gin.Context, code int, message string, object interface{}) {
	var resp ResponseJSON
	switch obj := object.(type) {
	case error:
		resp = ResponseJSON{code, message, obj.Error()}
	default:
		resp = ResponseJSON{code, message, object}

	}
	c.JSON(200, resp)

}
