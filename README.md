# 1、ubuntu环境开发

## 1.1、安装golang环境
```
wget http://go.dev/dl/go1.17.9.linux-amd64.tar.gz
sudo tar -xvzf go1.17.9.linux-amd64.tar.gz -C /usr/local/


vi /etc/profile
export GOPATH=/root/goPath
export GOROOT=/usr/local/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

source /etc/profile

go version
```

## 1.2、基础配置

国内代理
```
go env -w GOPROXY=https://goproxy.cn,direct
```

设置私库(没私库可忽略)
```
go env -w GOPRIVATE=gitlab.XX.com
```

关闭包校验并开启gomodule
```
go env -w GOSUMDB=off
go env -w GO111MODULE=on   
```
## 1.3、安装相关协议生成工具

安装protoc
```
sudo apt install libprotobuf-dev protobuf-compiler
protoc --version

```
安装protoc-gen-go(proto生成go结构)
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@master
```

安装protoc-gen-vkit(proto生成go rpc协议)
```
go install github.com/visonlv/protoc-gen-vkit@1.0.0
```

安装rotoc-gen-swagger(proto生成go swagger api文档)
```
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@v1.16.0
```

安装rotoc-gen-validate(proto生成参数校验逻辑)
```
go install github.com/envoyproxy/protoc-gen-validate@master
```

## 1.4、编译项目
```
mkdir goWork && cd goWork
git clone https://github.com/visonlv/vkit-example.git
cd vkit-example
go mod tidy
go build .
```
## 1.5、启动项目

配置数据库
```
go run main.go
```
