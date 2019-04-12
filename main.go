package main

import (
	"context"
	"flag"
	"fmt"

	"go.uber.org/zap"

	"net/http"
	_ "net/http/pprof"

	"deadlock/handler"
	hs "deadlock/hlserver"
	"deadlock/hlserver/restful"
	"deadlock/mockhl"
	log "deadlock/zaplog"
)

const (
	// DefaultRecvCache 消息接收通道缓冲
	DefaultRecvCache = 100
)

var logger *zap.Logger

func main() {

	var configPath string
	var fconfigPath string
	var help bool

	// 配置文件路径参数
	flag.BoolVar(&help, "help", false, "print help message")

	flag.Parse()

	// 展示命令使用帮助信息
	if help {
		flag.Usage()
	}

	go func() {
		fmt.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	// 如果没有独立的sdk配置文件，则默认相关配置合并在应用配置中
	if len(fconfigPath) == 0 {
		fconfigPath = configPath
	}

	// 初始化日志模块
	log.InitLog(log.LogConfig{
		Level: "debug",
	})
	logger = log.MustGetLogger("main")

	hlf := &mockhl.MockHlfront{}
	logger.Info("Put Public Key Successful !!!")
	handler := handler.NewHandler(hlf, context.Background())

	restServer := restful.NewServer("0.0.0.0:8000")

	// 注册事件监听
	logger.Info("Starting Event Server ...")

	// 启动接口数据接收及处理

	recvChan := make(chan hs.RecvData, DefaultRecvCache)
	logger.Info("Handler cache length", zap.Int("length", cap(recvChan)))
	logger.Info("Starting HLServer Server ...")

	err := restServer.Start(recvChan)
	if err != nil {
		logger.Error("start hlserver error", zap.Error(err))
		panic(err)
	}

	logger.Info("Starting Message Handler ...")
	handler.HandlerMsg(recvChan)
}
