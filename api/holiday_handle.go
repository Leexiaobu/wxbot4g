package api

import (
	"github.com/eatmoreapple/openwechat"
	"strings"
	"wxbot4g/holiday"
)

type HolidayHandler struct {
	Next
}

// Do 校验参数的逻辑
func (h *HolidayHandler) Do(ctx *openwechat.MessageContext) (bool, error) {
	if strings.EqualFold(ctx.Content, "#fj") {
		day := holiday.GetToDay()
		ctx.ReplyText(day)
		return true, nil
	}
	return false, nil
}
