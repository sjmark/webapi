package models

import (
	"github.com/go-xorm/xorm"
)

var (
	sessionRead  *xorm.Session
	sessionWrite *xorm.Session
)

func InitXormSession(write, read *xorm.Session) {
	sessionWrite = write
	sessionRead = read
}

type UserTb struct {
	Id         int64  `xorm:"not null pk autoincr BIGINT(20)"`
	Sex        int8   `xorm:"default 0 TINYINT(4)"`
	City       string `xorm:"VARCHAR(20)"`
	Mobile     string `xorm:"unique VARCHAR(20)"`
	Avatar     string `xorm:"VARCHAR(200)"`
	Password   string `xorm:"VARCHAR(32)"`
	Nickname   string `xorm:"INDEX VARCHAR(20)"`
	Sid        string `xorm:"VARCHAR(32)"`
	Level      int    `xorm:"default 0 INT(11)"`
	Balance    int    `xorm:"default 0 INT(11)"`
	Exp        int    `xorm:"default 0 INT(11)"`
	CreateTime int64  `xorm:"default 0 BIGINT(20)"`
	SidExpTime int64  `xorm:"default 0 BIGINT(20)"`
}

type UserInfo struct {
	Uid        int64
	Sex        int8
	City       string
	Mobile     string
	Avatar     string
	Nickname   string
	Sid        string
	SidExpTime int64
}

func (u *UserTb) NewUserToTb() (int64, error) {
	defer sessionWrite.Cols()

	if err := sessionWrite.Sync2(u); err != nil {
		return u.Id, err
	}

	_, err := sessionWrite.InsertOne(u)
	return u.Id, err
}

func (u *UserTb) GetUserInfoByMobile() error {
	defer sessionRead.Cols()
	_, err := sessionRead.Where("mobile=?", u.Mobile).Cols("id", "sid", "sid_exp_time").Get(u)
	return err
}

func (u *UserTb) GetUserInfoByUid() error {
	defer sessionRead.Cols()
	_, err := sessionRead.ID(u.Id).Cols("sex", "city", "mobile", "avatar", "nickname", "sid", "sid_exp_time").Get(u)
	return err
}

func (u *UserTb) UpdateUserSid() error {
	defer sessionRead.Cols()
	_, err := sessionRead.ID(u.Id).Cols("sid", "sid_exp_time").Update(u)
	return err
}

func (u *UserTb) UpdateUserInfo() error {
	defer sessionWrite.Cols()
	_, err := sessionWrite.ID(u.Id).Cols("city", "avatar", "nickname").Update(u)
	return err
}
