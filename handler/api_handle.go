package handler

import (
	"github.com/eatmoreapple/openwechat"
	"strings"
)

func checkIsApiMsg(message *openwechat.Message) bool {
	// 通知消息和自己发的不处理
	return !message.IsSendBySelf() && strings.HasPrefix(message.Content, "#")
}

func apiMessageHandle(ctx *openwechat.MessageContext) {
	ctx.Message.ReplyText("123")
	ctx.Next()
}
