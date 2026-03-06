# Agent `buildGraph` 流程图（Mermaid）

```mermaid
flowchart TD
    A[输入 UserMessage] --> B[InputToQuery]
    A --> C[InputToHistory]

    B --> D[Retriever]
    D --> D1[遍历 KnowledgeIDs]
    D1 --> D2[每个KB构建 MilvusRetriever]
    D2 --> D3["Retrieve(query) 得到 docs"]
    D3 --> D4[聚合 allDocuments append]
    D4 --> D5[按 MetaData.score 降序排序]
    D5 --> D6[截断到 TopK]
    D6 --> E[输出 documents]

    C --> F[ChatTemplate]
    E --> F
    F --> G[模板渲染结果]

    H{tools 是否为空} -->|否| I[React Agent]
    H -->|是| J[ChatModel]

    G --> I
    G --> J

    I --> K[输出 Message]
    J --> K
    K --> L[END]

    M[buildGraph前置装配] --> M1[GetModel]
    M1 --> M2[GetLLMClient]
    M2 --> H
    N[MCP Servers] --> N1[NewSSEMCPClient Start Initialize]
    N1 --> N2[GetTools 聚合]
    N2 --> H
```

## 节点数据传递与聚合说明

1. `START` 输入同一份 `UserMessage`，并行分发到 `InputToQuery` 与 `InputToHistory`。
2. `InputToQuery` 只提取 `query`，传给 `Retriever` 做向量检索。
3. `Retriever` 在 `MultiKBRetriever` 内部执行“跨知识库聚合”：
   - 按 `KnowledgeIDs` 循环检索；
   - 每个知识库各自召回文档；
   - 统一 `append` 到 `allDocuments`；
   - 按 `score` 排序后做全局 `TopK` 截断。
4. `InputToHistory` 产出 `map`，其中 `history/query/date` 作为 `ChatTemplate` 的变量输入。
5. `Retriever` 输出通过 `WithOutputKey("documents")` 注入模板变量 `documents`，与 `history/query` 在 `ChatTemplate` 聚合。
6. `ChatTemplate` 生成最终消息数组后：
   - 若存在 MCP tools：进入 `React Agent`（可多步工具调用，`MaxStep=10`）；
   - 否则：直接进入 `ChatModel` 单步生成。
7. 两条分支最终统一输出 `*schema.Message` 到 `END`。

## 图与代码对应关系

- 图中主流程对应 `buildGraph`：`internal/service/agent.go`。
- 聚合检索细节对应 `MultiKBRetriever.Retrieve`：`internal/component/retriever/milvus/multi_retriever.go`。
