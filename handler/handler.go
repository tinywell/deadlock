package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"deadlock/hlserver"
	"deadlock/mockhl"
)

type ConsumeHandler struct {
	hlf      *mockhl.MockHlfront
	ctx      context.Context
	sendChan chan<- *hlserver.EventData
}

func NewHandler(hlf *mockhl.MockHlfront, ctx context.Context) *ConsumeHandler {
	return &ConsumeHandler{
		hlf: hlf,
		ctx: ctx,
	}
}

func (h *ConsumeHandler) HandlerMsg(recvChan <-chan hlserver.RecvData) {
	for data := range recvChan {
		go func(data hlserver.RecvData) {
			rsp := h.Handle(data.Data)
			// logger.Debugf("HandleMsg result:%+v", rsp)
			data.RspChan <- *rsp
			close(data.RspChan)
		}(data)
	}
}

func (h *ConsumeHandler) Handle(data hlserver.QIData) *hlserver.RspData {
	zap.L().Named("handler").Info("Handle Message Start")
	defer zap.L().Named("handler").Info("Handle Message End")

	zap.L().Named("handler").Info("Handle PutRatio", zap.String("messageid", data.MessageID), zap.String("trancode", data.TranCode))
	defer zap.L().Named("handler").Info("Handle PutRatio End")

	msg, err := json.Marshal(data.TranData)
	if err != nil {
		return hlserver.CommonRspWithError(fmt.Errorf("marshal trandata error: %s", err.Error()))
	}
	return h.hlf.TransactionInvoke("mychannel", "kyc", "putratio",
		string(msg), data.Crypto, data.Members)

	zap.L().Named("handler").Warn("trancode not supported", zap.String("trancode", data.TranCode))
	return hlserver.CommonRspWithError(fmt.Errorf("error trancode :%s", data.TranCode))
}
