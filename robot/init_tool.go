package robot

import (
	. "wxbot4g/db"
	"wxbot4g/logger"
	"wxbot4g/protocol"
)

// InitBotWithStart 系统启动的时候从Redis加载登录信息自动登录
func InitBotWithStart() {
	logger.Log.Infof("登录管理员账户")
	key := "wxboot:login:leexiaobu"
	// 提取出AppKey
	appKey := key[13:]
	// 调用热登录
	logger.Log.Debugf("当前热登录AppKey: %v", appKey)
	bot := InitWechatBotHandle()
	storage := protocol.NewRedisHotReloadStorage(key)
	if err := bot.HotLogin(storage, true); err != nil {
		logger.Log.Infof("[%v] 热登录失败，错误信息：%v", appKey, err.Error())
		// 登录失败，删除热登录数据
		if err := RedisClient.Del(key); err != nil {
			logger.Log.Errorf("[%v] Redis缓存删除失败，错误信息：%v", key, err.Error())
		}
	}
	loginUser, _ := bot.GetCurrentUser()
	logger.Log.Infof("[%v]初始化自动登录成功，用户名：%v", appKey, loginUser.NickName)
	// 登录成功，写入到WechatBots
	SetBot(appKey, bot)
}
