package views

import (
	"github.com/feature/handler/cores"
	"github.com/feature/sdk/errno"
	"github.com/gin-gonic/gin"
)

func DemoHandler(c *gin.Context) {
	cores.NormalResponse(c, errno.OK, "hello world!")
}
