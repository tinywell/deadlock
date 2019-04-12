# deadlock
日志死锁研究

## 问题说明
进行大并发测试时，会偶发性导致程序死锁，本项目提取原程序核心逻辑框架，用于研究其死锁原因。

死锁时相关 goroutine 信息如下
```
goroutine 2712950 [semacquire, 3 minutes]:
sync.runtime_SemacquireMutex(0xc000104164, 0x0)
	/usr/local/go/src/runtime/sema.go:71 +0x3d
sync.(*Mutex).Lock(0xc000104160)
	/usr/local/go/src/sync/mutex.go:134 +0x109
go.uber.org/zap/zapcore.(*lockedWriteSyncer).Write(0xc000104160, 0xc003ade800, 0x22f, 0x400, 0xbbd5c0, 0x855b21, 0x7)
	/go/pkg/mod/go.uber.org/zap@v1.9.1/zapcore/write_syncer.go:65 +0x31
go.uber.org/zap/zapcore.(*ioCore).Write(0xc0000fa3f0, 0x2, 0xbf2432629218f4c5, 0xc4bcad557a, 0xbbd5c0, 0x855b21, 0x7, 0x85d76b, 0x17, 0x1, ...)
	/go/pkg/mod/go.uber.org/zap@v1.9.1/zapcore/core.go:90 +0x107
go.uber.org/zap/zapcore.(*CheckedEntry).Write(0xc003aa6630, 0x0, 0x0, 0x0)
	/go/pkg/mod/go.uber.org/zap@v1.9.1/zapcore/entry.go:215 +0x119
go.uber.org/zap.(*Logger).Error(0xc004f94180, 0x85d76b, 0x17, 0x0, 0x0, 0x0)
	/go/pkg/mod/go.uber.org/zap@v1.9.1/logger.go:203 +0x7f
deadlock/hlserver/restful.(*RestServer).WaitResponse(0xc000104300, 0x8e9200, 0xc002398380, 0xc00199d9a0)
	/deadlock/agent/hlserver/restful/server.go:126 +0x63c
deadlock/hlserver/restful.(*RestServer).SaveData(0xc000104300, 0x8e9200, 0xc002398380, 0xc001fcb600, 0x0, 0x0, 0x0)
	/deadlock/agent/hlserver/restful/server.go:93 +0x5d0
github.com/julienschmidt/httprouter.(*Router).ServeHTTP(0xc0000ee2c0, 0x8e9200, 0xc002398380, 0xc001fcb600)
	/go/pkg/mod/github.com/julienschmidt/httprouter@v1.2.0/router.go:334 +0x948
net/http.serverHandler.ServeHTTP(0xc000112410, 0x8e9200, 0xc002398380, 0xc001fcb600)
	/usr/local/go/src/net/http/server.go:2774 +0xa8
net/http.(*conn).serve(0xc007fdee60, 0x8e9980, 0xc00c3047c0)
	/usr/local/go/src/net/http/server.go:1878 +0x851
created by net/http.(*Server).Serve
	/usr/local/go/src/net/http/server.go:2884 +0x2f4
```


## 服务启动
- 容器方式
```
make docker

cd ./example/docker
docekr-compose up -d
```

- 本地服务
```
make go 
cd ./build/bin
./deadlock &
```

## 测试用客户端
`./example/benchmark.go` 为测试用客户端
示例：
```
go run benchmark.go -para 500 -count 600
```
