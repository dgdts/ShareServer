package main

import (
	"context"
	"flag"
	"net/http"
	"path/filepath"

	"github.com/cloudwego/hertz/pkg/app"
	global_init "github.com/dgdts/ShareServer/init"
	"github.com/dgdts/UniversalServer/pkg/config"
)

func main() {
	// 1. read and parse config
	configFilePath := flag.String("config", "../../conf/dev/conf.yaml", "config file path")
	absPath, err := filepath.Abs(*configFilePath)
	if err != nil {
		panic(err)
	}

	err = config.InitConfigFromLocal(absPath)
	if err != nil {
		panic(err)
	}

	global_init.InitLogger(config.GetGlobalStaticConfig().Log)

	server := global_init.InitServer(config.GetGlobalStaticConfig())

	server.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		c.String(http.StatusOK, "pong")
	})

	server.Spin()
}
