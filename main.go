package main

import (
	"deadlock/zaplog"
	"fmt"

	"go.uber.org/zap"

	"net/http"
	_ "net/http/pprof"

	"deadlock/handler"
	hs "deadlock/hlserver"
	"deadlock/hlserver/restful"
)

const (
	// DefaultRecvCache 消息接收通道缓冲
	DefaultRecvCache = 200
)

var logger *zap.Logger

func main() {

	zaplog.InitLog()
	logger = zap.L()

	go func() {
		fmt.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	recvChan := make(chan hs.RecvData, DefaultRecvCache)
	logger.Info("Handler cache length", zap.Int("length", cap(recvChan)))

	handler := &handler.ConsumeHandler{}
	restServer := restful.NewServer("0.0.0.0:8000")

	logger.Info("Starting Event Server ...")
	logger.Info("Starting HLServer Server ...")
	_ = restServer.Start(recvChan)

	logger.Info("Starting Message Handler ...")
	handler.HandlerMsg(recvChan)
}
