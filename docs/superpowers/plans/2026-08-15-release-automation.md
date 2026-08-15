# file-helper Release 自动化实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为 file-helper 添加 goreleaser + GitHub Actions 的多平台 Release 自动化，推送 `v*` tag 时自动构建 5 个平台的二进制并发布到 GitHub Release。

**Architecture:** 新增 `.goreleaser.yaml` 定义构建目标和产物格式，新增 `.github/workflows/release.yml` 定义触发条件和 CI 步骤。不修改任何现有源码。

**Tech Stack:** goreleaser v2，GitHub Actions（`actions/checkout@v4`、`actions/setup-go@v5`、`goreleaser/goreleaser-action@v6`）

## Global Constraints

- 不修改任何现有文件（`main.go`、`oss.go`、`go.mod`、`.gitignore` 等）
- 构建 5 个平台：`linux/amd64`、`linux/arm64`、`darwin/amd64`、`darwin/arm64`、`windows/amd64`
- Linux/macOS 产物格式：`.tar.gz`；Windows 产物格式：`.zip`
- checksum 文件名：`checksums.txt`，算法：`sha256`
- GitHub Actions 仅在 `v*` tag push 时触发
- 使用内置 `GITHUB_TOKEN`，无需手动配置 Secrets
- workflow 步骤顺序：checkout（`fetch-depth: 0`）→ `setup-go@v5`（Go `1.21`）→ `goreleaser-action@v6`
- release 直接发布（`draft: false`）

---

## 文件结构

| 文件 | 职责 |
|---|---|
| `.goreleaser.yaml` | 定义构建目标平台、产物打包格式、checksum 配置、release 配置 |
| `.github/workflows/release.yml` | 定义 CI trigger（`v*` tag）、安装环境、运行 goreleaser |

---

### Task 1: 创建 .goreleaser.yaml

**Files:**
- Create: `.goreleaser.yaml`

**Interfaces:**
- Produces: goreleaser 可读取的构建配置，定义 5 个平台目标

- [ ] **Step 1: 创建 .goreleaser.yaml**

```yaml
version: 2

builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ignore:
      - goos: windows
        goarch: arm64
    binary: file-helper

archives:
  - formats:
      - tar.gz
    format_overrides:
      - goos: windows
        formats:
          - zip
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"

checksum:
  name_template: "checksums.txt"
  algorithm: sha256

release:
  draft: false
```

- [ ] **Step 2: 验证 goreleaser 配置语法**

本地安装 goreleaser 并验证（如未安装则跳过，CI 会验证）：

```bash
# macOS 安装
brew install goreleaser/tap/goreleaser

# 验证配置（--snapshot 模式不需要 tag，不会实际发布）
goreleaser check
```

预期输出：`config is valid` 或无报错。若未安装 goreleaser，直接跳到 Step 3。

- [ ] **Step 3: 提交**

```bash
git add .goreleaser.yaml
git commit -m "chore: add goreleaser config for multi-platform release"
```

---

### Task 2: 创建 GitHub Actions release workflow

**Files:**
- Create: `.github/workflows/release.yml`

**Interfaces:**
- Consumes: `.goreleaser.yaml`（由 goreleaser-action 自动读取）
- Produces: 推送 `v*` tag 时自动触发的 CI pipeline

- [ ] **Step 1: 创建 .github/workflows 目录并写入 release.yml**

```bash
mkdir -p .github/workflows
```

```yaml
name: Release

on:
  push:
    tags:
      - "v*"

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: "1.21"

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v6
        with:
          distribution: goreleaser
          version: "~> v2"
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

- [ ] **Step 2: 验证 YAML 格式**

```bash
cat .github/workflows/release.yml
```

确认缩进正确，无 tab 字符（YAML 要求空格缩进）。

- [ ] **Step 3: 提交**

```bash
git add .github/workflows/release.yml
git commit -m "ci: add GitHub Actions release workflow"
```

---

### Task 3: 推送配置并触发首次 Release

**Files:**
- 不新增文件

**Interfaces:**
- Consumes: `.goreleaser.yaml`、`.github/workflows/release.yml`

- [ ] **Step 1: 推送 commits 到远程**

```bash
git push origin main
```

- [ ] **Step 2: 打首个 tag 并推送**

```bash
git tag v1.0.0
git push origin v1.0.0
```

- [ ] **Step 3: 在 GitHub Actions 页面确认 workflow 已触发**

访问 `https://github.com/fleqing/file-helper/actions`，确认名为 `Release` 的 workflow 正在运行。

- [ ] **Step 4: 确认 GitHub Release 页面产物**

workflow 完成后，访问 `https://github.com/fleqing/file-helper/releases`，确认：
- Release `v1.0.0` 已创建
- 包含以下 5 个产物：
  - `file-helper_1.0.0_linux_amd64.tar.gz`
  - `file-helper_1.0.0_linux_arm64.tar.gz`
  - `file-helper_1.0.0_darwin_amd64.tar.gz`
  - `file-helper_1.0.0_darwin_arm64.tar.gz`
  - `file-helper_1.0.0_windows_amd64.zip`
- 包含 `checksums.txt`
