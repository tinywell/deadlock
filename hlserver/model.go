package hlserver

type QIData struct {
	MessageID string `json:"messageid"`
	TranCode  string `json:"trancode"`
	// TranData  string   `json:"trandata"`
	TranData interface{} `json:"trandata"`
	Crypto   bool        `json:"crypto"`
	Members  []string    `json:"members"`
}

type RecvData struct {
	Data    QIData
	RspChan chan<- RspData
}

type RspData struct {
	MessageID string `json:"messageid"`
	Code      int    `json:"code"`
	Data      struct {
		TxID    string `json:"txid"`
		Message string `json:"message"`
		OriMsg  string `json:"orimsg"` // QIData
	} `json:"data"`
}

type EventData struct {
	Event string      `json:"event"`
	TxID  string      `json:"txid"`
	Data  interface{} `json:"data"`
}

func CommonRspWithError(err error) *RspData {
	rsp := &RspData{}
	rsp.Code = RSPCODE_COMMONERROR
	rsp.Data.Message = err.Error()
	return rsp
}

const (
	RSPCODE_SUCCESS     = 200
	RSPCODE_INNERERROR  = 500
	RSPCODE_TIMEOUT     = 502
	RSPCODE_COMMONERROR = 100
)
