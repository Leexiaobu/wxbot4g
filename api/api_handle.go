package api

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"io"
	"log"
	"net/http"
	"strings"
)

var Head = &NullHandler{}

func init() {
	Head.
		SetNext(&WeatherHandle{}).
		SetNext(&AiPicHandler{}).
		SetNext(&HolidayHandler{}).
		SetNext(&XiaoAiHandler{})
}
func CheckIsApiMsg(message *openwechat.Message) bool {
	// 通知消息和自己发的不处理
	return !message.IsSendBySelf() && strings.HasPrefix(message.Content, "#")
}

func ApiMessageHandle(ctx *openwechat.MessageContext) {
	_, err := Head.Run(ctx)
	if err != nil {
		fmt.Println("Api error:" + err.Error())
	}
	ctx.Next()
}

type ApiHandle interface {
	Do(ctx *openwechat.MessageContext) (bool, error)
	SetNext(h ApiHandle) ApiHandle
	Run(ctx *openwechat.MessageContext) (bool, error)
}
type Next struct {
	nextHandler ApiHandle
}

func (n *Next) SetNext(h ApiHandle) ApiHandle {
	n.nextHandler = h
	return h
}

// Run 执行
func (n *Next) Run(ctx *openwechat.MessageContext) (breakHandle bool, err error) {
	if n.nextHandler != nil {
		isHandled, err := (n.nextHandler).Do(ctx)
		if err != nil || isHandled {
			return isHandled, err
		}
		return (n.nextHandler).Run(ctx)
	}
	return false, nil
}

type NullHandler struct {
	Next
}

// Do 空Handler的Do
func (h *NullHandler) Do(ctx *openwechat.MessageContext) (err error) {
	return
}
func getStringResult(url string) string {
	response, _ := http.Get(url)
	defer response.Body.Close()
	body, err2 := io.ReadAll(response.Body)
	if err2 != nil {
		log.Println("ioutil read error")
	}
	return string(body)
}
