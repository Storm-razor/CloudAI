package database

import (
	"context"
	"sync"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/wwwzy/CloudAI/config"
)

var (
	milvusOnce   sync.Once
	milvusErr    error
	milvusClient client.Client
)

// ---------------------------
// @brief Milvus初始化
// ---------------------------
func InitMilvus(ctx context.Context) (client.Client, error) {
	milvusOnce.Do(func() {
		milvusClient, milvusErr = client.NewClient(ctx, client.Config{
			Address: config.GetConfig().Milvus.Address,
		})
	})
	return milvusClient, milvusErr
}

// ---------------------------
// @brief 获取Milvus单例
// ---------------------------
func GetMilvusClient() client.Client {
	return milvusClient
}
