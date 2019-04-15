package restful

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	"deadlock/hlserver"
)

// var logger = zaplog.MustGetLogger("restful")

const (
	URLSaveData         = "/agent/transaction"
	WaitResponseTimeOut = time.Second * 60
)

// RestServer RESTful接口
type RestServer struct {
	addr     string             // 端口
	router   *httprouter.Router // 路由
	recvChan chan<- hlserver.RecvData
}

func NewServer(addr string) *RestServer {
	return &RestServer{
		addr:   addr,
		router: httprouter.New(),
	}
}

func (server *RestServer) Start(recvchan chan<- hlserver.RecvData) error {
	server.recvChan = recvchan
	server.router.POST(URLSaveData, server.SaveData)

	go func() {
		zap.L().Named("restful").Info("Starting RESTful Server ...")
		zap.L().Named("restful").Info(" listen on addr ", zap.String("addr", server.addr))
		zap.L().Named("restful").Info(" TLS  disable")

		err := http.ListenAndServe(server.addr, server.router)
		if err != nil {
			panic(fmt.Errorf("http.ListenAndServe error:%s", err.Error()))
		}

		zap.L().Named("restful").Info("RESTful server started", zap.String("address", server.addr), zap.Bool("TLS", false))
	}()

	return nil
}

// SaveData 数据上链保存
func (server *RestServer) SaveData(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// 从http请求body中读取数据
	zap.L().Named("restful").Info("New Request Connected ", zap.String("remoteaddr", r.RemoteAddr))

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Read Body error:%s", err.Error())
		return
	}

	zap.L().Named("restful").Info("Receive Message", zap.ByteString("message", data))

	qidata := &hlserver.QIData{MessageID: "", TranCode: "HL20"}

	zap.L().Named("restful").Debug("WaitResponse")
	defer zap.L().Named("restful").Debug("WaitResponse end")

	rspchan := make(chan hlserver.RspData)
	recvData := hlserver.RecvData{Data: *qidata, RspChan: rspchan}

	server.recvChan <- recvData

	zap.L().Named("restful").Debug("send data to recvChan")

	select {
	case rspdata := <-rspchan:
		zap.L().Named("restful").Debug("receive response", zap.Any("response", rspdata))

		msg, err := json.Marshal(rspdata)
		if err != nil {
			zap.L().Named("restful").Error(err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		fmt.Fprintf(w, string(msg))

		return
	}
}
