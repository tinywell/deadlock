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
	// ServiceConfigPath 客户端配置文件
	ServiceConfigPath = "./config/config.yaml"

	// DefaultRecvCache 消息接收通道缓冲
	DefaultRecvCache = 100
)

var logger *zap.Logger

func main() {

	var configPath string
	var fconfigPath string
	var help, pprof bool

	// 配置文件路径参数
	flag.BoolVar(&help, "help", false, "print help message")
	flag.BoolVar(&pprof, "pprof", false, "pprof enable")
	// flag.StringVar(&configPath, "config", ServiceConfigPath,
	// 	"config file path (Default use relative path '"+ServiceConfigPath+"')")
	// flag.StringVar(&fconfigPath, "fconfig", "",
	// 	"sdk config file path")
	flag.Parse()

	// 展示命令使用帮助信息
	if help {
		flag.Usage()
	}

	// if pprof {
	// 	fmt.Println("pprof enable")
	go func() {
		fmt.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()
	// }

	// 如果没有独立的sdk配置文件，则默认相关配置合并在应用配置中
	if len(fconfigPath) == 0 {
		fconfigPath = configPath
	}

	// 初始化日志模块
	log.InitLog(log.LogConfig{
		Level: "debug",
	})
	logger = log.MustGetLogger("main")

	// // new底层服务模块（fabric-sdk、加解密等）
	// service := sdk.NewService(fconfigPath)
	// hlf := hlfront.NewHLF(fconfigPath, cfg, service, service)

	// // 上传当前组织加解密公钥
	// err := hlf.PutPublicKey()
	// if err != nil {
	// 	logger.Errorf("Put public key error:%s", err.Error())
	// 	os.Exit(1)
	// }

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

	// runMonitor(cfg)

	logger.Info("Starting Message Handler ...")
	handler.HandlerMsg(recvChan)
}
