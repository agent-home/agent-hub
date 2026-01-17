# Security Policy

## Supported Versions

We release patches for security vulnerabilities in the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of AgentHub seriously. If you have discovered a security vulnerability, we appreciate your help in disclosing it to us responsibly.

### How to Report

**Please DO NOT report security vulnerabilities through public GitHub issues.**

Instead, please report them via one of the following methods:

1. **Email**: Send an email to [security@agenthub.dev](mailto:security@agenthub.dev)
2. **GitHub Security Advisories**: Use [GitHub's private vulnerability reporting](https://github.com/agenthub/agenthub/security/advisories/new)

### What to Include

Please include the following information in your report:

- **Description**: A clear description of the vulnerability
- **Impact**: The potential impact of the vulnerability
- **Steps to Reproduce**: Detailed steps to reproduce the issue
- **Affected Versions**: Which versions are affected
- **Suggested Fix**: If you have one, a suggested fix or mitigation

### Response Timeline

- **Initial Response**: Within 48 hours
- **Status Update**: Within 7 days
- **Resolution Target**: Within 90 days (depending on severity)

### What to Expect

1. **Acknowledgment**: We'll acknowledge your report within 48 hours
2. **Assessment**: Our security team will assess the vulnerability
3. **Updates**: We'll keep you informed of our progress
4. **Fix**: We'll work on a fix and coordinate disclosure
5. **Credit**: We'll credit you in our security advisory (unless you prefer to remain anonymous)

## Security Best Practices

### For Users

1. **Keep Updated**: Always use the latest version of AgentHub
2. **API Keys**: Never commit API keys or secrets to repositories
3. **Access Control**: Use appropriate access controls for your agents
4. **Audit Logs**: Regularly review audit logs for suspicious activity

### For Agent Developers

1. **Input Validation**: Always validate and sanitize inputs
2. **Least Privilege**: Request only necessary permissions
3. **Secure Dependencies**: Keep dependencies updated
4. **No Hardcoded Secrets**: Never hardcode secrets in agent code

## Security Features

AgentHub includes several security features:

- **Authentication**: JWT-based authentication with refresh tokens
- **Authorization**: Role-based access control (RBAC)
- **Rate Limiting**: API rate limiting to prevent abuse
- **Audit Logging**: Comprehensive audit logs
- **Input Validation**: Server-side input validation
- **HTTPS**: All communications encrypted with TLS
- **Agent Sandboxing**: Isolated execution environments for agents

## Vulnerability Disclosure Policy

- We follow responsible disclosure practices
- We will not take legal action against researchers who follow this policy
- We will work with researchers to understand and resolve issues quickly
- We will credit researchers who report valid vulnerabilities

## Bug Bounty Program

We are working on establishing a bug bounty program. Details will be announced soon.

---

Thank you for helping keep AgentHub and our users safe!
