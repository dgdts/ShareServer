package init

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	hertz_config "github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/dgdts/ShareServer/biz/router"
	"github.com/dgdts/UniversalServer/pkg/config"
)

func InitServer(config *config.GlobalConfig) *server.Hertz {
	if len(config.Hertz.Service) == 0 {
		hlog.Error("service address is empty")
		return nil
	}

	serverOptions := []hertz_config.Option{
		server.WithHostPorts(config.Hertz.Service[0].Address),
	}

	if config.Hertz.EnablePprof {
		serverOptions = append(serverOptions, server.WithHandleMethodNotAllowed(true))
	}

	s := server.New(serverOptions...)

	router.GeneratedRegister(s)

	return s
}
