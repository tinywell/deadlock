BUILDDIR = ./build
BUILDBIN = ./build/bin
BUILDDOCKER = ./build/docker
BIN = deadlock

IMAGETAG = 1.0
DOCKERTAG = deadlock/agent:$(IMAGETAG)

.PHONY: go
go:
	go build -o $(BUILDBIN)/$(BIN)
	go build -o $(BUILDBIN)/benchmark example/http/benchmark.go

.PHONY: linux
linux:
	GOOS=linux go build -o $(BUILDBIN)/$(BIN)

.PHONY: docker
docker: 
	docker build -f images/Dockerfile -t $(DOCKERTAG) .

.PHONY: clean
clean:
	rm -f *.o
	rm -rf $(BUILDDIR)/*
