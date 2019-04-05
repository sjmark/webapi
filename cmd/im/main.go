package main

import (
	"os"
	"flag"
	_ "webapi/cmd/im/handler"
	"webapi/config"
	"github.com/mkideal/log"
	"webapi/common/osutil/signal"
)

var flConfigFilename = flag.String("c", "../../config/conf.cnf", "配置文件")

func main() {
	var conf *config.Config
	var err error
	if *flConfigFilename == "" {
		conf = config.NewConfig("im_server")
	} else {
		conf, err = config.Load(*flConfigFilename, "im_server")
		if err != nil {
			panic(err)
		}
	}

	defer log.Uninit(log.Init(conf.LogProvider, conf.LogOption))
	log.SetLevel(conf.LogLevel)
	// 启动服务器
	go startup(conf)

	signal.Wait(os.Interrupt)
}
