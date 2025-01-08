.PHONY: all build run gotool clean help

BINARY="bluebell"

all: gotool build

build:
	@go build -o ${BINARY}.
#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY}
#SET CGO_ENABLED=0 SET GOOS=linux SET GOARCH=amd64 go build -o ${BINARY}

run:
	@go build -o ${BINARY} . && ${BINARY}
#	@@go run ./main.go config/config.yaml

gotool:
	go fmt ./
	go vet ./

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

help:
	@echo "make - 格式化 Go 代码, 并编译生成二进制文件"
	@echo "make build - 编译 Go 代码, 生成二进制文件"
	@echo "make run - 直接运行 Go 代码"
	@echo "make clean - 移除二进制文件和 vim swap files"
	@echo "make gotool - 运行 Go 工具 'fmt' and 'vet'"
