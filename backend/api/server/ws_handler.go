package server

import (
	"github.com/filedrive-team/filfind/backend/api/ws"
	"github.com/gin-gonic/gin"
)

// ws godoc
// @Summary websocket api
// @Tags ws
// @Accept  json
// @Produce  json
// @Success 200
// @Router /ws [get]
func (s *Server) ws(c *gin.Context) {
	ws.ServeWs(s.hub, c.Writer, c.Request)
}
