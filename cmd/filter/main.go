package main

import (
	"webapi/common/word_filter"
	"github.com/mkideal/log"
	"flag"
	"webapi/config"
	"fmt"
)

var flConfigFilename = flag.String("c", "../../config/conf.cnf", "配置文件")

func main() {
	var conf *config.Config
	var err error
	if *flConfigFilename == "" {
		conf = config.NewConfig("filter")
	} else {
		conf, err = config.Load(*flConfigFilename, "filter")

		if err != nil {
			panic(err)
		}
	}
	defer log.Uninit(log.Init(conf.LogProvider, conf.LogOption))
	log.SetLevel(conf.LogLevel)

	options := word_filter.NewOptions(fmt.Sprintf(":%d", conf.FilterGrpcProt))
	filter := word_filter.New(options)
	filter.Run(options.TCPAddr)
}
