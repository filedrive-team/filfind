package repo

import (
	"github.com/filedrive-team/filfind/backend/api/token"
	"github.com/filedrive-team/filfind/backend/kvcache"
	"github.com/filedrive-team/filfind/backend/models"
	"github.com/filedrive-team/filfind/backend/settings"
	"github.com/filedrive-team/filfind/backend/utils/jwttoken"
	logger "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"os"
	"time"
)

type Manager struct {
	db          *gorm.DB
	tkGen       *jwttoken.TokenGenerator
	publishDate time.Time
}

func NewManage(conf *settings.AppConfig, logWriter io.Writer) *Manager {
	m := &Manager{
		tkGen:       token.NewTokenGenerator(conf.App.JwtSecret),
		publishDate: conf.App.PublishDate,
	}
	m.setup(conf, logWriter)

	kvcache.InitGlobalCache(conf)

	return m
}

type debugLogWriter struct {
	logWriter io.Writer
}

func (w *debugLogWriter) Write(p []byte) (n int, err error) {
	os.Stdout.Write(p)
	return w.logWriter.Write(p)
}

// setup initializes the database instance
func (m *Manager) setup(conf *settings.AppConfig, logWriter io.Writer) {
	var err error
	var dialector gorm.Dialector
	switch conf.Database.Type {
	case "postgres":
		dialector = postgres.Open(conf.Database.Dsn)
	case "mysql":
		dialector = mysql.Open(conf.Database.Dsn)
	}

	writer := logWriter
	logLevel := gormlogger.Warn
	if conf.App.Debug {
		writer = &debugLogWriter{logWriter: logWriter}
	}
	if logger.GetLevel() >= logger.DebugLevel {
		logLevel = gormlogger.Info
	}

	newLogger := gormlogger.New(
		log.New(writer, "[gorm] ", log.LstdFlags), // io writer
		gormlogger.Config{
			SlowThreshold: 5 * time.Second, // slow SQL threshold
			LogLevel:      logLevel,        // Log level
			Colorful:      false,           // disable color printing
		},
	)
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   conf.Database.TablePrefix,
			SingularTable: true,
		},
	})

	if err != nil {
		logger.Fatalf("database setup err: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatalf("db.DB() err: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = db.AutoMigrate(models.Models()...)
	if err != nil {
		logger.Fatalf("db.AutoMigrate() err: %v", err)
	}
	m.db = db

	logger.Info("database setup finished!")
}
