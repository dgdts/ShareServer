package main

import (
	"flag"
	"path/filepath"

	"github.com/dgdts/ShareServer/init"
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

	server := init.InitServer(config.GetGlobalStaticConfig())

}
