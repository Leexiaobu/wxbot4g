package api

import (
	"github.com/eatmoreapple/openwechat"
	"wxbot4g/config"
)

type XiaoAiHandler struct {
	Next
}

var xiaoaiconf = config.Config.ApiConfig.XiaoAiConfig

// Do 校验参数的逻辑
func (h *XiaoAiHandler) Do(ctx *openwechat.MessageContext) (bool, error) {
	//txt := strings.Replace(ctx.Content, "#", "", 1)
	//requestUrl := fmt.Sprintf(xiaoaiconf.Url, txt)
	//fmt.Println(requestUrl)
	//result := getStringResult(requestUrl)
	ctx.ReplyText("小爱也知道怎么回答")
	return true, nil
}
