package handler

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"io/ioutil"
	"net/http"
	"wxbot4g/config"
	"wxbot4g/logger"
	"wxbot4g/oss"
)

type AppMessageData struct {
	XMLName xml.Name `xml:"msg"`
	Appmsg  struct {
		Appid             string `xml:"appid,attr"`
		Sdkver            string `xml:"sdkver,attr"`
		Title             string `xml:"title"`
		Des               string `xml:"des"`
		Action            string `xml:"action"`
		Type              string `xml:"type"`
		Showtype          string `xml:"showtype"`
		Content           string `xml:"content"`
		URL               string `xml:"url"`
		Dataurl           string `xml:"dataurl"`
		Lowurl            string `xml:"lowurl"`
		Lowdataurl        string `xml:"lowdataurl"`
		Recorditem        string `xml:"recorditem"`
		Thumburl          string `xml:"thumburl"`
		Messageaction     string `xml:"messageaction"`
		Md5               string `xml:"md5"`
		Extinfo           string `xml:"extinfo"`
		Sourceusername    string `xml:"sourceusername"`
		Sourcedisplayname string `xml:"sourcedisplayname"`
		Commenturl        string `xml:"commenturl"`
		Appattach         struct {
			Totallen          string `xml:"totallen"`
			Attachid          string `xml:"attachid"`
			Emoticonmd5       string `xml:"emoticonmd5"`
			Fileext           string `xml:"fileext"`
			Fileuploadtoken   string `xml:"fileuploadtoken"`
			OverwriteNewmsgid string `xml:"overwrite_newmsgid"`
			Filekey           string `xml:"filekey"`
			Cdnattachurl      string `xml:"cdnattachurl"`
			Aeskey            string `xml:"aeskey"`
			Encryver          string `xml:"encryver"`
		} `xml:"appattach"`
		Weappinfo struct {
			Pagepath       string `xml:"pagepath"`
			Username       string `xml:"username"`
			Appid          string `xml:"appid"`
			Appservicetype string `xml:"appservicetype"`
		} `xml:"weappinfo"`
		Websearch string `xml:"websearch"`
	} `xml:"appmsg"`
	Fromusername string `xml:"fromusername"`
	Scene        string `xml:"scene"`
	Appinfo      struct {
		Version string `xml:"version"`
		Appname string `xml:"appname"`
	} `xml:"appinfo"`
	Commenturl string `xml:"commenturl"`
}

// APP????????????
func appMessageHandle(ctx *openwechat.MessageContext) {
	// ???????????????
	sender, _ := ctx.Sender()
	senderUser := sender.NickName
	if ctx.IsSendByGroup() {
		// ????????????????????????????????????
		senderInGroup, _ := ctx.SenderInGroup()
		senderUser = fmt.Sprintf("%v[%v]", senderInGroup.NickName, senderUser)
	}
	// ????????????
	var data AppMessageData
	if err := xml.Unmarshal([]byte(ctx.Content), &data); err != nil {
		logger.Log.Errorf("??????????????????: %v", err.Error())
		logger.Log.Debugf("????????????: %v", ctx.Content)
		return
	} else {
		logger.Log.Infof("[?????????????????????] == ????????????%v ==> Type???%v ==> ?????????%v ==> ??????APP: %v",
			senderUser, data.Appmsg.Type, data.Appmsg.Title, data.Appinfo.Appname)
		tt := data.Appmsg.Type
		dealType := []string{"2", "3", "4", "6", "15"}
		if !checkIsExist(dealType, tt) {
			logger.Log.Infof("????????????????????????????????????????????????????????????: %v", tt)
			return
		}
		// ??????????????????
		fileResp, err := ctx.GetFile()
		if err != nil {
			logger.Log.Errorf("??????????????????: %v", err.Error())
			return
		}
		defer fileResp.Body.Close()
		imgFileByte, err := ioutil.ReadAll(fileResp.Body)
		if err != nil {
			logger.Log.Errorf("??????????????????: %v", err.Error())
			return
		} else {
			// ????????????????????????
			contentType := http.DetectContentType(imgFileByte)
			fileName := fmt.Sprintf("%v_%v", ctx.MsgId, data.Appmsg.Title)
			if user, err := ctx.Bot.GetCurrentUser(); err == nil {
				uin := user.Uin
				fileName = fmt.Sprintf("%v/%v", uin, fileName)
			}
			// ????????????(reader2????????????????????????BUG,??????http.Response.Body??????????????????)
			reader2 := ioutil.NopCloser(bytes.NewReader(imgFileByte))
			flag := oss.SaveToOss(reader2, contentType, fileName)
			if flag {
				//fileUrl := fmt.Sprintf("https://%v/%v/%v", config.Config.OssConfig.Endpoint, config.Config.OssConfig.BucketName, fileName)
				fileUrl := fmt.Sprintf(config.Config.OssConfig.Endpoint, config.Config.OssConfig.BucketName, fileName)
				logger.Log.Infof("?????????????????????????????????: %v", fileUrl)
				ctx.Content = fileUrl
			} else {
				logger.Log.Error("??????????????????")
			}
		}
	}
	ctx.Next()
}

// ????????????????????????????????????
func checkIsExist(a []string, k string) bool {
	for _, item := range a {
		if item == k {
			return true
		}
	}
	return false
}
