# Contributing to AgentHub

First off, thank you for considering contributing to AgentHub! It's people like you that make AgentHub such a great tool for the AI agent community.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Pull Request Process](#pull-request-process)
- [Style Guides](#style-guides)
- [Community](#community)

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

### Issues

- **Bug Reports**: If you find a bug, please create an issue using the bug report template
- **Feature Requests**: Have an idea? Open an issue using the feature request template
- **Questions**: Use GitHub Discussions for questions

### Good First Issues

Looking for a place to start? Check out issues labeled [`good first issue`](https://github.com/agenthub/agenthub/labels/good%20first%20issue).

## How Can I Contribute?

### 1. Code Contributions

- Fix bugs
- Implement new features
- Improve documentation
- Write tests
- Optimize performance

### 2. Non-Code Contributions

- Report bugs
- Suggest features
- Improve documentation
- Help other users
- Spread the word

### 3. Agent Contributions

- Create and share agents on AgentHub
- Write tutorials for building agents
- Review community agents

## Development Setup

### Prerequisites

- Go 1.21+
- Node.js 18+
- Docker & Docker Compose
- PostgreSQL 15+

### Clone and Setup

```bash
# Clone the repository
git clone https://github.com/agenthub/agenthub.git
cd agenthub

# Backend setup
cd server
cp .env.example .env
docker-compose up -d postgres redis
go mod download
go run cmd/server/main.go

# Frontend setup (new terminal)
cd web
npm install
npm run dev

# CLI setup (new terminal)
cd cli
go install .
```

### Running Tests

```bash
# Backend tests
cd server
go test ./...

# Frontend tests
cd web
npm test

# E2E tests
npm run test:e2e
```

## Pull Request Process

### 1. Before You Start

- Check existing issues and PRs to avoid duplicate work
- For major changes, open an issue first to discuss

### 2. Branch Naming

```
feature/short-description
fix/issue-number-description
docs/what-you-documented
refactor/what-you-refactored
```

### 3. Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add agent version comparison
fix: resolve API rate limiting issue
docs: update CLI installation guide
refactor: simplify storage interface
test: add handler unit tests
chore: update dependencies
```

### 4. PR Checklist

- [ ] Code follows the style guidelines
- [ ] Self-reviewed the code
- [ ] Added/updated tests
- [ ] Updated documentation
- [ ] No new warnings
- [ ] Related issues linked

### 5. Review Process

1. Submit PR
2. Automated checks run (CI, linting, tests)
3. Maintainer reviews
4. Address feedback
5. Approval and merge

## Style Guides

### Go Style Guide

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt` and `golangci-lint`
- Write clear, idiomatic Go code

```go
// Good
func (s *Storage) GetAgent(ctx context.Context, id string) (*Agent, error) {
    if id == "" {
        return nil, ErrInvalidID
    }
    // ...
}

// Avoid
func (s *Storage) GetAgent(id string) *Agent {
    // Missing context, error handling
}
```

### TypeScript/React Style Guide

- Use TypeScript strict mode
- Prefer functional components with hooks
- Use ESLint and Prettier

```tsx
// Good
interface AgentCardProps {
  agent: Agent;
  onSelect?: (agent: Agent) => void;
}

export function AgentCard({ agent, onSelect }: AgentCardProps) {
  return (
    <div onClick={() => onSelect?.(agent)}>
      {agent.name}
    </div>
  );
}
```

### Documentation Style

- Use clear, concise language
- Include code examples
- Keep README up to date

## Project Structure

```
agent-hub/
├── server/          # Go backend
│   ├── cmd/         # Entry points
│   ├── internal/    # Private packages
│   └── pkg/         # Public packages
├── cli/             # Go CLI tool
├── web/             # Next.js frontend
├── spec/            # AgentSpec schema
└── docs/            # Documentation
```

## Community

- **Discord**: [Join our server](https://discord.gg/agenthub)
- **Twitter**: [@AgentHubAI](https://twitter.com/agenthubai)
- **Discussions**: [GitHub Discussions](https://github.com/agenthub/agenthub/discussions)

## Recognition

Contributors are recognized in:
- README.md Contributors section
- Release notes
- Annual contributor spotlight

---

Thank you for contributing to AgentHub! 🚀
