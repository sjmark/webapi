package main

import (
	"time"
	"net/http"

	"webapi/cmd/server/handler"
	"webapi/common/tools/tool"
	"webapi/config"
	"fmt"
	"google.golang.org/grpc"
)

func run(cfg *config.Config) {
	defer tool.PrintPanicStack("startup")

	mux := http.NewServeMux()
	mux.HandleFunc("/index", handler.HandlerCore)
	server := http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func grpcRun(cfg *config.Config) {
	defer tool.PrintPanicStack("grpcRun")
	// 连接
	conn, err := grpc.Dial("127.0.0.1:8289", grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	// 初始化客户端
	cfg.FilterConn = conn
}
