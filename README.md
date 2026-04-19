# SA - Simple Agent

一个包含三个子项目的 Go 项目，使用 go.work 进行多模块管理。

## 项目结构

```
sa/
├── go.work                      # Go workspace 多模块管理
├── pkg/                         # 公共可复用代码
│   ├── config/                  # 配置加载
│   ├── llm/                     # LLM 客户端封装
│   └── response/                # 统一响应格式
├── services/
│   ├── cli/                     # 1️⃣ 用户端 CLI Agent
│   │   ├── main.go              # 入口
│   │   └── tui/                 # BubbleTea TUI 界面
│   ├── gateway/                 # 2️⃣ 网关代理控制台
│   │   ├── main.go              # 入口
│   │   ├── app/                 # 路由配置
│   │   ├── handler/             # HTTP handlers
│   │   └── middleware/          # 中间件
│   └── knowledge/               # 3️⃣ 向量知识库服务
│       ├── main.go              # 入口
│       ├── app/                 # 路由配置
│       ├── handler/             # HTTP handlers
│       ├── embedding/           # DashScope Embedding 封装
│       └── qdrant/              # Qdrant HTTP 客户端
└── configs/
    └── config.yaml              # 配置示例
```

## 子项目说明

### 1. CLI Agent (`services/cli`)
用户端 TUI 聊天 Agent，使用 BubbleTea 构建交互界面。

**技术栈:**
- BubbleTea (TUI)
- Glamour (Markdown 渲染)
- OpenAI SDK (兼容 DashScope)

**启动:**
```bash
cd services/cli
go run .
```

### 2. Gateway (`services/gateway`)
API 网关服务，提供聊天接口和代理功能。

**技术栈:**
- Gin (HTTP 框架)
- OpenAI SDK

**API 端点:**
- `GET /health` - 健康检查
- `POST /api/v1/chat` - 单次聊天
- `POST /api/v1/chat/stream` - 流式聊天 (SSE)

**启动:**
```bash
cd services/gateway
go run .
```

### 3. Knowledge (`services/knowledge`)
向量知识库服务，支持文档入库和相似度检索。

**技术栈:**
- Gin (HTTP 框架)
- DashScope Embedding API
- Qdrant (向量数据库)

**API 端点:**
- `GET /health` - 健康检查
- `POST /api/v1/collections` - 创建集合
- `POST /api/v1/documents` - 添加文档
- `POST /api/v1/documents/search` - 搜索相似文档

**启动:**
```bash
cd services/knowledge
go run .
```

## 配置

编辑 `configs/config.yaml` 或使用环境变量:

```yaml
openai:
  api_key: ""        # 或使用 OPENAI_API_KEY 环境变量
  base_url: "https://dashscope.aliyuncs.com/compatible-mode/v1"
  model: "qwen-plus"

qdrant:
  host: "localhost"
  port: 6333
  api_key: ""
  use_http: false

server:
  host: "0.0.0.0"
  port: 8080
  mode: "debug"

embedding:
  model: "text-embedding-v3"
```

## 构建

```bash
# 构建 CLI
go build -o bin/cli ./services/cli

# 构建 Gateway
go build -o bin/gateway ./services/gateway

# 构建 Knowledge
go build -o bin/knowledge ./services/knowledge
```

## 依赖

- Go 1.24+
- Qdrant (可选，用于 Knowledge 服务)