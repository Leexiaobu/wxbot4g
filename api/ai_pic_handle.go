package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	"wxbot4g/logger"
)

type AiPicHandler struct {
	Next
}

var isRun = false

// Do 校验参数的逻辑
func (h *AiPicHandler) Do(ctx *openwechat.MessageContext) (bool, error) {
	if strings.HasPrefix(ctx.Content, "#hh") && isRun {
		params := strings.Replace(ctx.Content, "#hh", "", 1)
		ctx.ReplyText("小爱正在作画中。。预计80s")
		logger.Log.Info(fmt.Sprintf("作画参数:%s", params))
		var base64Data, unix = genPicData(params)
		dist, _ := base64.StdEncoding.DecodeString(base64Data)
		fileName := fmt.Sprintf("ai_%d.png", unix)
		f, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
		f.Write(dist)
		open, _ := os.Open(fileName)
		ctx.ReplyImage(open)
		return true, nil
	}
	return false, nil
}

var url = "http://192.168.1.133:6969/generate-stream"

func genPicData(msg string) (string, int64) {
	defaultReq.Prompt = msg
	unix := time.Now().Unix()
	defaultReq.Seed = unix
	marshal, _ := json.Marshal(defaultReq)
	params := string(marshal)
	logger.Log.Info(params)
	resp, _ := http.Post(url, "application/json", strings.NewReader(params))
	body, _ := io.ReadAll(resp.Body)
	base64Data := strings.Replace(string(body), "event: newImage\nid: 1\ndata:", "", 1)
	return base64Data, unix
}

var defaultReq = &payload{
	Height:   512,
	NSamples: 1,
	Sampler:  "k_euler_ancestral",
	Scale:    12,
	Steps:    28,
	Uc:       "lowres, bad anatomy, text, error, worst quality, low quality, normal quality, jpeg artifacts, signature,",
	UcPreset: 0,
	Width:    896,
}

type payload struct {
	UcPreset int    `json:"ucPreset" gorm:"column:ucPreset"`
	Seed     int64  `json:"seed" gorm:"column:seed"`
	NSamples int    `json:"n_samples" gorm:"column:n_samples"`
	Width    int    `json:"width" gorm:"column:width"`
	Scale    int    `json:"scale" gorm:"column:scale"`
	Prompt   string `json:"prompt" gorm:"column:prompt"`
	Steps    int    `json:"steps" gorm:"column:steps"`
	Uc       string `json:"uc" gorm:"column:uc"`
	Height   int    `json:"height" gorm:"column:height"`
	Sampler  string `json:"sampler" gorm:"column:sampler"`
}
