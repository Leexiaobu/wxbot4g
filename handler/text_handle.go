package handler

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"wxbot4g/logger"
)

// 处理文本消息
func textMessageHandle(ctx *openwechat.MessageContext) {
	sender, _ := ctx.Sender()
	senderUser := sender.NickName
	if ctx.IsSendByGroup() {
		// 取出消息在群里面的发送者
		senderInGroup, _ := ctx.SenderInGroup()
		senderUser = fmt.Sprintf("%v[%v]", senderInGroup.NickName, senderUser)
	}
	logger.Log.Infof("[收到新文字消息] == 发信人：%v ==> 内容：%v", senderUser, ctx.Content)
	ctx.Next()
}
