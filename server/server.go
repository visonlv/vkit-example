package server

import (
	"fmt"
	"net/http"

	"github.com/visonlv/go-vkit/errorsx/neterrors"
	"github.com/visonlv/go-vkit/httphandler"
	"github.com/visonlv/go-vkit/logger"
	"github.com/visonlv/vkit-example/app"
	"github.com/visonlv/vkit-example/handler"
	"github.com/visonlv/vkit-example/handler/auth"
	"github.com/visonlv/vkit-example/utils"
)

var (
	whiteList = []string{
		"/rpc/vkit-example/AuthService.Login",
	}
	authObj = auth.NewAuth(whiteList)
)

func Start() {
	//初始化权限数据
	authObj.Start()
	//统一鉴权逻辑
	tokenCheck := func(w http.ResponseWriter, r *http.Request) error {
		// 判断白名单
		if authObj.IsWhite(r.RequestURI) {
			return nil
		}

		// 判断token
		tokenStr := r.Header.Get("AuthToken")
		if tokenStr == "" {
			return neterrors.Unauthorized("没有请求令牌，请求失败!")
		}

		tt, err := utils.ParseToken(tokenStr)
		if err != nil {
			return neterrors.Unauthorized("请求令牌失效，请求失败!")
		}

		if err := auth.UserTokenExist(tt.Id, tokenStr); err != nil {
			return neterrors.Unauthorized(err.Error())
		}

		if !authObj.IsPemission(tt.Roles, r.RequestURI) {
			return neterrors.Forbidden("资源没有权限!")
		}
		r.Header.Set("user_id", tt.Id)
		r.Header.Set("role_code", tt.Roles[0])
		return nil
	}

	//http 转发
	httpHandler := httphandler.NewHandler()
	httpHandler.WithAuthFunc(tokenCheck)
	err := httpHandler.RegisterApiEndpoint(handler.GetList(), handler.GetApiEndpoint())
	if err != nil {
		logger.Errorf("[main] RegisterApiEndpoint fail %s", err)
		panic(err)
	}
	http.HandleFunc("/rpc/", func(w http.ResponseWriter, r *http.Request) {
		httpHandler.Handle(w, r)
	})

	logger.Infof("[main] Listen port:%d", app.Cfg.Server.HttpPort)
	err = http.ListenAndServe(fmt.Sprintf(":%d", app.Cfg.Server.HttpPort), nil)
	if err != nil {
		logger.Errorf("[main] ListenAndServe fail %s", err)
		panic(err)
	}
}
