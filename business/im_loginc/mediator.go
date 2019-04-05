package im_loginc

import (
	"webapi/business/im_loginc/msg"
)

var mediator *Mediator

type Mediator struct {
	MsgMod   msgInterface
}

func init() {
	mediator = &Mediator{
		MsgMod:   &msg.MsgBussiness{},
	}
}

func GetMadiator() *Mediator {
	return mediator
}
