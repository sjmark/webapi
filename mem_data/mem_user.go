package mem_data

import (
	"sync"
	"webapi/models"
	"github.com/mkideal/log"
	"time"
	"webapi/consts"
)

var userMap sync.Map

type userStore *models.UserInfo
type userOpt func(u userStore) error

func LoadUserStore(uid int64, opts ...userOpt) (mem models.UserInfo, err error) {
	var userInfo = new(models.UserInfo)
	var isUp bool
	if user, ok := userMap.Load(uid); ok {
		userInfo = user.(*models.UserInfo)
	} else {
		userTb := new(models.UserTb)
		userTb.Id = uid

		err = userTb.GetUserInfoByUid()

		if err == nil {
			userInfo.Uid = uid
			userInfo.Sex = userTb.Sex
			userInfo.City = userTb.City
			userInfo.Mobile = userTb.Mobile
			userInfo.Avatar = userTb.Avatar
			userInfo.Nickname = userTb.Nickname
			userInfo.Sid = userTb.Sid
			userInfo.SidExpTime = userTb.SidExpTime
		}

		isUp = true
	}

	for _, opt := range opts {
		if err := opt(userInfo); err != nil {
			log.Error("LoadUserStore opt error:%v,uid:%d", err, uid)
			continue
		}
		isUp = true
	}

	if isUp {
		userMap.Store(uid, userInfo)
	}

	return *userInfo, nil
}

func LoadUserStoreByMobile(mobile string) (mem models.UserInfo, err error) {

	userTb := new(models.UserTb)
	userTb.Mobile = mobile

	err = userTb.GetUserInfoByMobile()

	if err != nil {
		log.Error("userTb.GetUserInfoByMobile error:%v,mobile:%s", err, mobile)
		err = nil
		return
	}

	mem.Uid = userTb.Id
	mem.Sex = userTb.Sex
	mem.City = userTb.City
	mem.Mobile = userTb.Mobile
	mem.Avatar = userTb.Avatar
	mem.Nickname = userTb.Nickname
	mem.Sid = userTb.Sid
	mem.SidExpTime = userTb.SidExpTime
	return
}

func NewUser(mobile, sid string) userOpt {
	return func(u userStore) error {
		var expTime = time.Now().Add(time.Second*consts.SidExpTime).Unix()
		var userTb = new(models.UserTb)

		userTb.Mobile = mobile
		userTb.Sid = sid
		userTb.SidExpTime = expTime
		uid, err := userTb.NewUserToTb()

		if err != nil {
			return err
		}

		u.Uid = uid
		u.Mobile = mobile
		u.Sid = sid
		u.SidExpTime = expTime

		return nil
	}
}

func UpdateUserSid(sid string) userOpt {
	return func(u userStore) error {
		var expTime = time.Now().Add(time.Second*consts.SidExpTime).Unix()
		var userTb = new(models.UserTb)
		userTb.Id = u.Uid
		userTb.Sid = sid
		userTb.SidExpTime = expTime
		if err := userTb.UpdateUserSid(); err != nil {
			return err
		}
		u.Sid = sid
		u.SidExpTime = expTime
		return nil
	}
}

func UpdateUserInfo(city, avatar, nickname string) userOpt {
	return func(u userStore) error {
		userTb := new(models.UserTb)
		userTb.Id = u.Uid
		userTb.City = city
		userTb.Avatar = avatar
		userTb.Nickname = nickname

		if err := userTb.UpdateUserInfo(); err != nil {
			return err
		}
		u.City = city
		u.Avatar = avatar
		u.Nickname = nickname
		return nil
	}
}
