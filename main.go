package main

import (
	"sync"
	"wxbot4g/robot"
)

func main() {
	// 初始化WechatBotMap
	robot.InitWechatBotsMap()
	// 初始化OSS
	// 初始化Redis连接
	// 初始化Redis里登录的数据
	robot.InitBotWithStart()
	// 定时更新 Bot 的热登录数据
	robot.UpdateHotLoginData()
	// 保活
	//关闭保活
	//robot.KeepAliveHandle()

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
