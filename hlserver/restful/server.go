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

// RESTful接口地址
const (
	// URLSaveData 数据上链 RESTful接口地址 POST
	URLSaveData         = "/agent/transaction"
	WaitResponseTimeOut = time.Second * 60
)

// RestServer RESTful接口
type RestServer struct {
	addr     string             // 端口
	router   *httprouter.Router // 路由
	recvChan chan<- hlserver.RecvData
}

// NewServer 生成新RestServer对象
// params：
//   - addr string  服务监听地址
// return:
//   - RestServer  RESTful接口对象
func NewServer(addr string) *RestServer {
	return &RestServer{
		addr:   addr,
		router: httprouter.New(),
	}
}

// Start 启动restful服务
// params：
//   - recvchan 消息处理 channel
// return:
//   - error  错误信息
func (server *RestServer) Start(recvchan chan<- hlserver.RecvData) error {
	server.recvChan = recvchan

	server.router.POST(URLSaveData, server.SaveData)

	go func() {
		zap.L().Named("restful").Info("Starting RESTful Server ...")
		// logger.Infof(" listen on '%s' ", server.addr)
		var err error

		zap.L().Named("restful").Info(" TLS  disable")
		err = http.ListenAndServe(server.addr, server.router)

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
	logRequest(r)
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Read Body error:%s", err.Error())
		return
	}
	zap.L().Named("restful").Info("Receive Message", zap.ByteString("message", data))
	qidata := &hlserver.QIData{}
	err = json.Unmarshal(data, qidata)
	if err != nil {
		reterr := fmt.Errorf("Unmarshal data (%s) error:%s", string(data), err)
		// logger.Error(reterr)
		rspData := *hlserver.CommonRspWithError(reterr)
		rspData.MessageID = qidata.MessageID
		SendReturn(w, rspData)
		return
	}

	// 此处不能用 goroutine 处理，当前方法结束后，w会关闭，导致返回信息无法写入
	server.WaitResponse(w, qidata)
}

// WaitResponse 将接口数据交给 handler 处理并等待处理结果
func (server *RestServer) WaitResponse(w http.ResponseWriter, qidata *hlserver.QIData) {
	zap.L().Named("restful").Debug("WaitResponse")
	defer zap.L().Named("restful").Debug("WaitResponse end")
	rspchan := make(chan hlserver.RspData)
	recvData := hlserver.RecvData{
		Data:    *qidata,
		RspChan: rspchan,
	}
	select {
	case server.recvChan <- recvData:
		// logger.Debug("send data to recvChan")
	case <-time.After(WaitResponseTimeOut):
		// logger.Error("send data to recvChan timeout")
		rspData := *hlserver.CommonRspWithError(fmt.Errorf("send data to recvChan timeout"))
		rspData.MessageID = qidata.MessageID
		SendReturn(w, rspData)
		return
	}
	// server.recvChan <- recvData
	zap.L().Named("restful").Debug("send data to recvChan")
	select {
	case rspdata := <-rspchan:
		zap.L().Named("restful").Debug("receive response", zap.Any("response", rspdata))
		if rspdata.Code != hlserver.RSPCODE_SUCCESS {
			zap.L().Named("restful").Error("Handle message faild ", zap.Int("code", rspdata.Code), zap.Any("message", rspdata.Data))
		}
		SendReturn(w, rspdata)
		return
	case <-time.After(WaitResponseTimeOut): // TODO: 超时时间配置
		zap.L().Named("restful").Error("Receive rspdata timeout")
		rspData := *hlserver.CommonRspWithError(fmt.Errorf("Receive rspdata timeout"))
		rspData.MessageID = qidata.MessageID
		SendReturn(w, rspData)
		return
	}
}

func logRequest(r *http.Request) {
	zap.L().Named("restful").Info("New Request Connected ", zap.String("remoteaddr", r.RemoteAddr))
}

// SendReturn 将处理结果返回请求端
// params：
// - w http.ResponseWriter        响应输出
// - data interface{}				 返回数据
func SendReturn(w http.ResponseWriter, data interface{}) {
	msg, err := json.Marshal(data)
	if err != nil {
		zap.L().Named("restful").Error(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	fmt.Fprintf(w, string(msg))
}
