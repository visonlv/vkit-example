// Code generated by protoc-gen-infore.
// versions:
// - protoc-gen-infore v1.0.0

package handler

import (
	context "context"
	"time"

	"github.com/visonlv/go-vkit/errorsx"
	"github.com/visonlv/vkit-example/handler/auth"
	"github.com/visonlv/vkit-example/model"
	pb "github.com/visonlv/vkit-example/proto/vkit_example"
	"github.com/visonlv/vkit-example/utils"
	"gorm.io/gorm"
)

type AuthService struct {
}

func (the *AuthService) Login(ctx context.Context, req *pb.LoginReq, resp *pb.LoginResp) error {
	// 获取用户信息
	userModel, err := model.UserGetByEmail(nil, req.Username)
	if err != nil {
		resp.Code = errorsx.FAIL.Code
		if err == gorm.ErrRecordNotFound {
			resp.Msg = "邮箱不存在"
		} else {
			resp.Msg = err.Error()
		}
		return nil
	}
	// 校验密码
	if req.Password != userModel.Password {
		resp.Code = errorsx.FAIL.Code
		resp.Msg = "密码错误"
		return nil
	}

	roles := make([]string, 0)
	roles = append(roles, "ADMIN")

	//生成token
	t := &utils.Account{
		Id:       userModel.Id,
		Roles:    roles,
		Metadata: "{}",
	}
	tt := utils.GenToken(t, time.Hour*24)
	auth.AddToken(userModel.Id, tt)

	resp.Token = tt

	return nil
}

func (the *AuthService) Logout(ctx context.Context, req *pb.LogoutReq, resp *pb.LogoutResp) error {
	auth.DeleteToken(utils.GetUserIdFromContext(ctx))
	return nil
}
