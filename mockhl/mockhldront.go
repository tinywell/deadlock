package mockhl

import (
	"deadlock/hlserver"
	"time"
)

type MockHlfront struct {
}

func (mh *MockHlfront) TransactionInvoke(channel, cc, fcn, data string, crypto bool, member []string) *hlserver.RspData {
	rsp := &hlserver.RspData{
		MessageID: "",
		Code:      200,
	}
	rsp.Data.TxID = time.Now().UTC().String()
	rsp.Data.Message = "测试"
	time.Sleep(time.Microsecond * 100)
	return rsp
}
