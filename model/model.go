package model

import (
	"github.com/visonlv/go-vkit/mysqlx"
	"github.com/visonlv/vkit-example/app"
)

const ()

// InitTable 初始化数据库表
func InitTable() {
	// 自动建表
	app.Mysql.AutoMigrate(&UserModel{})
	// 创建第一个测试用户
	exist, _ := UserEmailExists(nil, "test")
	if !exist {
		UserAdd(nil, &UserModel{
			Name:     "测试",
			Email:    "test",
			Password: "123456",
		})
	}
}

func getTx(tx *mysqlx.MysqlClient) *mysqlx.MysqlClient {
	if tx == nil {
		return app.Mysql
	}
	return tx
}
