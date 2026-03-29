package oss

//实现上传文件、生成预签名 URL

import (
	"bytes"

	"geekedu-project/common/config"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// 初始化OSS客户端
func NewOSSClient() (*oss.Client, *oss.Bucket, error) {
	cfg := config.InitConfig()
	// 创建OSS客户端
	client, err := oss.New(cfg.OSSEndpoint, cfg.OSSAccessKey, cfg.OSSSecretKey)
	if err != nil {
		return nil, nil, err
	}

	// 获取Bucket实例
	bucket, err := client.Bucket(cfg.OSSBucket)
	if err != nil {
		return nil, nil, err
	}

	return client, bucket, nil
}

// 上传文件到OSS（封面/视频）
func UploadFile(objectKey string, fileBytes []byte) error {
	_, bucket, err := NewOSSClient()
	if err != nil {
		return err
	}

	// 上传字节流到OSS
	err = bucket.PutObject(objectKey, bytes.NewReader(fileBytes))
	if err != nil {
		return err
	}
	return nil
}

// 生成预签名URL（用于播放视频，有效期3600秒）
func GeneratePresignedURL(objectKey string) (string, error) {
	_, bucket, err := NewOSSClient()
	if err != nil {
		return "", err
	}

	// 设置URL有效期1小时
	expireSeconds := int64(3600)
	signedURL, err := bucket.SignURL(objectKey, oss.HTTPGet, expireSeconds)
	if err != nil {
		return "", err
	}

	return signedURL, nil
}
