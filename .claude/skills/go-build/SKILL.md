---
name: go-build
description: Builds Go projects with specified output path. Activates when modifying *.go files or running builds. Always runs golangci-lint before compilation.
---

## When to Use

This skill activates when:

- Modifying Go source files (\*.go)
- Need to compile/build Go projects
- Checking code quality before deployment

## Workflow

1. **Code Quality Check**: Run `golangci-lint run` before building
2. **Build**: Compile with specified output path

## Instructions

Before building:

- Always run `golangci-lint run` to check code quality
- Address any linting errors or warnings
- Build output to `tmp/main` or specified path

## Build Commands

**Main application:**

```bash
go build -o tmp/main main.go
```

**With verbose output:**

```bash
go build -v -o tmp/main main.go
```
