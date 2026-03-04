# AI-Cloud-Go 项目指南

本项目是一个基于 Go 语言的 AI 云服务后端，集成了 RAG（检索增强生成）、向量数据库（Milvus）、对象存储（MinIO）以及大模型（LLM）调用能力。

## 1. 项目概览与技术栈

### 核心技术栈
- **语言**: Go 1.23+
- **Web框架**: Gin
- **数据库**: MySQL 8.0 (元数据), Milvus 2.4 (向量数据)
- **对象存储**: MinIO (文件存储)
- **配置管理**: Viper
- **ORM**: Gorm
- **AI框架**: CloudWeGo/eino (用于构建 AI 应用)
- **LLM支持**: OpenAI, Ollama, DeepSeek

### 目录结构规划
```text
CloudAI/
├── cmd/
│   └── main.go           # 程序入口
├── config/               # 配置加载与定义
├── internal/
│   ├── component/        # AI组件 (Embedding, LLM, Indexer, Parser)
│   ├── controller/       # HTTP 接口处理
│   ├── dao/              # 数据库访问层 (Data Access Object)
│   ├── database/         # 数据库初始化 (MySQL, Milvus)
│   ├── middleware/       # 中间件 (Auth, CORS)
│   ├── model/            # 数据库模型定义
│   ├── router/           # 路由定义
│   ├── service/          # 业务逻辑层
│   └── utils/            # 工具函数
├── pkgs/                 # 公共包 (Error, Response)
├── docker-compose.yml    # 基础设施编排
└── go.mod                # 依赖管理
```

---

## 2. 环境准备与基础设施 (Phase 1)

### 2.1 初始化项目
在 `d:\Github\golang\src\CloudAI` 目录下执行：
```bash
go mod init cloud-ai
```

### 2.2 搭建基础设施
创建 `docker-compose.yml` 文件，定义 MySQL, MinIO, Milvus 服务。
*参考原项目 `d:\Github\golang\src\AI-Cloud-Go\docker-compose.yml`*

**关键配置点**:
- **MySQL**: 创建 `ai_cloud` 数据库。
- **MinIO**: 设置 buckets (`ai-cloud`)。
- **Milvus**: 向量数据库，依赖 Etcd 和 MinIO。

启动服务：
```bash
docker-compose up -d
```

---

## 3. 基础框架搭建 (Phase 2)

### 3.1 配置模块 (`config/`)
使用 `viper` 读取 `config.yaml`。
- 定义 `Config` 结构体，映射 yaml 字段。
- 实现 `InitConfig()` 函数加载配置。

### 3.2 数据库连接 (`internal/database/`)
- **MySQL**: 使用 `gorm` 连接 MySQL。
- **Milvus**: 使用 `milvus-sdk-go` 连接 Milvus。

### 3.3 公共组件 (`pkgs/`)
- **Response**: 统一 API 返回格式 (`code`, `msg`, `data`)。
- **Errcode**: 定义业务错误码。

### 3.4 数据模型 (`internal/model/`)
定义核心 Struct，并使用 Gorm Tag：
- `User`: 用户信息 (ID, Username, Password, ...).
- `File`: 文件记录 (Name, Size, URL, ...).
- `KnowledgeBase`: 知识库元数据.
- `Document`: 知识库中的文档.
- `Model`: LLM 模型配置信息.
- `Agent`: 智能体配置.

---

## 4. 核心业务功能 (Phase 3)

### 4.1 用户模块
- **Dao**: `UserDao` 实现 CRUD。
- **Service**: `UserService` 实现注册、登录（JWT签发）。
- **Controller**: `UserController` 处理 HTTP 请求。

### 4.2 文件模块(需要补全流式传输)
- **存储接口**: 定义 `Storage` 接口，实现 `MinIO` 和 `Local` 两种策略。
- **上传流程**: 接收文件 -> 上传 MinIO -> 记录数据库 `File` 表。

### 4.3 认证中间件 (`internal/middleware/`)
- 实现 `JWTAuth` 中间件，拦截 `/api` 请求，验证 Token。
- 实现 `CORS` 中间件，处理跨域。

---

## 5. AI 引擎与 RAG 实现 (Phase 4 - 核心难点)

### 5.1 LLM 工厂 (`internal/component/llm/`)
基于 `cloudwego/eino` 或直接封装 SDK。
- 实现 `GetLLMClient` 工厂方法。
- 支持 `OpenAI` 和 `Ollama` 适配器。

### 5.2 Agent组件 (`internal/component/`)
- **Embedding**: 文本向量化组件 (OpenAI, Ollama, DeepSeek)。
- **Indexer**: 向量索引组件 (Milvus)。
- **LLM**: 大模型调用组件 (OpenAI, Ollama, DeepSeek)。
- **Parser**: 文档解析组件 (docconv, PDF/TXT)(TODO: 支持更多格式,降低内存占用)。
- **Retriever**: 向量检索组件 (Milvus)。

### 5.3 知识库服务 (`internal/service/kb_service.go`)
**RAG 流程**:
1.  **解析**: 使用 `docconv` 或其他库解析 PDF/TXT 内容。
2.  **切片 (Chunking)**: 将长文本切分为小段 (Chunk Size ~500-1000)。
3.  **向量化**: 调用 Embedding 组件获取向量。
4.  **存储**: 将 文本 + 向量 存入 Milvus。

### 5.4 检索服务
- 实现 `Retrieve` 方法：Query -> Embedding -> Milvus Search -> 返回相关 Chunks。

---

## 6. 智能体与对话 (Phase 5)

### 6.1 Agent 模块
- 定义 `Agent` 模型，包含 Prompt 模板、关联的知识库 ID、使用的模型 ID。
- 实现 `AgentService`：组装 System Prompt + 用户问题 + 知识库上下文。

### 6.2 对话历史 (`internal/dao/history/`)
- 记录 `Conversation` (会话) 和 `Message` (消息)。
- 在调用 LLM 前，查询历史消息作为 Context 传入。

### 6.3 流式响应
- Controller 中使用 `c.Stream()` 或 SSE (Server-Sent Events) 返回 LLM 的流式输出。

---

## 7. 接口组装 (Phase 6)

### 7.1 路由 (`internal/router/`)
- 定义 `/api/users`, `/api/files`, `/api/knowledge`, `/api/chat` 等路由组。
- 绑定对应的 Controller 方法。

### 7.2 入口文件 (`cmd/main.go`)
1.  初始化 Config。
2.  初始化 MySQL & Milvus 连接。
3.  依赖注入：初始化所有 DAO -> Service -> Controller。
4.  启动 Gin Server。

---

## 8. 验证与测试
1.  注册登录，获取 Token。
2.  上传 PDF 文件。
3.  创建知识库，关联文件（触发解析入库）。
4.  创建 Agent，关联知识库。
5.  发起对话，验证是否能基于知识库回答。

