package ginmodel

// AddFileRequest 添加文件到知识库请求
// 将已上传的文件关联到指定知识库
type AddFileRequest struct {
	FileID string `json:"file_id" binding:"required"` // 文件系统中的文件ID
	KBID   string `json:"kb_id" binding:"required"`   // 目标知识库ID
}

// BatchDeleteDocsReq 批量删除文档请求
// 从知识库中移除文档(同时删除向量数据)
type BatchDeleteDocsReq struct {
	KBID   string   `json:"kb_id" binding:"required"`   // 知识库ID
	DocIDs []string `json:"doc_ids" binding:"required"` // 文档ID列表
}

// CreateKBRequest 创建知识库请求
type CreateKBRequest struct {
	Name         string `json:"name" binding:"required"`       // 知识库名称
	Description  string `json:"description"`                   // 描述
	EmbedModelID string `json:"embed_model_id" binding:"required"` // 使用的Embedding模型ID(创建后不可修改)
}

// RetrieveRequest 检索测试请求
// 用于在知识库管理界面测试召回效果
type RetrieveRequest struct {
	KBID  string `json:"kb_id" binding:"required"` // 知识库ID
	Query string `json:"query" binding:"required"` // 测试查询语句
	TopK  int    `json:"top_k,default=3"`          // 返回数量
}
