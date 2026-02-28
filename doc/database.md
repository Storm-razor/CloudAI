# 数据库模型设计文档 (Database Schema)

## 1. 用户模块 (User)

**表名**: `users` (默认)
**描述**: 存储系统的用户信息，支持用户名、手机号和邮箱登录。

| 字段名 | 类型 | GORM 标签/约束 | 说明 |
| :--- | :--- | :--- | :--- |
| `ID` | `uint` | `primaryKey` | 用户唯一标识 (自增主键) |
| `Username` | `string` | `uniqueIndex;size:50;not null` | 用户名 (唯一，必填) |
| `Phone` | `string` | `uniqueIndex;size:20;not null` | 手机号 (唯一，必填，E.164格式) |
| `Email` | `string` | `uniqueIndex;size:100` | 邮箱地址 (唯一) |
| `Password` | `string` | `not null` | 加密后的密码 |
| `CreatedAt` | `time.Time` | `autoCreateTime` | 账户创建时间 |
| `UpdatedAt` | `time.Time` | `autoUpdateTime` | 账户最后更新时间 |

## 2. 模型管理 (Model)

**表名**: `models` (默认)
**描述**: 管理接入的外部 AI 模型配置，包括 LLM (大语言模型) 和 Embedding (向量模型)。

| 字段名 | 类型 | GORM 标签/约束 | 说明 |
| :--- | :--- | :--- | :--- |
| `ID` | `string` | `primaryKey;type:char(36)` | 模型唯一标识 (UUID) |
| `UserID` | `uint` | `index` | 创建该配置的用户 ID |
| `Type` | `string` | `not null` | 模型类型: `embedding` 或 `llm` |
| `ShowName` | `string` | `not null` | 在前端显示的名称 (如 "GPT-4 Turbo") |
| `Server` | `string` | `not null` | 供应商标识 (如 `openai`, `ollama`) |
| `BaseURL` | `string` | `not null` | API 基础地址 (如 `https://api.openai.com/v1`) |
| `ModelName` | `string` | `not null` | 实际调用的模型标识符 (如 `gpt-4-turbo`) |
| `APIKey` | `string` | - | API 访问密钥 (Ollama 可为空) |
| `Dimension` | `int` | - | 向量维度 (仅 Embedding 模型有效) |
| `MaxOutputLength`| `int` | `default:4096` | 最大输出长度 (仅 LLM 有效) |
| `Function` | `bool` | `default:false` | 是否支持 Function Calling (仅 LLM 有效) |
| `MaxTokens` | `int` | `default:1024` | 允许的最大输入 Token 限制 |

## 3. 智能体模块 (Agent)

**表名**: `agents` (默认)
**描述**: 存储智能体的配置信息，每个 Agent 包含特定的 Prompt、关联的模型和知识库等。

| 字段名 | 类型 | GORM 标签/约束 | 说明 |
| :--- | :--- | :--- | :--- |
| `ID` | `string` | `primaryKey;type:char(36)` | 智能体唯一标识 (UUID) |
| `UserID` | `uint` | `index` | 创建者的用户 ID |
| `Name` | `string` | `not null` | 智能体名称 |
| `Description` | `string` | `type:text` | 智能体描述 |
| `AgentSchema` | `string` | `type:json` | **核心配置字段** (JSON格式)，包含 LLM配置、知识库关联、Prompt 等 |
| `CreatedAt` | `time.Time` | `autoCreateTime` | 创建时间 |
| `UpdatedAt` | `time.Time` | `autoUpdateTime` | 更新时间 |

> **AgentSchema 结构说明**:
> 包含 `LLMConfig` (模型参数), `MCP` (MCP服务器), `Tools` (工具集), `Prompt` (系统提示词), `Knowledge` (关联知识库ID) 等详细配置。

## 4. 知识库模块 (Knowledge & RAG)

### 4.1 知识库 (KnowledgeBase)
**表名**: `knowledge_bases` (默认)
**描述**: 知识库元数据，用于逻辑上分组管理文档。

| 字段名 | 类型 | GORM 标签/约束 | 说明 |
| :--- | :--- | :--- | :--- |
| `ID` | `string` | `primaryKey;type:char(36)` | 知识库 ID (UUID) |
| `Name` | `string` | `not null` | 知识库名称 |
| `Description` | `string` | - | 描述信息 |
| `UserID` | `uint` | `index` | 所属用户 ID |
| `EmbedModelID` | `string` | `index` | 关联的 Embedding 模型配置 ID |
| `MilvusCollection`| `string` | `not null` | 对应的 Milvus 向量数据库 Collection 名称 |

### 4.2 文档 (Document)
**表名**: `documents` (默认)
**描述**: 存储上传到知识库的具体文档记录。

| 字段名 | 类型 | GORM 标签/约束 | 说明 |
| :--- | :--- | :--- | :--- |
| `ID` | `string` | `primaryKey;type:char(36)` | 文档 ID (UUID) |
| `UserID` | `uint` | `index` | 上传用户 ID |
| `KnowledgeBaseID`| `string` | `index` | 所属知识库 ID |
| `FileID` | `string` | `index` | 关联的原始文件 ID (对应文件系统) |
| `Title` | `string` | - | 文档标题 |
| `DocType` | `string` | - | 文档类型 (pdf, txt, md 等) |
| `Status` | `int` | - | 处理状态: 0(待处理), 1(处理中), 2(已完成), 3(失败) |

> **注意**: 文档的具体内容切片 (Chunks) 和向量数据存储在 **Milvus** 向量数据库中，不存储在 MySQL。

## 5. 文件系统 (File)

**表名**: `files` (默认)
**描述**: 统一文件管理系统，支持本地存储和对象存储 (OSS/MinIO)。

| 字段名 | 类型 | GORM 标签/约束 | 说明 |
| :--- | :--- | :--- | :--- |
| `ID` | `string` | `primaryKey;type:char(36)` | 文件 ID (UUID) |
| `UserID` | `uint` | `index` | 上传用户 ID |
| `Name` | `string` | `not null` | 文件名 |
| `Size` | `int64` | - | 文件大小 (字节) |
| `Hash` | `string` | `index;size:64` | 文件 SHA-256 哈希值 (用于秒传/去重) |
| `MIMEType` | `string` | - | 文件 MIME 类型 |
| `IsDir` | `bool` | `default:false` | 是否为文件夹 |
| `ParentID` | `*string` | `type:char(36);index` | 父目录 ID (支持嵌套目录) |
| `StorageType` | `string` | `default:'local'` | 存储类型: `local` 或 `oss` |
| `StorageKey` | `string` | - | 文件在存储系统中的唯一标识 (路径或 Key) |

## 6. 对话系统 (Conversation)

### 6.1 会话 (Conversation)
**表名**: `conversations`
**描述**: 记录用户与 Agent 的一次完整对话会话。

| 字段名 | 类型 | GORM 标签/约束 | 说明 |
| :--- | :--- | :--- | :--- |
| `ID` | `uint64` | `primaryKey` | 内部主键 |
| `ConvID` | `string` | `uniqueIndex;type:varchar(255)` | 会话业务 ID |
| `UserID` | `uint` | `index` | 用户 ID |
| `AgentID` | `string` | `index;type:varchar(255)` | 关联的 Agent ID |
| `Title` | `string` | `type:varchar(255)` | 会话标题 |
| `Settings` | `json` | `type:json` | 会话特定设置 (JSON) |
| `IsArchived` | `bool` | `default:0` | 是否归档 |
| `IsPinned` | `bool` | `default:0` | 是否置顶 |

### 6.2 消息 (Message)
**表名**: `messages`
**描述**: 会话中的每一条消息记录。

| 字段名 | 类型 | GORM 标签/约束 | 说明 |
| :--- | :--- | :--- | :--- |
| `ID` | `uint64` | `primaryKey` | 内部主键 |
| `MsgID` | `string` | `uniqueIndex;type:varchar(255)` | 消息业务 ID |
| `UserID` | `uint` | `index` | 发送者 ID |
| `ConvID` | `string` | `type:varchar(255)` | 所属会话 ID |
| `ParentID` | `string` | `default:''` | 父消息 ID (用于树状对话结构) |
| `Role` | `string` | `enum(...)` | 角色: `user`, `assistant`, `system`, `function` |
| `Content` | `string` | `type:text` | 消息内容 |
| `TokenCount` | `int` | `default:0` | 消息 Token 数 |
| `Status` | `string` | `enum(...)` | 状态: `sent`, `pending`, `error` |
| `Metadata` | `json` | `type:json` | 元数据 (JSON) |

### 6.3 附件 (Attachment)
**表名**: `attachments`
**描述**: 消息中携带的附件 (图片、文件等)。

| 字段名 | 类型 | GORM 标签/约束 | 说明 |
| :--- | :--- | :--- | :--- |
| `ID` | `uint64` | `primaryKey` | 内部主键 |
| `AttachID` | `string` | `uniqueIndex` | 附件业务 ID |
| `MessageID` | `string` | - | 关联的消息 ID |
| `AttachmentType`| `string` | `enum(...)` | 类型: `file`, `image`, `code` 等 |
| `StoragePath` | `string` | `type:varchar(1024)` | 存储路径 |
| `Vectorized` | `bool` | `default:0` | 是否已向量化 (用于检索) |

### 6.4 消息附件关联 (MessageAttachment)
**表名**: `message_attachments`
**描述**: 消息与附件的多对多关联表。

| 字段名 | 类型 | GORM 标签/约束 | 说明 |
| :--- | :--- | :--- | :--- |
| `ID` | `uint64` | `primaryKey` | 主键 |
| `MessageID` | `string` | - | 消息 ID |
| `AttachmentID` | `string` | - | 附件 ID |
