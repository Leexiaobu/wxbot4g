package handler

import (
	"encoding/json"
	"github.com/eatmoreapple/openwechat"
	"strconv"
	"strings"
	"time"
	"wxbot4g/db"
)

// 检查是否需要保存
func checkNeedSave(message *openwechat.Message) bool {
	return message.IsText() || message.IsEmoticon() || message.IsPicture() || message.IsMedia()
}

func saveToDb(ctx *openwechat.MessageContext) {
	// TODO 需要解析成支持的结构体
	slew, _ := ctx.Bot.GetCurrentUser()
	sender, _ := ctx.Sender()
	senderUser := sender.NickName
	groupName := ""
	if ctx.IsSendByGroup() {
		// 取出消息在群里面的发送者
		senderInGroup, _ := ctx.SenderInGroup()
		senderUser = senderInGroup.NickName
		groupName = sender.NickName
	}
	msgStr, _ := json.Marshal(ctx)
	command := 0
	if strings.HasPrefix(ctx.Content, "#") {
		command = 1
	}
	msg := Message{
		Uin:          strconv.FormatInt(slew.Uin, 10),
		MsgID:        ctx.MsgId,
		MsgType:      int(ctx.MsgType),
		Content:      ctx.Content,
		SendUserName: senderUser,
		GroupName:    groupName,
		Read:         0,
		Command:      int8(command),
		BaseStr:      string(msgStr),
		DataTime:     time.Now().In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05"),
	}
	db.MysqlCon.Save(&msg)
	ctx.Next()
}

type Message struct {
	ID           int    `json:"id" gorm:"column:id"`
	Uin          string `json:"uin" gorm:"column:uin"`
	MsgID        string `json:"msg_id" gorm:"column:msg_id"`
	MsgType      int    `json:"msg_type" gorm:"column:msg_type"`
	Content      string `json:"content" gorm:"column:content"`
	SendUserName string `json:"send_user_name" gorm:"column:send_user_name"`
	GroupName    string `json:"group_name" gorm:"column:group_name"`
	Read         int8   `json:"read" gorm:"column:read"`
	Command      int8   `json:"command" gorm:"column:command"`
	BaseStr      string `json:"base_str" gorm:"column:base_str"`
	DataTime     string `json:"data_time" gorm:"column:data_time"`
}

func (m *Message) TableName() string {
	return "message"
}
