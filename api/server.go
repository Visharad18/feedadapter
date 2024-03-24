package api

import (
	"fmt"
	"net/http"

	"github.com/Visharad18/feedadapter/app"
	"github.com/Visharad18/feedadapter/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg    *config.Config
	engine *gin.Engine
	app    *app.App
}

func NewServer(cfg *config.Config, app *app.App) (*Server, error) {
	server := &Server{cfg: cfg, app: app}
	server.engine = gin.New()
	server.engine.GET("/get", server.get)
	return server, server.engine.Run(fmt.Sprintf(":%s", server.cfg.HTTPPort))
}

func (s *Server) get(c *gin.Context) {
	data := s.app.GetData()
	c.JSON(http.StatusAccepted, data)
}
