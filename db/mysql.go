package db

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
	. "wxbot4g/config"
	"wxbot4g/logger"
)

var MysqlCon mysqlCon

type mysqlCon struct {
	client *gorm.DB
}

func InitMysqlCon() {
	logger.Log.Info(fmt.Sprintf("Mysql:%s:%s/%s ", Config.MySQLConfig.Host, Config.MySQLConfig.Port, Config.MySQLConfig.DbName))
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Config.MySQLConfig.Username,
		Config.MySQLConfig.Password,
		Config.MySQLConfig.Host,
		Config.MySQLConfig.Port,
		Config.MySQLConfig.DbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	} else {
		logger.Log.Info("mysql init success")
		MysqlCon = mysqlCon{
			client: db,
		}
	}
}
func (m *mysqlCon) Save(data interface{}) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // 在调用WithTimeout之后defer cancel()
	err := MysqlCon.client.Create(data).WithContext(ctx).Error

	if err != nil {
		logger.Log.Errorf("保存数据到数据库失败: %v", err.Error())
		return false
	}
	return true
}
