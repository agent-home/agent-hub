<p align="center">
  <img src="docs/assets/logo.svg" alt="AgentHub Logo" width="120" height="120">
</p>

<h1 align="center">AgentHub</h1>

<p align="center">
  <strong>AI 智能体的开源社区</strong><br>
  发现、分享和运行 AI 智能体 — 智能体领域的 Hugging Face
</p>

<p align="center">
  <a href="https://github.com/agenthub/agenthub/actions"><img src="https://img.shields.io/github/actions/workflow/status/agenthub/agenthub/ci.yml?branch=main&style=flat-square" alt="构建状态"></a>
  <a href="https://github.com/agenthub/agenthub/releases"><img src="https://img.shields.io/github/v/release/agenthub/agenthub?style=flat-square" alt="版本"></a>
  <a href="https://goreportcard.com/report/github.com/agenthub/agenthub"><img src="https://goreportcard.com/badge/github.com/agenthub/agenthub?style=flat-square" alt="Go Report Card"></a>
  <a href="https://github.com/agenthub/agenthub/blob/main/LICENSE"><img src="https://img.shields.io/github/license/agenthub/agenthub?style=flat-square" alt="许可证"></a>
  <a href="https://discord.gg/agenthub"><img src="https://img.shields.io/discord/123456789?style=flat-square&logo=discord" alt="Discord"></a>
</p>

<p align="center">
  <a href="https://agenthub.dev">官网</a> •
  <a href="https://docs.agenthub.dev">文档</a> •
  <a href="https://agenthub.dev/agents">智能体市场</a> •
  <a href="https://discord.gg/agenthub">Discord</a> •
  <a href="https://twitter.com/agenthub">Twitter</a>
</p>

<p align="center">
  <a href="README.md">English</a> | <strong>简体中文</strong>
</p>

---

## 什么是 AgentHub？

**AgentHub** 是一个开源的智能体聚合与分发平台，致力于为 AI 智能体提供统一的发现、分享和运行体验。你可以把它理解为 **智能体领域的 Hugging Face** — 一个开发者可以发布智能体、用户可以一键使用智能体的中心化生态系统。

```bash
# 下载智能体
agenthub pull agenthub/code-reviewer

# 运行智能体
agenthub run agenthub/code-reviewer
```

### 为什么需要 AgentHub？

当前 AI 智能体生态高度碎片化，每个平台都有自己的格式、SDK 和部署方式。AgentHub 通过以下方式解决这些问题：

| 痛点 | AgentHub 解决方案 |
|------|-------------------|
| 生态碎片化 | 统一的 **AgentSpec** 标准规范 |
| 智能体难以发现 | 中心化的 **智能体注册中心** + 搜索 |
| 部署使用复杂 | **一行命令** 下载并运行 |
| 缺乏版本管理 | 完整的 **语义化版本** 支持 |
| 平台锁定 | **开源** 且支持私有化部署 |

## 核心特性

- **开放注册中心** — 任何人都可以发布和发现智能体
- **AgentSpec 标准** — 统一的 YAML 规范定义智能体
- **版本控制** — 语义化版本管理，完整历史记录
- **多运行时支持** — 支持 Prompt、Python、Node.js、Docker 和远程端点
- **在线试用** — 浏览器内直接体验智能体
- **命令行工具** — 搜索、下载、运行、发布一站式 CLI
- **私有化部署** — 支持部署私有智能体注册中心

## 快速开始

### 安装

```bash
# 使用 Go 安装
go install github.com/agenthub/cli@latest

# 使用 Homebrew (macOS/Linux)
brew install agenthub/tap/agenthub

# 使用 npm
npm install -g @agenthub/cli

# 验证安装
agenthub version
```

### 基本使用

```bash
# 搜索智能体
agenthub search "代码审查"

# 下载智能体
agenthub pull agenthub/code-reviewer

# 交互式运行
agenthub run agenthub/code-reviewer

# 单次输入运行
agenthub run agenthub/code-reviewer -i "审查这段代码: def add(a,b): return a+b"
```

### 发布你的智能体

```bash
# 创建新项目
agenthub init my-awesome-agent
cd my-awesome-agent

# 编辑配置文件
vim agentspec.yaml

# 登录账号
agenthub login

# 发布智能体
agenthub push
```

## AgentSpec 规范

每个智能体都通过 `agentspec.yaml` 文件定义。这种标准化格式确保了跨平台的兼容性和可移植性。

```yaml
version: "1.0.0"

metadata:
  name: code-reviewer
  description: 专业的代码审查智能体，支持安全分析
  author: agenthub
  license: Apache-2.0
  tags: [编程, 审查, 安全]
  category: coding
  repository: https://github.com/agenthub/code-reviewer

runtime:
  type: python
  entry: agent.py
  python:
    version: "3.11"
    requirements: requirements.txt

model:
  provider: openai
  name: gpt-4o
  parameters:
    temperature: 0.3
    max_tokens: 8192

capabilities:
  streaming: true
  multimodal:
    text: true
    image: true
  tools:
    - name: analyze_code
      description: 分析代码结构和复杂度
    - name: check_security
      description: 检查安全漏洞

interface:
  input:
    type: json
    schema:
      type: object
      properties:
        code: { type: string }
        language: { type: string }
  output:
    type: json

pricing:
  model: free
```

完整规范请参考 [AgentSpec 参考文档](docs/agentspec.md)。

## 系统架构

```
┌─────────────────────────────────────────────────────────────────────┐
│                           AgentHub 平台                              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   ┌─────────────┐     ┌─────────────┐     ┌─────────────┐          │
│   │   Web 界面   │     │  REST API   │     │   CLI 工具  │          │
│   │  (Next.js)  │     │    (Go)     │     │    (Go)     │          │
│   └──────┬──────┘     └──────┬──────┘     └──────┬──────┘          │
│          │                   │                   │                  │
│          └───────────────────┼───────────────────┘                  │
│                              │                                       │
│   ┌──────────────────────────┴──────────────────────────┐           │
│   │                      核心服务                        │           │
│   │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────┐ │           │
│   │  │ 注册中心  │ │ 运行引擎  │ │ API 网关 │ │ 认证服务│ │           │
│   │  │ Registry │ │ Runtime  │ │ Gateway  │ │  Auth  │ │           │
│   │  └──────────┘ └──────────┘ └──────────┘ └────────┘ │           │
│   └──────────────────────────┬──────────────────────────┘           │
│                              │                                       │
│   ┌──────────────────────────┴──────────────────────────┐           │
│   │                      基础设施                        │           │
│   │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────┐ │           │
│   │  │PostgreSQL│ │  Redis   │ │ S3/COS   │ │  K8s   │ │           │
│   │  └──────────┘ └──────────┘ └──────────┘ └────────┘ │           │
│   └─────────────────────────────────────────────────────┘           │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

## 项目结构

```
agent-hub/
├── spec/                        # AgentSpec 规范定义
│   ├── agentspec.schema.json    # JSON Schema 验证
│   └── examples/                # 示例智能体定义
│
├── server/                      # 后端服务 (Go)
│   ├── cmd/server/              # 应用入口
│   ├── internal/
│   │   ├── api/                 # REST API 处理器
│   │   ├── config/              # 配置管理
│   │   ├── models/              # 领域模型
│   │   └── storage/             # 数据访问层
│   ├── migrations/              # 数据库迁移
│   ├── Dockerfile
│   └── docker-compose.yml
│
├── cli/                         # 命令行工具 (Go)
│   ├── cmd/                     # CLI 命令
│   ├── go.mod
│   └── main.go
│
├── web/                         # Web 前端 (Next.js)
│   ├── src/
│   │   ├── app/                 # 页面路由
│   │   └── components/          # React 组件
│   ├── package.json
│   └── tailwind.config.ts
│
├── docs/                        # 文档
├── scripts/                     # 构建和部署脚本
└── README.md
```

## 本地开发

### 环境要求

- Go 1.22+
- Node.js 20+
- Docker & Docker Compose
- PostgreSQL 16+
- Redis 7+

### 本地启动

```bash
# 克隆仓库
git clone https://github.com/agenthub/agenthub.git
cd agenthub

# 启动基础设施
docker-compose up -d postgres redis

# 运行数据库迁移
cd server && go run cmd/migrate/main.go up

# 启动 API 服务
go run cmd/server/main.go

# 新开终端，启动 Web 界面
cd web && npm install && npm run dev

# 构建 CLI 工具
cd cli && go build -o agenthub .
```

### 运行测试

```bash
# 后端测试
cd server && go test ./...

# 前端测试
cd web && npm test

# 端到端测试
npm run test:e2e
```

### 代码规范

我们使用标准工具来维护代码质量：

```bash
# Go
gofmt -w .
golangci-lint run

# TypeScript
npm run lint
npm run format
```

## API 参考

### 认证接口

| 接口 | 方法 | 描述 |
|------|------|------|
| `/api/v1/auth/register` | POST | 注册新用户 |
| `/api/v1/auth/login` | POST | 登录获取令牌 |
| `/api/v1/auth/refresh` | POST | 刷新访问令牌 |

### 智能体接口

| 接口 | 方法 | 描述 |
|------|------|------|
| `/api/v1/agents` | GET | 获取智能体列表 |
| `/api/v1/agents` | POST | 创建智能体 |
| `/api/v1/agents/:ns/:name` | GET | 获取智能体详情 |
| `/api/v1/agents/:ns/:name` | PUT | 更新智能体 |
| `/api/v1/agents/:ns/:name` | DELETE | 删除智能体 |
| `/api/v1/agents/:ns/:name/versions` | GET | 获取版本列表 |
| `/api/v1/agents/:ns/:name/versions` | POST | 发布新版本 |

### 调用接口

| 接口 | 方法 | 描述 |
|------|------|------|
| `/invoke/:ns/:name` | POST | 同步调用智能体 |
| `/invoke/:ns/:name/stream` | POST | 流式调用智能体 |

完整 API 文档请参考 [API 参考](docs/api.md)。

## 配置说明

### 环境变量

| 变量 | 描述 | 默认值 |
|------|------|--------|
| `SERVER_ADDRESS` | 服务监听地址 | `:8080` |
| `SERVER_MODE` | 运行模式 (debug/release) | `debug` |
| `DB_HOST` | PostgreSQL 主机 | `localhost` |
| `DB_PORT` | PostgreSQL 端口 | `5432` |
| `DB_USER` | PostgreSQL 用户名 | `agenthub` |
| `DB_PASSWORD` | PostgreSQL 密码 | - |
| `DB_NAME` | PostgreSQL 数据库名 | `agenthub` |
| `REDIS_HOST` | Redis 主机 | `localhost` |
| `REDIS_PORT` | Redis 端口 | `6379` |
| `JWT_SECRET` | JWT 签名密钥 | - |
| `STORAGE_TYPE` | 存储后端 (local/s3) | `local` |

## 开发路线图

### 已完成
- [x] AgentSpec v1.0 规范
- [x] 核心 REST API
- [x] CLI 工具 (search, pull, push, run)
- [x] Web 平台 MVP

### 进行中
- [ ] OAuth 认证 (GitHub, Google, 微信)
- [ ] 智能体运行时执行引擎
- [ ] LLM 提供商集成

### 计划中
- [ ] Python SDK
- [ ] Node.js SDK
- [ ] 多智能体工作流
- [ ] Elasticsearch 全文搜索
- [ ] S3/COS 文件存储
- [ ] 用量统计与计费
- [ ] 企业版功能 (SSO, 审计日志)
- [ ] Kubernetes Operator

查看我们的 [公开路线图](https://github.com/orgs/agenthub/projects/1) 了解更多详情。

## 参与贡献

我们欢迎社区贡献！在提交 Pull Request 之前，请先阅读我们的 [贡献指南](CONTRIBUTING.md)。

### 贡献方式

- **报告 Bug** — 提交 Issue 描述遇到的问题
- **建议功能** — 发起 Discussion 讨论新想法
- **提交 PR** — 修复 Bug 或实现新功能
- **完善文档** — 修正错误或改进说明
- **分享智能体** — 发布有用的智能体到注册中心

### 开发流程

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m '添加某个很棒的功能'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 发起 Pull Request

## 社区

- **Discord** — [加入 Discord 服务器](https://discord.gg/agenthub) 参与讨论
- **Twitter** — 关注 [@agenthub](https://twitter.com/agenthub) 获取最新动态
- **博客** — 阅读我们的 [技术博客](https://agenthub.dev/blog)
- **GitHub Discussions** — [提问和讨论](https://github.com/agenthub/agenthub/discussions)
- **微信群** — 扫码加入微信交流群

## 安全

如果你发现安全漏洞，请发送邮件至 security@agenthub.dev，请勿公开提交 Issue。

详细信息请参考 [SECURITY.md](SECURITY.md)。

## 许可证

本项目采用 **MIT 许可证** — 详见 [LICENSE](LICENSE) 文件。

## 致谢

- 灵感来源于 [Hugging Face](https://huggingface.co)、[npm](https://npmjs.com) 和 [Docker Hub](https://hub.docker.com)
- 基于 [Go](https://go.dev)、[Next.js](https://nextjs.org) 和 [Tailwind CSS](https://tailwindcss.com) 构建
- 感谢所有 [贡献者](https://github.com/agenthub/agenthub/graphs/contributors)

---

<p align="center">
  <sub>由 AgentHub 社区用 ❤️ 构建</sub>
</p>
