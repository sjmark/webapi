package queuebiz

import "webapi/models"

type BusinessInterface interface {
	CommonBusiness(param models.MfwTask) error
}
type GivenBusinessInterface interface {
	GivenBusiness(areaID, name, title, url string, source, urlType uint8) error
}

// 短信发送什么的可以写到这里，发送消息通知等
