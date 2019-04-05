package word_filter

import (
	log "github.com/sirupsen/logrus"
)

type Options struct {
	//自定义脏词库
	UserDictDataPath string
	//TCP监听地址
	TCPAddr string
}

func NewOptions(addr string) *Options {
	return &Options{
		UserDictDataPath:    "./user_dict.txt",
		TCPAddr:             addr,
	}
}

func init() {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
}
