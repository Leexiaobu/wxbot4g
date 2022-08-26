package oss

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	. "wxbot4g/config"
	"wxbot4g/logger"
)

var minioClient *minio.Client

// InitOssConnHandle 初始化OSS连接
func InitOssConnHandle() {
	ctx := context.Background()
	// 初始化OSS配置
	//config.InitOssConfig()
	// 初使化 minio client对象。
	client, err := minio.New(Config.OssConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(Config.OssConfig.AccessKeyID, Config.OssConfig.SecretAccessKey, ""),
		Secure: Config.OssConfig.UseSsl,
	})
	if err != nil {
		logger.Log.Panicf("OSS初始化失败: %v", err.Error())
	}
	logger.Log.Info("OSS连接成功，开始判断桶是否存在")
	// 判断捅是否存在，不存在就创建
	exists, err := client.BucketExists(ctx, Config.OssConfig.BucketName)
	if err != nil {
		logger.Log.Panicf("判断桶失败: %v", err)
	}
	if !exists {
		logger.Log.Info("桶不存在，开始创建")
		// 创建桶
		err = client.MakeBucket(ctx, Config.OssConfig.BucketName, minio.MakeBucketOptions{Region: "us-east-1"})
		if err != nil {
			logger.Log.Panicf("OSS桶创建失败: %v", err.Error())
		}
		logger.Log.Info("桶创建成功")
	} else {
		logger.Log.Info("桶已存在")
	}
	minioClient = client
	logger.Log.Info("OSS初始化成功")
}

// SaveToOss 保存文件到OSS
func SaveToOss(b io.Reader, contentType, fileName string) bool {
	logger.Log.Debugf("开始上传文件: %v", fileName)
	ctx := context.Background()
	_, err := minioClient.PutObject(ctx, Config.OssConfig.BucketName, fileName, b, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		logger.Log.Errorf("文件上传错误: %v", err)
		return false
	}
	logger.Log.Debugf("文件上传完毕: %v", fileName)
	return true
}
