package auth

import (
	"path"
	"sync"

	"github.com/visonlv/go-vkit/logger"
)

type AuthRole struct {
	Urls []string //角色对应的Api列表
	Code string   //角色代码
}

type Auth struct {
	WhiteUrls []string  //白名单 支持正则匹配，如:/rpc/sso/UserService.*,/rpc/sso/*
	Roles     *sync.Map // 每个角色对应的ApiId列表 [角色代码]=角色权限信息
}

func NewAuth(whiteUrls []string) *Auth {
	return &Auth{
		WhiteUrls: whiteUrls,
		Roles:     new(sync.Map),
	}
}

func (a *Auth) Start() {
	a.ResetAll()
}

func (a *Auth) ResetAll() {
	cc := &AuthRole{
		Code: "ADMIN",
		Urls: []string{
			"/rpc/vkit-example/UserService.Add",
			"/rpc/vkit-example/UserService.Update",
			"/rpc/vkit-example/UserService.Page",
			"/rpc/vkit-example/UserService.Del",
			"/rpc/vkit-example/UserService.Get",
		},
	}
	a.Roles.Store("ADMIN", cc)
}

func (a *Auth) IsWhite(url string) bool {
	//判断白名单
	for _, v := range a.WhiteUrls {
		b, err := path.Match(v, url)
		if err != nil {
			logger.Errorf("[auth] IsWhite match url:%s err:%s", url, err)
			return false
		}
		if b {
			return true
		}
	}
	return false
}

func (a *Auth) IsPemission(roleList []string, url string) bool {
	//判断各个角色权限 匹配一个则满足
	for _, code := range roleList {
		r, ok := a.Roles.Load(code)
		if ok {
			item := r.(*AuthRole)
			for _, v := range item.Urls {
				b, err := path.Match(v, url)
				if err != nil {
					logger.Errorf("[auth] IsPemission match url:%s err:%s", url, err)
					return false
				}
				if b {
					return true
				}
			}
		}
	}
	return false
}
