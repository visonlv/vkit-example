package model

import (
	"time"

	"github.com/visonlv/go-vkit/mysqlx"
	"github.com/visonlv/go-vkit/utilsx"
	"gorm.io/gorm"
)

// 用户
var (
	user_model = &UserModel{}
)

type UserModel struct {
	Id         string    `gorm:"primaryKey;type:varchar(64);comment:主键id"`
	CreatedAt  time.Time `gorm:"comment:创建时间"` // 在创建时，如果该字段值为F零值，则使用当前时间填充
	UpdatedAt  time.Time `gorm:"comment:更新时间"` // 在创建时该字段值为零值或者在更新时，使用当前时间戳秒数填充
	CreateUser string    `gorm:"type:varchar(128);comment:创建用户"`
	UpdateUser string    `gorm:"type:varchar(128);comment:更新用户"`
	IsDelete   int       `gorm:"type:tinyint;not null;default:0;comment:删除状态 0正常 1删除"`

	Name     string `gorm:"type:varchar(128);comment:用户昵称"`
	Password string `gorm:"type:varchar(128);comment:密码"`
	Email    string `gorm:"type:varchar(64);comment:邮箱(唯一)"`
}

func (a *UserModel) BeforeCreate(tx *gorm.DB) error {
	a.Id = utilsx.GenUuid()
	return nil
}

func (*UserModel) TableName() string {
	return "t_user"
}

func UserAddWithTransaction(tx *mysqlx.MysqlClient, m1 *UserModel, m2 *UserModel) error {
	return getTx(tx).Transaction(func(newTx *mysqlx.MysqlClient) error {
		if err := UserAdd(newTx, m1); err != nil {
			return err
		}
		if err := UserAdd(newTx, m2); err != nil {
			return err
		}
		return nil
	})
}

func UserAdd(tx *mysqlx.MysqlClient, m *UserModel) error {
	err := getTx(tx).Model(user_model).Insert(m)
	return err
}

func UserDel(tx *mysqlx.MysqlClient, id string) error {
	result := getTx(tx).Model(user_model).Where("id = ?", id).Update("is_delete", 1)
	return result.GetDB().Error
}

func UserGet(tx *mysqlx.MysqlClient, id string) (*UserModel, error) {
	item := &UserModel{}
	result := getTx(tx).Where("id = ? AND is_delete = ?", id, 0).First(item)
	return item, result.GetDB().Error
}

func UserGetByEmail(tx *mysqlx.MysqlClient, email string) (*UserModel, error) {
	item := &UserModel{}
	result := getTx(tx).Where("email = ? AND is_delete = ?", email, 0).First(item)
	return item, result.GetDB().Error
}

func UserEmailExists(tx *mysqlx.MysqlClient, email string) (bool, error) {
	item := &UserModel{}
	has, err := getTx(tx).Where("email = ? AND is_delete = ?", email, 0).Exists(item)
	return has, err
}

func UserUpdate(tx *mysqlx.MysqlClient, m *UserModel) (*UserModel, error) {
	err := getTx(tx).UpdateEx(m)
	return m, err
}

func UserPage(tx *mysqlx.MysqlClient, userId, name string, pageIndex int32, pageSize int32) (list []*UserModel, total int32, err error) {
	query := getTx(tx).Model(user_model).Where("is_delete = ?", 0)
	if name != "" {
		query = query.Where("name like ?", "%"+name+"%")
	}
	if userId != "" {
		query = query.Where("id = ?", userId)
	}
	query = query.Order("created_at desc")
	err = query.FindPage(pageIndex, pageSize, &list, &total)
	return
}
