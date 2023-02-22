package main

import (
	"github.com/visonlv/vkit-example/app"
	"github.com/visonlv/vkit-example/model"
	"github.com/visonlv/vkit-example/server"
)

func main() {
	// 1. 初始化配置
	app.Init("./config.toml")

	model.InitTable()

	server.Start()

}
