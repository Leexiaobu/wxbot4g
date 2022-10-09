package robot

import (
	. "wxbot4g/db"
	"wxbot4g/logger"
	"wxbot4g/protocol"
)

// InitBotWithStart 系统启动的时候从Redis加载登录信息自动登录
func InitBotWithStart() {
	//keys, err := RedisClient.GetKeys("wechat:login:*")
	//if err != nil {
	//	logger.Log.Error("获取Key失败")
	//	return
	//}
	logger.Log.Infof("登录管理员账户")
	appKey := "wxboot:login:leexiaobu"
	bot := InitWechatBotHandle()
	storage := protocol.NewRedisHotReloadStorage(appKey)
	if err := bot.HotLogin(storage, true); err != nil {
		logger.Log.Infof("[%v] 热登录失败，错误信息：%v", appKey, err.Error())
		// 登录失败，删除热登录数据
		if err := RedisClient.Del(appKey); err != nil {
			logger.Log.Errorf("[%v] Redis缓存删除失败，错误信息：%v", appKey, err.Error())
		}
	}
	loginUser, _ := bot.GetCurrentUser()
	logger.Log.Infof("[%v]初始化自动登录成功，用户名：%v", appKey, loginUser.NickName)
	// 登录成功，写入到WechatBots
	SetBot(appKey, bot)
}
