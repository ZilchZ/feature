package cores

import (
	"github.com/feature/sdk/errno"
	"github.com/feature/sdk/httputil"
	"github.com/gin-gonic/gin"
)

func NormalResponse(c *gin.Context, err *errno.Errno, object interface{}) {
	httputil.NormalResponseDetail(c, err.Code, err.Message, object)
}
