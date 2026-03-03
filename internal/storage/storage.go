package storage

import (
	"fmt"

	"github.com/wwwzy/CloudAI/config"
)

// 存储驱动接口
type Driver interface {
	Upload(data []byte, key string, contentType string) error // 上传文件
	Download(key string) ([]byte, error)                      // 下载文件
	Delete(key string) error                                  // 删除文件
	GetURL(key string) (string, error)                        // 获取访问URL
}

func NewDriver(cfg config.StorageConfig) (Driver, error) {
	switch cfg.Type {
	case "local":
		return NewLocalStorage(cfg.Local.BaseDir)
	case "oss":
		return NewOSSStorage(cfg.OSS)
	case "minio":
		return NewMinioStorage(cfg.Minio)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.Type)
	}
}
