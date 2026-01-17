<p align="center">
  <img src="docs/assets/logo.svg" alt="AgentHub Logo" width="120" height="120">
</p>

<h1 align="center">AgentHub</h1>

<p align="center">
  <strong>The Open Platform for AI Agents</strong><br>
  Discover, Share, and Run AI Agents — The Hugging Face for Agents
</p>

<p align="center">
  <a href="https://github.com/agenthub/agenthub/actions"><img src="https://img.shields.io/github/actions/workflow/status/agenthub/agenthub/ci.yml?branch=main&style=flat-square" alt="Build Status"></a>
  <a href="https://github.com/agenthub/agenthub/releases"><img src="https://img.shields.io/github/v/release/agenthub/agenthub?style=flat-square" alt="Release"></a>
  <a href="https://goreportcard.com/report/github.com/agenthub/agenthub"><img src="https://goreportcard.com/badge/github.com/agenthub/agenthub?style=flat-square" alt="Go Report Card"></a>
  <a href="https://github.com/agenthub/agenthub/blob/main/LICENSE"><img src="https://img.shields.io/github/license/agenthub/agenthub?style=flat-square" alt="License"></a>
  <a href="https://discord.gg/agenthub"><img src="https://img.shields.io/discord/123456789?style=flat-square&logo=discord" alt="Discord"></a>
</p>

<p align="center">
  <a href="https://agenthub.dev">Website</a> •
  <a href="https://docs.agenthub.dev">Documentation</a> •
  <a href="https://agenthub.dev/agents">Agent Hub</a> •
  <a href="https://discord.gg/agenthub">Discord</a> •
  <a href="https://twitter.com/agenthub">Twitter</a>
</p>

---

## What is AgentHub?

**AgentHub** is an open-source platform for discovering, sharing, and running AI agents. Think of it as the **Hugging Face for AI Agents** — a centralized ecosystem where developers can publish their agents and users can find and use them with a single command.

```bash
# Install an agent
agenthub pull agenthub/code-reviewer

# Run it
agenthub run agenthub/code-reviewer
```

### Why AgentHub?

The AI agent ecosystem is fragmented. Each platform has its own format, SDK, and deployment method. AgentHub solves this by providing:

| Problem | AgentHub Solution |
|---------|-------------------|
| Fragmented ecosystems | Unified **AgentSpec** standard |
| Hard to discover agents | Centralized **Agent Registry** with search |
| Complex deployment | **One-command** pull and run |
| No version control | **Semantic versioning** for agents |
| Vendor lock-in | **Open-source** and self-hostable |

## Features

- **Open Registry** — Publish and discover agents from the community
- **AgentSpec Standard** — Unified YAML specification for defining agents
- **Version Control** — Semantic versioning with full history
- **Multi-Runtime** — Support for Prompt, Python, Node.js, Docker, and Remote endpoints
- **Playground** — Try agents directly in the browser
- **CLI Tool** — Search, pull, run, and publish from the terminal
- **Self-Hostable** — Deploy your own private registry

## Quick Start

### Installation

```bash
# Using Go
go install github.com/agenthub/cli@latest

# Using Homebrew (macOS/Linux)
brew install agenthub/tap/agenthub

# Using npm
npm install -g @agenthub/cli

# Verify installation
agenthub version
```

### Basic Usage

```bash
# Search for agents
agenthub search "code review"

# Pull an agent
agenthub pull agenthub/code-reviewer

# Run an agent interactively
agenthub run agenthub/code-reviewer

# Run with a single input
agenthub run agenthub/code-reviewer -i "Review this Python code: def add(a,b): return a+b"
```

### Publishing Your Agent

```bash
# Create a new agent project
agenthub init my-awesome-agent
cd my-awesome-agent

# Edit the agentspec.yaml (see specification below)
vim agentspec.yaml

# Login to AgentHub
agenthub login

# Publish your agent
agenthub push
```

## AgentSpec Specification

Every agent is defined by an `agentspec.yaml` file. This standardized format ensures compatibility and portability across different platforms.

```yaml
version: "1.0.0"

metadata:
  name: code-reviewer
  description: Professional code review agent with security analysis
  author: agenthub
  license: Apache-2.0
  tags: [coding, review, security]
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
      description: Analyze code structure and complexity
    - name: check_security
      description: Check for security vulnerabilities

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

For the complete specification, see [AgentSpec Reference](docs/agentspec.md).

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                           AgentHub Platform                          │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   ┌─────────────┐     ┌─────────────┐     ┌─────────────┐          │
│   │   Web UI    │     │  REST API   │     │     CLI     │          │
│   │  (Next.js)  │     │    (Go)     │     │    (Go)     │          │
│   └──────┬──────┘     └──────┬──────┘     └──────┬──────┘          │
│          │                   │                   │                  │
│          └───────────────────┼───────────────────┘                  │
│                              │                                       │
│   ┌──────────────────────────┴──────────────────────────┐           │
│   │                    Core Services                     │           │
│   │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────┐ │           │
│   │  │ Registry │ │ Runtime  │ │ Gateway  │ │  Auth  │ │           │
│   │  │ Service  │ │ Engine   │ │ Service  │ │Service │ │           │
│   │  └──────────┘ └──────────┘ └──────────┘ └────────┘ │           │
│   └──────────────────────────┬──────────────────────────┘           │
│                              │                                       │
│   ┌──────────────────────────┴──────────────────────────┐           │
│   │                   Infrastructure                     │           │
│   │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────┐ │           │
│   │  │PostgreSQL│ │  Redis   │ │ S3/COS   │ │  K8s   │ │           │
│   │  └──────────┘ └──────────┘ └──────────┘ └────────┘ │           │
│   └─────────────────────────────────────────────────────┘           │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

## Project Structure

```
agent-hub/
├── spec/                        # AgentSpec specification
│   ├── agentspec.schema.json    # JSON Schema for validation
│   └── examples/                # Example agent definitions
│
├── server/                      # Backend services (Go)
│   ├── cmd/server/              # Application entry point
│   ├── internal/
│   │   ├── api/                 # REST API handlers
│   │   ├── config/              # Configuration management
│   │   ├── models/              # Domain models
│   │   └── storage/             # Data access layer
│   ├── migrations/              # Database migrations
│   ├── Dockerfile
│   └── docker-compose.yml
│
├── cli/                         # Command-line interface (Go)
│   ├── cmd/                     # CLI commands
│   ├── go.mod
│   └── main.go
│
├── web/                         # Web frontend (Next.js)
│   ├── src/
│   │   ├── app/                 # App router pages
│   │   └── components/          # React components
│   ├── package.json
│   └── tailwind.config.ts
│
├── docs/                        # Documentation
├── scripts/                     # Build and deployment scripts
└── README.md
```

## Development

### Prerequisites

- Go 1.22+
- Node.js 20+
- Docker & Docker Compose
- PostgreSQL 16+
- Redis 7+

### Local Setup

```bash
# Clone the repository
git clone https://github.com/agenthub/agenthub.git
cd agenthub

# Start infrastructure
docker-compose up -d postgres redis

# Run database migrations
cd server && go run cmd/migrate/main.go up

# Start the API server
go run cmd/server/main.go

# In another terminal, start the web UI
cd web && npm install && npm run dev

# Build the CLI
cd cli && go build -o agenthub .
```

### Running Tests

```bash
# Backend tests
cd server && go test ./...

# Frontend tests
cd web && npm test

# E2E tests
npm run test:e2e
```

### Code Style

We use standard tooling to maintain code quality:

```bash
# Go
gofmt -w .
golangci-lint run

# TypeScript
npm run lint
npm run format
```

## API Reference

### Authentication

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/v1/auth/register` | POST | Register a new user |
| `/api/v1/auth/login` | POST | Login and get token |
| `/api/v1/auth/refresh` | POST | Refresh access token |

### Agents

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/v1/agents` | GET | List agents |
| `/api/v1/agents` | POST | Create agent |
| `/api/v1/agents/:ns/:name` | GET | Get agent details |
| `/api/v1/agents/:ns/:name` | PUT | Update agent |
| `/api/v1/agents/:ns/:name` | DELETE | Delete agent |
| `/api/v1/agents/:ns/:name/versions` | GET | List versions |
| `/api/v1/agents/:ns/:name/versions` | POST | Publish version |

### Invocation

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/invoke/:ns/:name` | POST | Invoke agent (sync) |
| `/invoke/:ns/:name/stream` | POST | Invoke agent (streaming) |

For complete API documentation, see [API Reference](docs/api.md).

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_ADDRESS` | Server bind address | `:8080` |
| `SERVER_MODE` | Run mode (debug/release) | `debug` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | PostgreSQL user | `agenthub` |
| `DB_PASSWORD` | PostgreSQL password | - |
| `DB_NAME` | PostgreSQL database | `agenthub` |
| `REDIS_HOST` | Redis host | `localhost` |
| `REDIS_PORT` | Redis port | `6379` |
| `JWT_SECRET` | JWT signing secret | - |
| `STORAGE_TYPE` | Storage backend (local/s3) | `local` |

## Roadmap

### Released
- [x] AgentSpec v1.0 specification
- [x] Core REST API
- [x] CLI tool (search, pull, push, run)
- [x] Web platform MVP

### In Progress
- [ ] OAuth authentication (GitHub, Google)
- [ ] Agent runtime execution engine
- [ ] LLM provider integrations

### Planned
- [ ] Python SDK
- [ ] Node.js SDK
- [ ] Multi-agent workflows
- [ ] Elasticsearch integration
- [ ] S3/COS file storage
- [ ] Usage analytics & billing
- [ ] Enterprise features (SSO, audit logs)
- [ ] Kubernetes operator

See our [public roadmap](https://github.com/orgs/agenthub/projects/1) for more details.

## Contributing

We welcome contributions from the community! Please read our [Contributing Guide](CONTRIBUTING.md) before submitting a Pull Request.

### Ways to Contribute

- **Report bugs** — File an issue describing the bug
- **Suggest features** — Open a discussion for new ideas
- **Submit PRs** — Fix bugs or implement features
- **Improve docs** — Fix typos or clarify documentation
- **Share agents** — Publish useful agents to the registry

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Community

- **Discord** — [Join our Discord server](https://discord.gg/agenthub) for discussions
- **Twitter** — Follow [@agenthub](https://twitter.com/agenthub) for updates
- **Blog** — Read our [engineering blog](https://agenthub.dev/blog)
- **GitHub Discussions** — [Ask questions](https://github.com/agenthub/agenthub/discussions)

## Security

If you discover a security vulnerability, please send an email to security@agenthub.dev. Do not open a public issue.

See [SECURITY.md](SECURITY.md) for our security policy.

## License

This project is licensed under the **MIT License** — see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by [Hugging Face](https://huggingface.co), [npm](https://npmjs.com), and [Docker Hub](https://hub.docker.com)
- Built with [Go](https://go.dev), [Next.js](https://nextjs.org), and [Tailwind CSS](https://tailwindcss.com)
- Thanks to all [contributors](https://github.com/agenthub/agenthub/graphs/contributors)

---

<p align="center">
  <sub>Built with ❤️ by the AgentHub community</sub>
</p>
