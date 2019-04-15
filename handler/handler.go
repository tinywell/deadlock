package handler

import (
	"time"

	"go.uber.org/zap"

	"deadlock/hlserver"
)

type ConsumeHandler struct {
}

func (h *ConsumeHandler) HandlerMsg(recvChan <-chan hlserver.RecvData) {
	for data := range recvChan {
		go func(data hlserver.RecvData) {
			// rsp := h.Handle(data.Data)
			zap.L().Named("handler").Info("Handle Message Start")
			defer zap.L().Named("handler").Info("Handle Message End")

			zap.L().Named("handler").Info("Handle PutRatio", zap.String("messageid", data.Data.MessageID), zap.String("trancode", data.Data.TranCode))
			defer zap.L().Named("handler").Info("Handle PutRatio End")

			rsp := &hlserver.RspData{MessageID: "", Code: 200}
			rsp.Data.TxID = time.Now().UTC().String()
			rsp.Data.Message = "测试"
			time.Sleep(time.Microsecond * 100)

			zap.L().Named("handler").Debug("HandleMsg completed", zap.Any("result", rsp))
			data.RspChan <- *rsp
			close(data.RspChan)
		}(data)
	}
}
