package model

import "encoding/json"

// Conversation 对话表
type Conversation struct {
	ID         uint64          `gorm:"primaryKey;column:id"`
	ConvID     string          `gorm:"uniqueIndex;column:conv_id;type:varchar(255)"`
	UserID     uint            `gorm:"index;column:user_id"`
	AgentID    string          `gorm:"index;column:agent_id;type:varchar(255)"`
	Title      string          `gorm:"column:title;type:varchar(255)"`
	CreatedAt  int64           `gorm:"column:created_at"`
	UpdatedAt  int64           `gorm:"column:updated_at"`
	Settings   json.RawMessage `gorm:"column:settings;type:json"`
	IsArchived bool            `gorm:"column:is_archived;default:0"`
	IsPinned   bool            `gorm:"column:is_pinned;default:0"`
}

// TableName 设置表名
func (Conversation) TableName() string {
	return "conversations"
}

// Message 消息表
type Message struct {
	ID            uint64          `gorm:"primaryKey;column:id"`
	MsgID         string          `gorm:"uniqueIndex;column:msg_id;type:varchar(255)"`                      // 消息的业务唯一标识(UUID)，用于前端引用
	UserID        uint            `gorm:"index;column:user_id"`                                             // 发送消息的用户ID
	ConvID        string          `gorm:"column:conv_id;type:varchar(255)"`                                 // 所属会话ID(关联Conversation.ConvID)
	ParentID      string          `gorm:"column:parent_id;type:varchar(255);default:''"`                    // 父消息ID，用于构建消息树(支持分支对话/重新生成)
	Role          string          `gorm:"column:role;type:enum('user','assistant','system','function')"`    // 角色：user(用户), assistant(AI), system(系统指令)
	Content       string          `gorm:"column:content;type:text"`                                         // 消息内容(Markdown格式)
	CreatedAt     int64           `gorm:"column:created_at"`                                                // 创建时间戳(秒/毫秒)
	OrderSeq      int             `gorm:"column:order_seq;default:0"`                                       // 消息在当前分支中的顺序序号
	TokenCount    int             `gorm:"column:token_count;default:0"`                                     // 该消息消耗的Token数量(用于计费/统计)
	Status        string          `gorm:"column:status;type:enum('sent','pending','error');default:'sent'"` // 状态：sent(已发送/完成), pending(生成中), error(失败)
	Metadata      json.RawMessage `gorm:"column:metadata;type:json"`                                        // 元数据(JSON)：存储额外信息，如RAG引用的文档、消耗时间、模型参数等
	IsContextEdge bool            `gorm:"column:is_context_edge;default:0"`                                 // 是否为上下文边界(用于长对话截断，标记从这里开始不带入历史)
	IsVariant     bool            `gorm:"column:is_variant;default:0"`                                      // 是否为变体(例如用户重新编辑了问题，或者AI重新生成了回答)
}

// TableName 设置表名
func (Message) TableName() string {
	return "messages"
}

// Attachment 附件表
type Attachment struct {
	ID             uint64 `gorm:"primaryKey;column:id"`
	AttachID       string `gorm:"uniqueIndex;column:attach_id;type:varchar(255)"`                          // 附件唯一业务ID
	UserID         uint   `gorm:"index;column:user_id"`                                                    // 上传者ID
	MessageID      string `gorm:"column:message_id;type:varchar(255)"`                                     // 关联的消息ID(Attachment依附于Message)
	AttachmentType string `gorm:"column:attachment_type;type:enum('file','image','code','audio','video')"` // 附件类型
	FileName       string `gorm:"column:file_name;type:varchar(255)"`                                      // 原始文件名
	FileSize       int64  `gorm:"column:file_size"`                                                        // 文件大小(字节)
	StorageType    string `gorm:"column:storage_type;type:enum('path','blob','cloud')"`                    // 存储方式：path(本地路径), blob(数据库二进制), cloud(OSS/MinIO)
	StoragePath    string `gorm:"column:storage_path;type:varchar(1024)"`                                  // 存储路径或URL
	Thumbnail      []byte `gorm:"column:thumbnail;type:mediumblob"`                                        // 缩略图数据(仅图片/视频有)，直接存二进制方便快速预览
	Vectorized     bool   `gorm:"column:vectorized;default:0"`                                             // 是否已向量化(RAG用)：如果是文档/代码，是否已经切片存入Milvus
	DataSummary    string `gorm:"column:data_summary;type:text"`                                           // 内容摘要(AI生成的总结)，用于快速了解附件内容
	MimeType       string `gorm:"column:mime_type;type:varchar(255)"`                                      // 具体的MIME类型，如 application/pdf, image/png
	CreatedAt      int64  `gorm:"column:created_at"`                                                       // 上传时间
}

// TableName 设置表名
func (Attachment) TableName() string {
	return "attachments"
}

// MessageAttachment 消息附件关联表
type MessageAttachment struct {
	ID           uint64 `gorm:"primaryKey;column:id"`
	MessageID    string `gorm:"column:message_id;type:varchar(255)"`
	AttachmentID string `gorm:"column:attachment_id;type:varchar(255)"`
}

// TableName 设置表名
func (MessageAttachment) TableName() string {
	return "message_attachments"
}
