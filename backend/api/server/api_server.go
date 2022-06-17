package server

import (
	"context"
	"fmt"
	"github.com/filedrive-team/filfind/backend/api/middleware/accesslog"
	"github.com/filedrive-team/filfind/backend/api/middleware/auth"
	"github.com/filedrive-team/filfind/backend/api/middleware/cors"
	"github.com/filedrive-team/filfind/backend/api/ws"
	"github.com/filedrive-team/filfind/backend/docs"
	"github.com/filedrive-team/filfind/backend/filclient"
	"github.com/filedrive-team/filfind/backend/jobs"
	"github.com/filedrive-team/filfind/backend/log"
	"github.com/filedrive-team/filfind/backend/models"
	"github.com/filedrive-team/filfind/backend/repo"
	"github.com/filedrive-team/filfind/backend/settings"
	"github.com/filedrive-team/filfind/backend/smtp"
	"github.com/filedrive-team/filfind/backend/utils/jwttoken"
	"github.com/filedrive-team/filfind/backend/utils/utils"
	"github.com/filedrive-team/filfind/backend/validator"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type Server struct {
	httpSrv    *http.Server
	repo       *repo.Manager
	filClient  *filclient.FilClient
	hub        *ws.Hub
	conf       *settings.AppConfig
	dir        string
	initSystem bool
}

func NewServer(conf *settings.AppConfig, dir string, loglevel string, initSystem bool) *Server {
	logWriter, err := log.InitLogger(dir, loglevel)
	if err != nil {
		panic(fmt.Errorf("init logger failed %+v", err))
	}
	r := repo.NewManage(conf, logWriter)
	s := &Server{
		repo:       r,
		filClient:  filclient.NewFileClient(conf.App.FilecoinApi),
		hub:        ws.NewHub(r),
		conf:       conf,
		dir:        dir,
		initSystem: initSystem,
	}
	return s
}

func (s *Server) registerRouter(ctx context.Context) http.Handler {
	conf := s.conf
	if logger.GetLevel() < logger.DebugLevel {
		println("gin mode: ReleaseMode")
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.AddCorsHeaders())
	r.Use(accesslog.AccessLog(ctx, s.repo))
	r.Use(func(c *gin.Context) {
		// Do not use gzip compression for static resources
		if strings.HasPrefix(c.FullPath(), "/avatars") {
			c.Next()
			return
		}
		handler := gzip.Gzip(gzip.DefaultCompression)
		handler(c)
	})

	if conf.App.Swag {
		// programatically set swagger info
		docs.SwaggerInfo.Title = settings.ProductName + " backend API"
		//docs.SwaggerInfo.Version = settings.Version
		docs.SwaggerInfo.BasePath = "/api/v1"
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	cacheDir := filepath.Join(s.dir, "cache")
	r.Static("/avatars", cacheDir)

	baseGroup := r.Group("/api/v1")
	baseGroup.POST("/userSignUp", s.userSignUp)
	baseGroup.POST("/userLogin", s.userLogin)
	baseGroup.POST("/userResetPwd", s.userResetPwd)
	baseGroup.POST("/vcodeByEmailToResetPwd", s.vcodeByEmailToResetPwd)
	baseGroup.GET("/providers", s.providerList)
	baseGroup.GET("/spOwnerProfile", s.spOwnerProfile)
	baseGroup.GET("/spServiceDetail", s.spServiceDetail)
	baseGroup.GET("/spOwnerReviews", s.spOwnerReviews)
	baseGroup.GET("/clientProfile", s.clientProfile)
	baseGroup.GET("/clientDetail", s.clientDetail)
	baseGroup.GET("/clientHistoryDealStats", s.clientHistoryDealStats)
	baseGroup.GET("/clientReviews", s.clientReviews)
	baseGroup.GET("/clients", s.clientList)

	userGroup := baseGroup.Group("/user", auth.CheckAuthority(s.repo.GetTokenVerify()))
	{
		userGroup.POST("/modifyPassword", s.modifyPassword)
		userGroup.POST("/token", s.Token)
		userGroup.POST("/profile", s.modifyProfile)
	}

	clientGroup := baseGroup.Group("/client", auth.CheckAuthority(s.repo.GetTokenVerify(), models.ClientRole))
	{
		clientGroup.POST("/detail", s.modifyClientDetail)
		clientGroup.POST("/review", s.submitReview)
	}

	providerGroup := baseGroup.Group("/provider", auth.CheckAuthority(s.repo.GetTokenVerify(), models.SPOwnerRole))
	{
		providerGroup.POST("/detail", s.modifyProviderDetail)
	}

	chatGroup := baseGroup.Group("/chat", auth.CheckAuthority(s.repo.GetTokenVerify()))
	{
		chatGroup.GET("/history", s.chatHistory)
	}

	// websocket entrypoint
	baseGroup.GET("/ws", s.ws)

	adminGroup := baseGroup.Group("/admin")
	{
		adminGroup.POST("/userLogin", s.adminUserLogin)
	}
	adminUserGroup := adminGroup.Group("/user", auth.CheckAuthority(s.repo.GetTokenVerify()))
	{
		adminUserGroup.POST("/modifyPassword", s.modifyAdminPassword)
	}

	metricsGroup := adminGroup.Group("/metrics", auth.CheckAuthority(s.repo.GetTokenVerify()))
	{
		metricsGroup.GET("/overview", s.metricsOverview)
		metricsGroup.GET("/spOverview", s.metricsSpOverview)
		metricsGroup.GET("/clientOverview", s.metricsClientOverview)
		metricsGroup.GET("/spToDealNewClientDetail", s.metricsSpToDealNewClientDetail)
		metricsGroup.GET("/clientToDealSpDetail", s.metricsClientToDealSpDetail)
	}

	return r
}

func (s *Server) Run(ctx context.Context) {
	// init system
	if s.initSystem {
		// create admin user
		exist, err := s.repo.ExistAdminUserByName(settings.DefaultAdminUser)
		if err != nil {
			logger.WithError(err).Error("call ExistAdminUserByName failed")
			return
		}
		if !exist {
			hashedPassword, err := utils.GenerateHashedPassword(settings.DefaultAdminPassword + s.conf.App.PasswordSalt)
			if err != nil {
				logger.WithError(err).Error("call GenerateHashedPassword failed")
				return
			}
			admin := &models.AdminUser{
				Name:           settings.DefaultAdminUser,
				HashedPassword: hashedPassword,
			}
			err = s.repo.CreateAdminUser(admin)
			if err != nil {
				return
			}
		}

		// create system user to send system message
		exist, err = s.repo.ExistUserByEmail(settings.SystemUser)
		if err != nil {
			logger.WithError(err).Error("call ExistUserByEmail failed")
			return
		}
		if !exist {
			system := &models.User{
				Name:          settings.SystemUser,
				Email:         settings.SystemUser,
				Type:          models.SystemRole,
				AddressRobust: settings.SystemUser,
			}
			err = s.repo.CreateUser(system)
			if err != nil {
				logger.WithError(err).Error("call CreateUser failed")
				return
			}
		}

		crawler := jobs.NewCrawler(s.conf.App.FilrepApi, s.repo)
		crawler.Init(ctx)

		syncer := jobs.NewFilecoinSyncer(s.filClient, s.repo)
		syncer.Init(ctx)
		return
	}

	crawler := jobs.NewCrawler(s.conf.App.FilrepApi, s.repo)
	go crawler.Run(ctx)

	syncer := jobs.NewFilecoinSyncer(s.filClient, s.repo)
	go syncer.Run(ctx)

	go s.hub.Run(ctx)

	smtp.Setup(s.conf)

	validator.InitExtendValidation()

	router := s.registerRouter(ctx)

	addr := fmt.Sprintf(":%d", s.conf.Server.HttpPort)
	maxHeaderBytes := 1 << 20

	s.httpSrv = &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    time.Duration(s.conf.Server.ReadTimeout * 1e9),
		WriteTimeout:   time.Duration(s.conf.Server.WriteTimeout * 1e9),
		MaxHeaderBytes: maxHeaderBytes,
	}

	// service connections
	if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("listen: %s\n", err)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.initSystem {
		return nil
	}
	return s.httpSrv.Shutdown(ctx)
}

func MustGetToken(c *gin.Context) *jwttoken.JwtPayload {
	v, ok := c.Get(settings.TokenKey)
	if ok {
		return v.(*jwttoken.JwtPayload)
	}
	logger.Fatal("can't get user token from gin.Context")
	return nil
}
