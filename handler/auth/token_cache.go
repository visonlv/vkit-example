package auth

import (
	"fmt"
	"sync"
)

var tokens = new(sync.Map)

type cacheValue struct {
	expTime int64
	userId  string
}

func AddToken(userId, token string) {
	tokens.Store(userId, token)
}

func DeleteToken(userId string) {
	tokens.Delete(userId)
}

func UserTokenExist(userId, token string) error {
	r, ok := tokens.Load(userId)
	if !ok {
		return fmt.Errorf("请求令牌失效，请求失败!!")

	}
	cacheToken := r.(string)
	if cacheToken != token {
		return fmt.Errorf("账号在其他设备上登录，请确认!")
	}
	return nil
}
