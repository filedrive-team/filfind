package accesslog

import (
	"context"
	"github.com/filedrive-team/filfind/backend/models"
	"github.com/filedrive-team/filfind/backend/repo"
	"github.com/filedrive-team/filfind/backend/utils/utils"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	logger "github.com/sirupsen/logrus"
	"time"
)

func AccessLog(ctx context.Context, repo *repo.Manager) gin.HandlerFunc {
	logCh := make(chan *models.AccessLog, 1024)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(logCh)
				return
			case log, ok := <-logCh:
				if !ok {
					return
				}

				n := len(logCh)
				logs := make([]*models.AccessLog, 0, n+1)
				logs = append(logs, log)
				for i := 0; i < n; i++ {
					logs = append(logs, <-logCh)
				}
				err := repo.CreateAccessLog(logs)
				if err != nil {
					logger.WithError(err).Error("call CreateAccessLog failed")
				}
			}
		}
	}()

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		// filter out the unknown path and the swagger uri
		fullPath := c.FullPath()
		if fullPath == "" || fullPath == "/swagger/*any" {
			return
		}

		ua := user_agent.New(c.GetHeader("User-Agent"))
		browser, browserVersion := ua.Browser()
		log := &models.AccessLog{
			Method:         c.Request.Method,
			Uri:            c.Request.RequestURI,
			FullPath:       fullPath,
			Cost:           time.Since(start),
			Ip:             utils.GetClientIP(c),
			Browser:        browser,
			BrowserVersion: browserVersion,
			OS:             ua.OS(),
			Platform:       ua.Platform(),
			Mobile:         ua.Mobile(),
			Bot:            ua.Bot(),
			Status:         c.Writer.Status(),
		}
		logCh <- log
	}
}
