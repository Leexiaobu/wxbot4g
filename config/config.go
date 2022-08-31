package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"wxbot4g/logger"
)

type config struct {
	RedisConfig redisConfig `yaml:"redis"`
	MySQLConfig mysqlConfig `yaml:"mysql"`
	OssConfig   ossConfig   `yaml:"oss"`
	ApiConfig   apiConfig   `yaml:"api"`
}

var Config = &config{}

func init() {
	file, err := os.ReadFile("conf.yaml")
	if err != nil {
		log.Fatal("fail to read file:", err)
	}
	err = yaml.Unmarshal(file, &Config)
	if err != nil {
		logger.Log.Errorln("error init config,%s", err.Error())
	}
}

// Redis配置
type redisConfig struct {
	Host      string `yaml:"host"`       // Redis主机
	Port      string `yaml:"port"`       // Redis端口
	Password  string `yaml:"password"`   // Redis密码
	Db        int    `yaml:"dbname"`     // Redis库
	KeepAlive int    `yaml:"keep_alive"` // Redis库
}

// MySQL配置
type mysqlConfig struct {
	Host     string `yaml:"host"`     // 主机
	Port     string `yaml:"port"`     // 端口
	Username string `yaml:"username"` // 用户名
	Password string `yaml:"password"` // 密码
	DbName   string `yaml:"dbname"`   // 数据库名称
}

type ossConfig struct {
	Endpoint        string `yaml:"endpoint"`          // 接口地址
	AccessKeyID     string `yaml:"access_key_id"`     // 账号
	SecretAccessKey string `yaml:"secret_access_key"` // 密码
	BucketName      string `yaml:"bucket_name"`       // 桶名称
	UseSsl          bool   `yaml:"use_ssl"`           // 是否使用SSL
}
type apiConfig struct {
	XiaoAiConfig  xiaoaiConfig  `yaml:"xiaoai"`
	WeatherConfig weatherConfig `yaml:"weather"`
	Url           string        `yaml:"url"`
}

type xiaoaiConfig struct {
	Url string `yaml:"url"`
}
type weatherConfig struct {
	Url      string `yaml:"url"`
	Location string `yaml:"location"`
	Key      string `yaml:"key"`
}
