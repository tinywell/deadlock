package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	count   = 1
	paraNum = 10
)

var (
	wg       sync.WaitGroup
	paraChan chan struct{}
	sucess   int
	sucChan  chan struct{}
	sucSig   chan struct{}
)

func main() {
	flag.IntVar(&paraNum, "para", 10, "并发数")
	flag.IntVar(&count, "count", 10, "轮次。（并发数 * 轮次 = 总交易量）")
	flag.Parse()
	fmt.Printf("Para: %d  Count:%d  Total: %d\n", paraNum, count, paraNum*count)

	sucSig = make(chan struct{})
	sucChan = make(chan struct{}, count*paraNum)
	go countSuc()
	start := time.Now()
	runner()
	spend := time.Since(start)
	<-sucSig // 等待 countSuc() 统计交易成功数（否则 success 数据不准确，全局变量并发操作）
	fmt.Printf("Total: %d 笔  Success: %d 笔  Faild: %d 笔  SpendTime: %f s  TPS:%f 笔/s\n", count*paraNum, sucess, count*paraNum-sucess, spend.Seconds(), float64(sucess)/spend.Seconds())
}

func runner() {
	defer close(sucChan)
	url := "http://127.0.0.1:8000/agent/transaction"
	data := `{
		"messageid": "0002",
		"trancode": "HLF20",
		"trandata": {
			"bizType": "banklevel",
			"bizTypeDesc": "测试",
			"fromMspId": "test",
			"rewardRatio": 15
		},
		"crypto": false,
		"members": [
			""
		]
	}`
	paraChan = make(chan struct{}, paraNum)
	for i := 0; i < count*paraNum; i++ {
		paraChan <- struct{}{}
		wg.Add(1)
		go httpDo(url, data)
	}
	wg.Wait()
}

func countSuc() {
	for _ = range sucChan {
		sucess++
	}
	sucSig <- struct{}{}
}

func httpDo(url, data string) {
	defer wg.Done()
	defer func() { <-paraChan }()
	client := &http.Client{}

	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		// handle error
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
	sucChan <- struct{}{}
}
