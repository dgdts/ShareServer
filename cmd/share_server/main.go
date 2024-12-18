package main

import (
	"context"
	"flag"
	"net/http"
	"path/filepath"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/dgdts/ShareServer/biz/share"
	global_init "github.com/dgdts/ShareServer/init"
	"github.com/dgdts/UniversalServer/pkg/config"
)

func main() {
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

	global_init.InitRedis(config.GetGlobalStaticConfig())

	global_init.InitMongo(config.GetGlobalStaticConfig())

	share.InitShareCache()

	server.Spin()
}
