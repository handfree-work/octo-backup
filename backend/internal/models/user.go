package models

import (
	"handfree-work/web-restic/internal/base/db_"
)

type User struct {
	db_.BaseModel
	db_.LogicDelete
	AppKey       string `json:"appKey" gorm:"size:100;comment:应用Key;index"`
	Avatar       string `json:"avatar" gorm:"comment:用户头像"`
	Username     string `json:"userName" gorm:"size:100;comment:用户登录名;index:idx_username"`
	Password     string `json:"password"  gorm:"size:100;comment:用户登录密码"`
	NickName     string `json:"nickName" gorm:"size:100;default:系统用户;comment:用户昵称"`
	PhoneCode    string `json:"phoneCode"  gorm:"size:100;comment:手机区号"`
	Mobile       string `json:"mobile"  gorm:"size:100;comment:手机号;index:idx_mobile"`
	Email        string `json:"email"  gorm:"size:100;comment:邮箱;index:idx_email"`
	Gender       string `json:"gender"  gorm:"size:20;comment:性别"`
	UserType     int32  `json:"userType"  gorm:"comment:用户类型"`
	InviteCode   string `json:"inviteCode" gorm:"size:100;comment:邀请码;unique"`
	InviteUserId int64  `json:"inviteUserId" gorm:"comment:邀请人id"`
}

func (u *User) NewUser(id *int64) *User {
	return &User{BaseModel: db_.NewBaseModel(id)}
}

func (u *User) EncryptPassword(password string) error {
	return nil
}
