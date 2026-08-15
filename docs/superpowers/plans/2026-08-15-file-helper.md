# file-helper 实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 用 Go 语言实现一个基于阿里云 OSS 的文件上传/下载命令行工具，产出单一可执行文件。

**Architecture:** 两个源文件：`oss.go` 封装 OSS 操作，`main.go` 负责 CLI 参数解析和路由。使用标准库 `os.Args` 解析子命令，无第三方 CLI 框架。

**Tech Stack:** Go 1.21+，`github.com/aliyun/aliyun-oss-go-sdk/oss`

## Global Constraints

- Go 模块名：`file-helper`
- 只有两个子命令：`upload` 和 `download`
- 上传时 OSS 对象名 = `filepath.Base(localPath)`
- 下载文件保存到当前工作目录，文件名 = OSS 对象名
- 成功信息输出到 `stdout`，错误信息输出到 `stderr`，格式为 `error: <描述>`
- 成功退出码 `0`，失败退出码 `1`
- 配置仅来自环境变量：`OSS_ACCESS_KEY_ID`、`OSS_ACCESS_KEY_SECRET`、`OSS_BUCKET`、`OSS_ENDPOINT`

---

## 文件结构

| 文件 | 职责 |
|---|---|
| `go.mod` / `go.sum` | 模块定义和依赖锁定 |
| `oss.go` | OSS 客户端结构体、初始化、Upload、Download |
| `main.go` | 环境变量读取、参数解析、子命令路由、错误输出 |

---

### Task 1: 初始化 Go 模块并添加依赖

**Files:**
- Create: `go.mod`
- Create: `go.sum`（自动生成）

**Interfaces:**
- Produces: 可用的 `github.com/aliyun/aliyun-oss-go-sdk/oss` 包

- [ ] **Step 1: 初始化 Go 模块**

```bash
cd /Users/hero/Desktop/project/file-helper
go mod init file-helper
```

- [ ] **Step 2: 添加 OSS SDK 依赖**

```bash
go get github.com/aliyun/aliyun-oss-go-sdk/oss@latest
```

- [ ] **Step 3: 验证依赖已写入**

```bash
grep aliyun go.mod
```
预期输出包含：`github.com/aliyun/aliyun-oss-go-sdk`

- [ ] **Step 4: 提交**

```bash
git add go.mod go.sum
git commit -m "chore: initialize go module with aliyun-oss-go-sdk"
```

---

### Task 2: 实现 OSS 封装层（oss.go）

**Files:**
- Create: `oss.go`

**Interfaces:**
- Consumes: `github.com/aliyun/aliyun-oss-go-sdk/oss`
- Produces:
  - `type Client struct` — 持有 bucket 实例
  - `func NewClient(endpoint, accessKeyID, accessKeySecret, bucket string) (*Client, error)` — 初始化
  - `func (c *Client) Upload(localPath string) error` — 上传，对象名取 `filepath.Base(localPath)`
  - `func (c *Client) Download(ossKey string) error` — 下载到当前目录，文件名 = ossKey

- [ ] **Step 1: 编写 oss.go**

```go
package main

import (
	"fmt"
	"path/filepath"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Client struct {
	bucket *oss.Bucket
}

func NewClient(endpoint, accessKeyID, accessKeySecret, bucket string) (*Client, error) {
	c, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}
	b, err := c.Bucket(bucket)
	if err != nil {
		return nil, err
	}
	return &Client{bucket: b}, nil
}

func (c *Client) Upload(localPath string) error {
	key := filepath.Base(localPath)
	return c.bucket.PutObjectFromFile(key, localPath)
}

func (c *Client) Download(ossKey string) error {
	destPath := filepath.Base(ossKey)
	err := c.bucket.GetObjectToFile(ossKey, destPath)
	if err != nil {
		if ossErr, ok := err.(oss.ServiceError); ok && ossErr.StatusCode == 404 {
			return fmt.Errorf("object not found: %s", ossKey)
		}
		return err
	}
	return nil
}
```

- [ ] **Step 2: 验证语法**

```bash
go build ./...
```
预期：无错误（此时 main.go 不存在，会报 `no Go files`，属正常；只要无语法错误即可）

实际上需要先有一个临时的 main.go 才能编译，跳到 Task 3 后再回来验证。

- [ ] **Step 3: 提交**

```bash
git add oss.go
git commit -m "feat: add OSS client wrapper (oss.go)"
```

---

### Task 3: 实现 CLI 入口（main.go）

**Files:**
- Create: `main.go`

**Interfaces:**
- Consumes:
  - `NewClient(endpoint, accessKeyID, accessKeySecret, bucket string) (*Client, error)`
  - `(*Client).Upload(localPath string) error`
  - `(*Client).Download(ossKey string) error`
- Produces: 可执行二进制 `file-helper`

- [ ] **Step 1: 编写 main.go**

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: file-helper <upload|download> <path>")
		os.Exit(1)
	}

	cmd := os.Args[1]
	path := os.Args[2]

	client, err := newClientFromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	switch cmd {
	case "upload":
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "error: file not found: %s\n", path)
			os.Exit(1)
		}
		if err := client.Upload(path); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("uploaded: %s\n", path)

	case "download":
		if err := client.Download(path); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("downloaded: %s\n", path)

	default:
		fmt.Fprintln(os.Stderr, "usage: file-helper <upload|download> <path>")
		os.Exit(1)
	}
}

func newClientFromEnv() (*Client, error) {
	vars := []string{"OSS_ACCESS_KEY_ID", "OSS_ACCESS_KEY_SECRET", "OSS_BUCKET", "OSS_ENDPOINT"}
	vals := make(map[string]string, len(vars))
	for _, v := range vars {
		val := os.Getenv(v)
		if val == "" {
			return nil, fmt.Errorf("%s is not set", v)
		}
		vals[v] = val
	}
	return NewClient(
		vals["OSS_ENDPOINT"],
		vals["OSS_ACCESS_KEY_ID"],
		vals["OSS_ACCESS_KEY_SECRET"],
		vals["OSS_BUCKET"],
	)
}
```

- [ ] **Step 2: 构建验证**

```bash
go build -o file-helper .
```
预期：生成 `file-helper` 可执行文件，无错误。

- [ ] **Step 3: 验证无参数时显示用法**

```bash
./file-helper 2>&1 || true
```
预期输出：`usage: file-helper <upload|download> <path>`

- [ ] **Step 4: 验证未设置环境变量时报错**

```bash
./file-helper upload test.txt 2>&1 || true
```
预期输出：`error: OSS_ACCESS_KEY_ID is not set`

- [ ] **Step 5: 验证本地文件不存在时报错**

```bash
OSS_ACCESS_KEY_ID=x OSS_ACCESS_KEY_SECRET=x OSS_BUCKET=x OSS_ENDPOINT=x \
  ./file-helper upload nonexistent.txt 2>&1 || true
```
预期输出：`error: file not found: nonexistent.txt`

- [ ] **Step 6: 提交**

```bash
git add main.go
git commit -m "feat: add CLI entry point (main.go)"
```

---

### Task 4: 构建最终可执行文件并验证

**Files:**
- 不新增文件，验证整体可执行

**Interfaces:**
- Consumes: `main.go`、`oss.go`

- [ ] **Step 1: 清理并重新构建**

```bash
rm -f file-helper
go build -o file-helper .
```

- [ ] **Step 2: 验证帮助信息（无参数）**

```bash
./file-helper 2>&1; echo "exit: $?"
```
预期：
```
usage: file-helper <upload|download> <path>
exit: 1
```

- [ ] **Step 3: 验证错误的子命令**

```bash
OSS_ACCESS_KEY_ID=x OSS_ACCESS_KEY_SECRET=x OSS_BUCKET=x OSS_ENDPOINT=x \
  ./file-helper list something 2>&1; echo "exit: $?"
```
预期：
```
usage: file-helper <upload|download> <path>
exit: 1
```

- [ ] **Step 4: （可选）真实 OSS 集成测试**

如有真实 OSS 配置，执行以下命令端到端验证：

```bash
# 创建测试文件
echo "hello" > /tmp/test-upload.txt

# 上传
OSS_ACCESS_KEY_ID=<your-key> \
OSS_ACCESS_KEY_SECRET=<your-secret> \
OSS_BUCKET=<your-bucket> \
OSS_ENDPOINT=<your-endpoint> \
  ./file-helper upload /tmp/test-upload.txt
# 预期: uploaded: /tmp/test-upload.txt

# 下载（到当前目录）
OSS_ACCESS_KEY_ID=<your-key> \
OSS_ACCESS_KEY_SECRET=<your-secret> \
OSS_BUCKET=<your-bucket> \
OSS_ENDPOINT=<your-endpoint> \
  ./file-helper download test-upload.txt
# 预期: downloaded: test-upload.txt

# 验证内容一致
diff /tmp/test-upload.txt ./test-upload.txt && echo "OK"
```

- [ ] **Step 5: 提交二进制（忽略）并添加 .gitignore**

```bash
echo "file-helper" > .gitignore
git add .gitignore
git commit -m "chore: add .gitignore to exclude binary"
```
