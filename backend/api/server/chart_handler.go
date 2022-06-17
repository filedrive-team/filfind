package server

import (
	"github.com/filedrive-team/filfind/backend/api"
	"github.com/filedrive-team/filfind/backend/api/errormsg"
	"github.com/filedrive-team/filfind/backend/api/ws"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
)

// chatHistory godoc
// @Summary Get chat history
// @Tags chat
// @Accept  json
// @Produce  json
// @Param Authorization header string true "jwt access token" default(Bearer YOUR_JWT)
// @Param partner query string false "Partner uid"
// @Param limit query int false "Max records to return" mininum(1) maximum(100)
// @Param before query int false "Option, unix time(seconds). It returns last messages if you not set."
// @Success 200
// @Router /chat/history [get]
func (s *Server) chatHistory(c *gin.Context) {
	params := new(ws.HistoryMessageReq)
	err := c.ShouldBindQuery(params)
	if err != nil {
		errMsg := errormsg.ByCtx(c, errormsg.ParamsError)
		api.JSONError(c, errMsg)
		logger.Error(err)
		return
	}
	tk := MustGetToken(c)

	list, err := s.repo.MessagesByUser(tk.Uid, params.Partner, params.Before, params.Limit)
	if err != nil {
		logger.WithError(err).Error("call MessagesByUser failed")
		return
	}

	results := make([]*ws.Message, 0, len(list))
	for _, msg := range list {
		results = append(results, &ws.Message{
			Sender:    msg.Sender,
			Recipient: msg.Recipient,
			Type:      msg.Type,
			Content:   msg.Content,
			Timestamp: &msg.CreatedAt,
			Checked:   msg.Checked,
		})
	}
	api.JSONOk(c, results)
}
