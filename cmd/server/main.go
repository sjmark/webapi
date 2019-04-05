package main

import (
	"os"
	"flag"
	. "webapi/cmd/server/handler"
	"webapi/config"
	"github.com/mkideal/log"
	"webapi/common/osutil/signal"
	"webapi/business/api_loginc"
	"webapi/models"
)

var flConfigFilename = flag.String("c", "../../config/conf.cnf", "配置文件")

func main() {
	var conf *config.Config
	var err error
	if *flConfigFilename == "" {
		conf = config.NewConfig("server")
	} else {
		conf, err = config.Load(*flConfigFilename, "server")

		if err != nil {
			panic(err)
		}
	}

	defer log.Uninit(log.Init(conf.LogProvider, conf.LogOption))
	log.SetLevel(conf.LogLevel)
	// 初始化hanndler
	Init(conf)
	// 初始化业务层
	api_loginc.Init(conf)
	// 初始化models
	models.InitXormSession(conf.WriteDB.NewSession(), conf.ReadDB.NewSession())
	// 启动服务器
	go run(conf)
	// 脏词过滤
	go grpcRun(conf)
	signal.Wait(os.Interrupt)
}
