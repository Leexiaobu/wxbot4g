package holiday

import (
	"testing"
	"wxbot4g/db"
	"wxbot4g/logger"
)

func TestUpdateYear(t *testing.T) {
	logger.InitLogger()
	db.InitMysqlCon()
	// 初始化WechatBotMap
	UpdateYear("2022")
}
