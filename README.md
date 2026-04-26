# SA - Special Assiant

AI Agent 基于私有向量知识库的办公助手

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
- 自己实现 规划、上下文、历史会话，技能

**启动:**
```bash
cd services/cli
go run .
```

### 2. Gateway (`services/gateway`)
API 网关服务，提供open AI和RAG项目库代理功能，直接跟用户端交互，鉴权，token审计。

**技术栈:**
- Gin (HTTP 框架)
- OpenAI SDK

**API 端点:**
- `GET /health` - 健康检查
- `POST /llm` - 调用大模型
- `POST /rag` - 调用rag知识库

**启动:**
```bash
cd services/gateway
go run .
```

### 3. Knowledge (`services/knowledge`)
向量知识库服务，支持文档入库和相似度检索。

**技术栈:**
- Gin (HTTP 框架)
- DashScope Embedding 
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

编辑 `configs/config.yaml` 

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