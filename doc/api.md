# CloudAI API 文档

本文档描述了 CloudAI 项目的 API 请求与响应结构。所有 API 均遵循 RESTful 风格，请求和响应体均为 JSON 格式。

## 目录

- [用户认证 (User Auth)](#用户认证-user-auth)
- [模型管理 (Model Management)](#模型管理-model-management)
- [知识库管理 (Knowledge Base)](#知识库管理-knowledge-base)
- [文件管理 (File Management)](#文件管理-file-management)
- [智能体管理 (Agent Management)](#智能体管理-agent-management)
- [对话系统 (Conversation)](#对话系统-conversation)

---

## 用户认证 (User Auth)

### 注册
**Endpoint**: `POST /api/v1/register`
**Request Body**: `UserRegisterReq`
```json
{
  "username": "user123",      // 必填，3-20字符
  "password": "password123",  // 必填，6-30字符
  "email": "user@example.com",// 必填，邮箱格式
  "phone": "+8613800000000"   // 必填，E.164格式
}
```

### 用户名登录
**Endpoint**: `POST /api/v1/login/username`
**Request Body**: `UserNameLoginReq`
```json
{
  "username": "user123",
  "password": "password123"
}
```

### 手机号登录
**Endpoint**: `POST /api/v1/login/phone`
**Request Body**: `UserPhoneLogin`
```json
{
  "phone": "+8613800000000",
  "password": "password123"
}
```

### 登录响应
**Response Body**: `LoginResponse`
```json
{
  "access_token": "eyJhbGciOiJIUzI1Ni...", // JWT Token
  "expires_in": 7200,                      // 过期时间(秒)
  "token_type": "Bearer"
}
```

---

## 模型管理 (Model Management)

### 注册模型
**Endpoint**: `POST /api/v1/models`
**Request Body**: `CreateModelRequest`
```json
{
  "type": "llm",               // 必填，"embedding" 或 "llm"
  "name": "GPT-4 Turbo",       // 显示名称
  "server": "openai",          // 供应商
  "base_url": "https://api.openai.com/v1",
  "model": "gpt-4-turbo",      // 实际模型ID
  "api_key": "sk-...",
  "max_tokens": 128000,        // 上下文长度
  // LLM 特有
  "max_output_length": 4096,
  "function": true,            // 是否支持函数调用
  // Embedding 特有
  "dimension": 1536            // 向量维度
}
```

### 更新模型
**Endpoint**: `PUT /api/v1/models`
**Request Body**: `UpdateModelRequest`
```json
{
  "id": "model_uuid",          // 必填
  "name": "New Name",
  // 其他字段同创建请求，可选
}
```

### 模型列表
**Endpoint**: `GET /api/v1/models`
**Query Params**: `PageModelRequest`
- `type`: `llm` (可选)
- `page`: `1` (默认)
- `size`: `10` (默认)

---

## 知识库管理 (Knowledge Base)

### 创建知识库
**Endpoint**: `POST /api/v1/kbs`
**Request Body**: `CreateKBRequest`
```json
{
  "name": "公司文档库",          // 必填
  "description": "包含HR政策和技术文档",
  "embed_model_id": "model_id" // 必填，关联的Embedding模型ID
}
```

### 添加文件到知识库
**Endpoint**: `POST /api/v1/kbs/files`
**Request Body**: `AddFileRequest`
```json
{
  "kb_id": "kb_uuid",          // 知识库ID
  "file_id": "file_uuid"       // 文件ID
}
```

### 批量删除文档
**Endpoint**: `DELETE /api/v1/kbs/docs`
**Request Body**: `BatchDeleteDocsReq`
```json
{
  "kb_id": "kb_uuid",
  "doc_ids": ["doc_id_1", "doc_id_2"]
}
```

### 检索测试
**Endpoint**: `POST /api/v1/kbs/retrieve`
**Request Body**: `RetrieveRequest`
```json
{
  "kb_id": "kb_uuid",
  "query": "请假流程是什么？",
  "top_k": 3
}
```

---

## 文件管理 (File Management)

### 创建文件夹
**Endpoint**: `POST /api/v1/files/folder`
**Request Body**: `CreateFolderReq`
```json
{
  "name": "新文件夹",
  "parent_id": "parent_uuid" // 可选，根目录不传
}
```

### 批量移动
**Endpoint**: `POST /api/v1/files/move`
**Request Body**: `BatchMoveRequest`
```json
{
  "files_pid": ["file_id_1", "folder_id_2"],
  "target_pid": "folder_uuid"
}
```

### 重命名
**Endpoint**: `POST /api/v1/files/rename`
**Request Body**: `RenameRequest`
```json
{
  "file_id": "file_uuid",
  "new_name": "renamed.pdf"
}
```

---

## 智能体管理 (Agent Management)

### 创建智能体
**Endpoint**: `POST /api/v1/agents`
**Request Body**: `CreateAgentRequest`
```json
{
  "name": "翻译助手",
  "description": "精通多国语言翻译"
}
```

### 更新智能体配置
**Endpoint**: `PUT /api/v1/agents`
**Request Body**: `UpdateAgentRequest`
```json
{
  "id": "agent_uuid",
  "name": "高级翻译助手",
  "prompt": "你是一个专业的翻译...",
  "llm_config": { ... },       // 模型参数
  "knowledge": { ... },        // 知识库配置
  "tools": { ... }             // 工具配置
}
```

### 执行智能体 (API调用)
**Endpoint**: `POST /api/v1/agents/execute`
**Request Body**: `ExecuteAgentRequest`
```json
{
  "id": "request_uuid",
  "agent_id": "agent_uuid",
  "message": {
    "query": "Translate this to English",
    "history": []
  }
}
```

---

## 对话系统 (Conversation)

### 创建会话
**Endpoint**: `POST /api/v1/conversations`
**Request Body**: `CreateConvRequest`
```json
{
  "agent_id": "agent_uuid"
}
```

### 发送消息 (对话)
**Endpoint**: `POST /api/v1/chat`
**Request Body**: `ChatRequest` (单轮/简单) 或 `ConvRequest` (多轮/复杂)
```json
{
  "agent_id": "agent_uuid",
  "conv_id": "conv_uuid",
  "message": "你好"
}
```

### 对话响应 (普通)
**Response Body**: `ChatResponse`
```json
{
  "response": "你好！有什么我可以帮你的吗？",
  "references": [ ... ] // 引用文档列表
}
```

### 对话响应 (流式 SSE)
**Response Body**: `ChatStreamResponse`
```json
{
  "id": "msg_uuid",
  "object": "chat.completion.chunk",
  "choices": [
    {
      "delta": { "content": "你好" },
      "finish_reason": null
    }
  ]
}
```
