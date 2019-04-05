package config

import (
	"fmt"

	"webapi/common/redis"

	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/config"
	"github.com/go-xorm/xorm"
	"github.com/mkideal/log/logger"
	"webapi/common/tools/tool"
	"google.golang.org/grpc"
)

var rc4Pwd = "Rc4PwdSalt..xdcyyouyuanbao.."
var pwdSalt = "xdcy"

type Config struct {
	//redis 切片
	RedisR1    redis.Cache
	RedisR2    redis.Cache
	ReadDB     *xorm.Engine
	WriteDB    *xorm.Engine
	FilterConn *grpc.ClientConn

	// 服务器 ID
	ServerId int
	// 服务器监听的 TCP 端口
	ImPort         int
	FilterGrpcProt int
	ServerPort     int
	// 每个连接同时能写的消息数量
	ConWriteSize int

	// 日志配置
	LogLevel    logger.Level
	LogProvider string
	LogOption   string
	QueueName   string
	PwdSalt     string
	Rc4PwdSalt  string
}

func NewConfig(logName string) *Config {
	// 创建默认配置
	return &Config{
		ImPort:         8282,
		ServerPort:     8287,
		FilterGrpcProt: 8289,
		ConWriteSize:   32,
		LogLevel:       logger.TRACE,
		LogProvider:    "console/file",                                                 // 默认日志输出到控制台和文件,生产环境只可以使用文件 file
		LogOption:      fmt.Sprintf("dir=../%s&filename=handler&suffix=.txt", logName), // 指定日志输出文件
		PwdSalt:        pwdSalt,
		Rc4PwdSalt:     tool.MD5f(tool.MD5f(rc4Pwd)),
	}
}

func initEngine(cfg *config.Config, d string) (*xorm.Engine, error) {
	driver, _ := cfg.String(d, "db.driver")
	dbname, _ := cfg.String(d, "db.dbname")
	user, _ := cfg.String(d, "db.user")
	pwd, _ := cfg.String(d, "db.password")
	host, _ := cfg.String(d, "db.host")
	sourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true", user, pwd, host, dbname)
	return xorm.NewEngine(driver, sourceName)
}

func initRedis(cfg *config.Config, d int) *redis.Cache {
	host, _ := cfg.String("redis", fmt.Sprintf("host-master%d", d))
	pwd, _ := cfg.String("redis", fmt.Sprintf("pwd-master%d", d))
	var redis = redis.NewRedisCache("1", host, pwd, 2000)
	return &redis
}

func Load(path, logName string) (*Config, error) {

	conf := new(Config)
	cfg, err := config.ReadDefault(path)
	if err != nil {
		return nil, err
	}
	//日志
	isprot, _ := cfg.Int("web", "isprot")
	conf.LogProvider = "console/file"
	conf.LogOption = fmt.Sprintf("dir=../../%s&filename=handler&suffix=.txt", logName)
	// 1 表示线上line 表示dev
	if isprot == 1 {
		conf.LogProvider = "file"
		conf.LogOption = fmt.Sprintf("dir=../%s&filename=handler&suffix=.txt", logName)
	}
	level, _ := cfg.Int("log", "log-level")
	conf.LogLevel = logLevel(level)

	//服务器
	conf.ImPort, _ = cfg.Int("web", "im_port")
	conf.ServerPort, _ = cfg.Int("web", "server_port")
	conf.FilterGrpcProt, _ = cfg.Int("web", "filter_port")
	conf.ConWriteSize, _ = cfg.Int("web", "conwritesize")

	//redis tcp 连接 redis 分片
	if v := initRedis(cfg, 1); v != nil {
		conf.RedisR1 = *v
	}
	if v := initRedis(cfg, 2); v != nil {
		conf.RedisR2 = *v
	}

	// mysql tcp 连接 读写分离
	conf.WriteDB, err = initEngine(cfg, "database")

	if err != nil {
		return nil, err
	}

	conf.ReadDB, err = initEngine(cfg, "database.readonly")
	if err != nil {
		return nil, err
	}

	conf.PwdSalt = pwdSalt
	conf.Rc4PwdSalt = tool.MD5f(tool.MD5f(rc4Pwd))
	//队列初始化

	return conf, nil
}

func logLevel(level int) logger.Level {
	switch level {
	case 0:
		return logger.FATAL
	case 1:
		return logger.ERROR
	case 2:
		return logger.WARN
	case 3:
		return logger.INFO
	case 4:
		return logger.DEBUG
	case 5:
		return logger.TRACE
	case 6:
		return logger.NumLevel
	default:
		return logger.INFO
	}
}
