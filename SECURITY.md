# Security Notice

## Known Vulnerabilities

This repository **intentionally** contains a vulnerable dependency for demonstration and testing purposes.

### CVE-2020-28483

**Package:** `github.com/gin-gonic/gin`
**Vulnerable Version:** v1.6.3 (and earlier)
**Severity:** Medium
**Type:** Directory Traversal

#### Description

The vulnerable version of Gin allows directory traversal attacks through the static file serving functionality. An attacker could potentially access files outside the intended static file directory by using specially crafted URLs with path traversal sequences (e.g., `../`).

#### Affected Code

The vulnerability is in the static file handler:
```go
r.Static("/static", "./static")
```

#### Detection

You can detect this vulnerability using:

```bash
# Using govulncheck
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# Using go list
go list -m github.com/gin-gonic/gin
# Should show: github.com/gin-gonic/gin v1.6.3

# Using dependency scanning tools
# - GitHub Dependabot
# - Snyk
# - OWASP Dependency-Check
```

#### Remediation

To fix this vulnerability, update to Gin v1.7.0 or later:

```bash
go get github.com/gin-gonic/gin@v1.7.0
go mod tidy
```

Or in go.mod:
```
require github.com/gin-gonic/gin v1.7.0
```

## Why This Exists

This vulnerable dependency is included intentionally to:
- Demonstrate security scanning capabilities
- Test dependency vulnerability detection tools
- Provide a safe environment for security training
- Validate CI/CD security gates

## Do Not Use In Production

⚠️ **WARNING:** This code should never be deployed to production environments. It is for educational and testing purposes only.

## Security Best Practices

When developing real applications:

1. Regularly update dependencies
2. Use `govulncheck` in CI/CD pipelines
3. Enable GitHub Dependabot or similar tools
4. Implement security scanning in your workflow
5. Follow the principle of least privilege
6. Validate all user inputs
7. Use security headers
8. Implement rate limiting
9. Add authentication and authorization
10. Use HTTPS in production

## Reporting Issues

This is an example repository. If you find additional security issues beyond the intentional CVE-2020-28483, please report them for educational purposes.
