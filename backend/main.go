package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/filedrive-team/filfind/backend/api/server"
	"github.com/filedrive-team/filfind/backend/settings"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var (
	showV      bool
	configFile string
	loglevel   string
	initSystem bool
)

func printVersion() {
	fmt.Printf("filfind-backend version: v%s.%s.%s-%s\ngithub.com/gin-gonic/gin version: %s\n",
		Major, Minjor,
		Patch, BuildVersion, gin.Version)
}

func main() {
	flag.BoolVar(&showV, "version", false, "print version")
	flag.StringVar(&configFile, "config", "conf/app.toml", "set config file path")
	flag.StringVar(&loglevel, "loglevel", "debug", "set log level")
	flag.BoolVar(&initSystem, "init", false, "init system")
	flag.Parse()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(fmt.Errorf("fetch work dir failed %+v", err))
	}

	if showV {
		printVersion()
		os.Exit(0)
	}

	conf := settings.LoadConfig(configFile)

	mainCtx, mainCancel := context.WithCancel(context.Background())
	srv := server.NewServer(conf, dir, loglevel, initSystem)
	if initSystem {
		srv.Run(mainCtx)
		os.Exit(0)
	}
	go srv.Run(mainCtx)

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown Server ...")

	mainCancel()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown: ", err)
	}

	logger.Info("Server exit")
}
