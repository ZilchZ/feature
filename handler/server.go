package handler

import (
	"context"
	"fmt"
	"github.com/feature/conf"
	"github.com/feature/handler/middlewares"
	"github.com/feature/sdk/log"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var (
	cfg    = conf.Config
	server *http.Server
)

type Server struct {
	server *http.Server
}

func (s *Server) Init() {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery(), requestid.New(), serverLog())
	setupRoute(engine)
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	s.server = &http.Server{
		Addr:    addr,
		Handler: engine,
	}
}

func (s *Server) Launch() {
	log.Logger.Infof("Server start at %s", s.server.Addr)
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			log.Logger.Error(err)
		}
	}()
}

func (s *Server) Stop() error {
	if s.server != nil {
		log.Logger.Warningf("Server shutdown now with timeout %d s.", cfg.Server.ShutdownTimeout)
		timeout := time.Duration(cfg.Server.ShutdownTimeout) * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := s.server.Shutdown(ctx); err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			log.Logger.Warning("Server shutdown timeout.")
		default:
			log.Logger.Warning("Server exited")
		}
	}
	return nil
}

func serverLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		middlewares.LogAccess(start, c)
	}
	
}