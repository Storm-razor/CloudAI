package config

import (
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

// AppConfig 应用配置
type AppConfig struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Storage  StorageConfig  `mapstructure:"storage"`
	CORS     CORSConfig     `mapstructure:"cors"`
	RAG      RAGConfig      `mapstructure:"rag"`
	LLM      LLMConfig      `mapstructure:"llm"`
	Milvus   MilvusConfig   `mapstructure:"milvus"`
}

// 服务器配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret          string `mapstructure:"secret"`
	ExpirationHours int    `mapstructure:"expiration_hours"`
}

// MinioConfig Minio配置
type MinioConfig struct {
	Endpoint        string `mapstructure:"endpoint"`          // MinIO服务地址 (如 localhost:9000)
	Bucket          string `mapstructure:"bucket"`            // 存储桶名称 (如 ai-cloud)
	AccessKeyID     string `mapstructure:"access_key_id"`     // 访问密钥ID (用户名)
	AccessKeySecret string `mapstructure:"access_key_secret"` // 访问密钥Secret (密码)
	UseSSL          bool   `mapstructure:"use_ssl"`           // 是否使用HTTPS连接
	Region          string `mapstructure:"region"`            // 区域 (可选，默认空)
}

// MilvusConfig Milvus向量数据库配置
type MilvusConfig struct {
	Address         string `mapstructure:"address"`          // Milvus连接地址 (如 localhost:19530)
	CollectionName  string `mapstructure:"collection_name"`  // 集合名称 (相当于SQL中的表名)
	VectorDimension int    `mapstructure:"vector_dimension"` // 向量维度 (如 1024, 必须与Embedding模型一致)
	IndexType       string `mapstructure:"index_type"`       // 索引类型 (如 IVF_FLAT, HNSW)
	MetricType      string `mapstructure:"metric_type"`      // 距离度量方式 (如 COSINE, L2, IP)
	Nlist           int    `mapstructure:"nlist"`            // 倒排列表大小 (影响索引构建速度和搜索精度)

	Nprobe int `mapstructure:"nprobe"` // 搜索时查找的倒排列表数量 (影响搜索速度和精度)

	// 字段长度限制 (Schema 定义)
	IDMaxLength      string `mapstructure:"id_max_length"`       // ID字段最大长度
	ContentMaxLength string `mapstructure:"content_max_length"`  // 内容字段最大长度 (Chunk文本内容)
	DocIDMaxLength   string `mapstructure:"doc_id_max_length"`   // 文档ID字段最大长度
	DocNameMaxLength string `mapstructure:"doc_name_max_length"` // 文档名称字段最大长度
	KbIDMaxLength    string `mapstructure:"kb_id_max_length"`    // 知识库ID字段最大长度
}

// ---------------------------
// @brief 获取向量度量类型
// ---------------------------
func (m *MilvusConfig) GetMetricType() entity.MetricType {
	var metricType entity.MetricType
	switch m.MetricType {
	case "L2":
		metricType = entity.L2 // 欧几里得距离：测量向量间的直线距离，适合图像特征等数值型向量
	case "IP":
		metricType = entity.IP // 内积距离：适合已归一化的向量，计算效率高
	default:
		metricType = entity.COSINE // 余弦相似度：测量向量方向的相似性，适合文本语义搜索
	}
	return metricType
}

// ---------------------------
// @brief 根据配置构建索引
// ---------------------------
func (m *MilvusConfig) GetMilvusIndex() (idx entity.Index, err error) {
	metricType := m.GetMetricType()
	if m.Nlist <= 0 {
		m.Nlist = 128
	}
	switch m.IndexType {
	case "IVF_FLAT":
		// IVF_FLAT: 倒排文件索引 + 原始向量存储
		// 优点：搜索精度高；缺点：内存占用较大
		// nlist: 聚类数量，值越大精度越高但速度越慢，通常设置为 sqrt(n) 到 4*sqrt(n)，其中n为向量数量
		idx, err = entity.NewIndexIvfFlat(metricType, m.Nlist)
	case "IVF_SQ8":
		// IVF_SQ8: 倒排文件索引 + 标量量化压缩存储（8位）
		// 优点：比IVF_FLAT节省内存；缺点：轻微精度损失
		// nlist: 与IVF_FLAT相同，根据数据规模调整
		idx, err = entity.NewIndexIvfSQ8(metricType, m.Nlist)
	case "HNSW":
		// HNSW: 层次可导航小世界图索引，高效且精确但内存占用大
		// M: 每个节点的最大边数，影响图的连通性和构建/查询性能
		//    - 值越大，构建越慢，内存占用越大，但查询越精确
		//    - 通常取值范围为8-64之间，默认值8在大多数场景下平衡了性能和精度
		// efConstruction: 构建索引时每层搜索的候选邻居数量
		//    - 值越大，构建越慢，索引质量越高
		//    - 通常取值范围为40-800，默认值40在大多数场景下表现良好
		// 注：这两个参数需要根据数据特性和性能要求综合调优，目前使用经验值
		idx, err = entity.NewIndexHNSW(metricType, 8, 40) // M=8, efConstruction=40
	default:
		// 默认使用IVF_FLAT，兼顾搜索精度和性能
		idx, err = entity.NewIndexIvfFlat(metricType, m.Nlist)
	}
	return
}

// StorageConfig 存储配置
type StorageConfig struct {
	Type  string      `mapstructure:"type"` // local/oss/minio
	Local LocalConfig `mapstructure:"local"`
	OSS   OSSConfig   `mapstructure:"oss"`
	Minio MinioConfig `mapstructure:"minio"`
}

// LocalConfig 本地存储配置
type LocalConfig struct {
	BaseDir string `mapstructure:"base_dir"` // 本地存储根目录（如 /data/storage）
}

// OSSConfig 阿里云OSS对象存储配置
type OSSConfig struct {
	Endpoint        string `mapstructure:"endpoint"`          // OSS访问域名 (如 oss-cn-hangzhou.aliyuncs.com)
	Bucket          string `mapstructure:"bucket"`            // 存储桶名称
	AccessKeyID     string `mapstructure:"access_key_id"`     // 阿里云 AccessKey ID
	AccessKeySecret string `mapstructure:"access_key_secret"` // 阿里云 AccessKey Secret
}

// CORSConfig 跨域资源共享(CORS)配置
// 用于解决前端(如Vue/React)与后端不同域名/端口时的访问受限问题
type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`     // 允许的源列表 (如 ["http://localhost:3000", "*"])
	AllowMethods     []string `mapstructure:"allow_methods"`     // 允许的HTTP方法 (如 ["GET", "POST", "PUT", "DELETE"])
	AllowHeaders     []string `mapstructure:"allow_headers"`     // 允许的请求头 (如 ["Content-Type", "Authorization"])
	ExposeHeaders    []string `mapstructure:"expose_headers"`    // 暴露给前端的响应头 (如 ["Content-Length", "X-Token"])
	AllowCredentials bool     `mapstructure:"allow_credentials"` // 是否允许携带凭证(Cookie/Auth Header)
	MaxAge           string   `mapstructure:"max_age"`           // 预检请求(OPTIONS)的缓存时间 (如 "12h")
}

// RAGConfig 检索增强生成(RAG)配置
type RAGConfig struct {
	ChunkSize   int `mapstructure:"chunk_size"`   // 文本切片大小 (如 500-1000字符)。过大可能丢失细节，过小可能割裂语义。
	OverlapSize int `mapstructure:"overlap_size"` // 切片重叠部分大小 (如 100-200字符)。用于保持上下文连贯性，防止语义被切断。
}

// LLMConfig 大语言模型(LLM)全局默认配置
// 这些配置用于初始化默认的 LLM 客户端，也可以在具体的 Agent 中被覆盖
type LLMConfig struct {
	APIKey      string  `mapstructure:"api_key"`     // 模型服务的 API Key (OpenAI/DeepSeek等)
	Model       string  `mapstructure:"model"`       // 模型名称 (如 gpt-4o, deepseek-chat)
	BaseURL     string  `mapstructure:"base_url"`    // 模型服务地址 (如 https://api.openai.com/v1 或本地 Ollama 地址)
	MaxTokens   int     `mapstructure:"max_tokens"`  // 单次回复最大 Token 数 (防止模型输出过长)
	Temperature float32 `mapstructure:"temperature"` // 温度值 (0.0-1.0)。值越高回复越随机/有创意，值越低回复越严谨/确定。
}
