package handler

import (
	"github.com/feature/handler/views"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func setupRoute(engine *gin.Engine) {
	applyApiRoutes(engine)
	engine.GET("/_health_check", func(c *gin.Context) {
		c.String(200, "ok")
	})
	pprof.Register(engine)
}

func applyApiRoutes(engine *gin.Engine) {
	group := engine.Group("/api/v1")
	{
		registerTest(group)
	}

}

func registerTest(rg *gin.RouterGroup) {
	test := rg.Group("/test")
	{
		test.GET("/demo", views.DemoHandler)
	}
}
