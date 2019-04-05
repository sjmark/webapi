package api_loginc

import (
	"webapi/config"
	"webapi/business/api_loginc/login"
	"webapi/business/api_loginc/user"
)

var mediator *Mediator

type Mediator struct {
	LoginMod loginInterface
	UserMod  userInterface
}

func init() {
	mediator = &Mediator{
		LoginMod: &login.LoginLoginc{},
		UserMod:  &user.UserLoginc{},
	}
}

func Init(conf *config.Config) {
	mediator.LoginMod.(*login.LoginLoginc).Config = conf
	mediator.UserMod.(*user.UserLoginc).Config = conf
}

func GetMadiator() *Mediator {
	return mediator
}
