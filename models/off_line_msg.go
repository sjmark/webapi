package models

import "webapi/protos/im_proto"

type OffLineMsg struct {
	MsgId       int64  `xorm:"not null pk autoincr BIGINT(20)"` // 消息发送者ID消息的 唯一指纹码（即消息ID），用于去重等场景，单机情况下此id可能是个自增值、分布式场景下可能是类似于UUID这样的东西
	SenderUid   int64  `xorm:"INDEX default 0 bigint(20)"`      // 消息发送者ID
	ReceiverUid int64  `xorm:"INDEX default 0  bigint(20)" `    // 消息接收者ID
	MsgType     uint32  `xorm:"default 0 TINYINT(4)"`            // 消息类型（标识此条消息是：文本、图片还是语音留言等）
	MsgContent  string `xorm:"default NULL VARCHAR(1024)"`      // 消息内容（如果是图片或语音留言等类型，由此字段存放的可能是对应文件的存储地址或CDN的访问URL）
	Status      uint8  `xorm:"INDEX default 0 TINYINT(4)"`      // 消息状态 0表示未发送 1 表示完成 后期定时移除这些过期的数据，或者保存到正常消息的库中
	CreatedTime int64  `xorm:"default 0 BIGINT(20)"`            // 消息发出时的时间戳（如果是个跨国IM，则此时间戳可能是GMT-0标准时间）
}

// proto文件 暂定这样
type OffLineMsgRes struct {
	MsgType    uint32  `json:"msg_type"`
	MsgContent string `json:"msg_content"`
}

// SELECT msg_id, send_time, msg_type, msg_content  FROM off_line_msg  WHERE receiver_uid=? and status =0
func GetOffLineMsgByUid(revUid int64) []*protos.OffLineMsgRes {
	return []*protos.OffLineMsgRes{}
}
