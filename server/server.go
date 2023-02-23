package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/visonlv/go-vkit/errorsx/neterrors"
	"github.com/visonlv/go-vkit/gate"
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

func tokenCheckFunc(w http.ResponseWriter, r *http.Request) error {
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

func logFunc(f gate.HandlerFunc) gate.HandlerFunc {
	return func(ctx context.Context, req *gate.HttpRequest, resp *gate.HttpResponse) error {
		startTime := time.Now()
		err := f(ctx, req, resp)
		costTime := time.Since(startTime)
		body, _ := req.Read()
		var logText string
		if err != nil {
			logText = fmt.Sprintf("fail cost:[%v] url:[%v] req:[%v] resp:[%v]", costTime.Milliseconds(), req.Uri(), string(body), err.Error())
		} else {
			logText = fmt.Sprintf("success cost:[%v] url:[%v] req:[%v] resp:[%v]", costTime.Milliseconds(), req.Uri(), string(body), string(resp.Content()))
		}
		logger.Infof(logText)
		return err
	}
}

func Start() {
	//初始化权限数据
	authObj.Start()

	h := gate.NewNativeHandler(
		gate.HttpAuthHandler(tokenCheckFunc),
		gate.HttpWrapHandler(logFunc),
	)
	err := h.RegisterApiEndpoint(handler.GetList(), handler.GetApiEndpoint())
	if err != nil {
		logger.Errorf("[main] RegisterApiEndpoint fail %s", err)
		panic(err)
	}
	http.HandleFunc("/rpc/", func(w http.ResponseWriter, r *http.Request) {
		h.Handle(w, r)
	})

	logger.Infof("[main] Listen port:%d", app.Cfg.Server.HttpPort)
	err = http.ListenAndServe(fmt.Sprintf(":%d", app.Cfg.Server.HttpPort), nil)
	if err != nil {
		logger.Errorf("[main] ListenAndServe fail %s", err)
		panic(err)
	}
}
