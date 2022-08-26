package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type config struct {
	RedisConfig redisConfig `yaml:"redis"`
	MySQLConfig mysqlConfig `yaml:"mysql"`
	OssConfig   ossConfig   `yaml:"oss"`
}

var Config = &config{}

func init() {
	file, err := os.ReadFile("conf.yml")
	if err != nil {
		log.Fatal("fail to read file:", err)
	}
	err = yaml.Unmarshal(file, &Config)
}

//// RedisConfig Redis配置
//var (
//	RedisConfig redisConfig
//	MySQLConfig mysqlConfig
//	OssConfig   ossConfig
//)

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

// InitRedisConfig 初始化Redis配置
//func InitRedisConfig() {
//	// RedisHost Redis主机
//	host := utils.GetEnvVal("REDIS_HOST", "192.168.31.192")
//	// RedisPort Redis端口
//	port := utils.GetEnvVal("REDIS_PORT", "6379")
//	// RedisPassword Redis密码
//	password := utils.GetEnvVal("REDIS_PWD", "")
//	// Redis库
//	db := utils.GetEnvIntVal("REDIS_DB", 0)
//
//	RedisConfig = redisConfig{
//		Host:     host,
//		Port:     port,
//		Password: password,
//		Db:       db,
//	}
//}
//func InitMysqlConfig() {
//	host := utils.GetEnvVal("MYSQL_HOST", "192.168.31.192")
//	port := utils.GetEnvVal("MYSQL_PORT", "3306")
//	password := utils.GetEnvVal("MYSQL_PWD", "root")
//	user := utils.GetEnvVal("MSYQL_USER", "root")
//	db := utils.GetEnvVal("MYSQL_DB", "wxbot4g")
//	MySQLConfig = mysqlConfig{
//		Host:     host,
//		Port:     port,
//		Username: user,
//		Password: password,
//		DbName:   db,
//	}
//}
//
//// InitOssConfig 初始化OSS配置
//func InitOssConfig() {
//	endpoint := utils.GetEnvVal("OSS_ENDPOINT", "192.168.31.192:9001")
//	accessKeyID := utils.GetEnvVal("OSS_KEY", "leexiaobu")
//	secretAccessKey := utils.GetEnvVal("OSS_SECRET", "MINIOdmm21~")
//	bucketName := utils.GetEnvVal("OSS_BUCKET", "wechat")
//	useSSL := utils.GetEnvBoolVal("OSS_SSL", false)
//
//	OssConfig = ossConfig{
//		Endpoint:        endpoint,
//		AccessKeyID:     accessKeyID,
//		SecretAccessKey: secretAccessKey,
//		BucketName:      bucketName,
//		UseSsl:          useSSL,
//	}
//}
