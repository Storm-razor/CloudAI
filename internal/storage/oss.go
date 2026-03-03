package storage

import (
	"bytes"
	"fmt"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/wwwzy/CloudAI/config"
)

// OSSStorage 阿里云OSS存储驱动结构体
type OSSStorage struct {
	bucket *oss.Bucket // OSS Bucket实例
}

// ---------------------------
// @brief 初始化OSS存储
// ---------------------------
func NewOSSStorage(cfg config.OSSConfig) (*OSSStorage, error) {
	client, err := oss.New(cfg.Endpoint, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create OSS client: %v", err)
	}

	// 获取Bucket实例
	bucket, err := client.Bucket(cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to get OSS bucket: %v", err)
	}

	return &OSSStorage{bucket: bucket}, nil
}

// ---------------------------
// @brief 上传文件到OSS
// ---------------------------
func (s *OSSStorage) Upload(data []byte, key string, contentType string) error {
	return s.bucket.PutObject(key, bytes.NewReader(data))
}

// ---------------------------
// @brief 从OSS下载文件
// ---------------------------
func (s *OSSStorage) Download(key string) ([]byte, error) {
	reader, err := s.bucket.GetObject(key)
	if err != nil {
		return nil, fmt.Errorf("failed to download from OSS: %v", err)
	}
	defer reader.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(reader); err != nil {
		return nil, fmt.Errorf("failed to read OSS data: %v", err)
	}
	return buf.Bytes(), nil
}

// ---------------------------
// @brief 删除OSS文件
// ---------------------------
func (s *OSSStorage) Delete(key string) error {
	return s.bucket.DeleteObject(key)
}

// ---------------------------
// @brief 生成带签名的临时访问URL（有效期1小时）
// ---------------------------
func (s *OSSStorage) GetURL(key string) (string, error) {
	expired := time.Now().Add(1 * time.Hour)
	return s.bucket.SignURL(key, oss.HTTPGet, int64(expired.Unix()))
}
