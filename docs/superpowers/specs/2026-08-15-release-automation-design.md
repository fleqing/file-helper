# file-helper Release 自动化设计规范

日期：2026-08-15

## 概述

为 file-helper 添加基于 goreleaser + GitHub Actions 的多平台 Release 自动化流程。推送语义化版本 tag（如 `v1.0.0`）时自动触发，构建 5 个目标平台的二进制并发布到 GitHub Release。

## 架构

新增两个配置文件，不修改任何现有源码：

```
file-helper/
├── .goreleaser.yaml              # goreleaser 构建配置
└── .github/
    └── workflows/
        └── release.yml           # GitHub Actions workflow
```

## 构建目标平台

| 平台 | GOOS | GOARCH | 产物格式 |
|---|---|---|---|
| Linux x86-64 | linux | amd64 | `.tar.gz` |
| Linux ARM64 | linux | arm64 | `.tar.gz` |
| macOS x86-64 | darwin | amd64 | `.tar.gz` |
| macOS ARM64 | darwin | arm64 | `.tar.gz` |
| Windows x86-64 | windows | amd64 | `.zip` |

每次 Release 自动生成 `checksums.txt`（所有产物的 SHA256 哈希值）。

## 版本号规范

语义化版本（SemVer）：`vMAJOR.MINOR.PATCH`

| 情况 | 示例 |
|---|---|
| 不兼容的 API 变更 | `v1.0.0 → v2.0.0` |
| 新增功能（向后兼容） | `v1.0.0 → v1.1.0` |
| Bug 修复 | `v1.0.0 → v1.0.1` |

## 发布流程

```bash
git tag v1.0.0
git push origin v1.0.0
```

GitHub Actions 仅在 tag 匹配 `v*` 时触发，普通 push 不受影响。

## GitHub Actions Workflow

触发条件：`push` 到 `v*` tag。

步骤：
1. `actions/checkout@v4`（`fetch-depth: 0` 获取完整 tag 历史，goreleaser 需要）
2. `actions/setup-go@v5`（Go 1.21）
3. `goreleaser/goreleaser-action@v6`（使用内置 `GITHUB_TOKEN`，无需手动配置 Secrets）

## .goreleaser.yaml 关键配置

- `builds`：指定 5 个目标平台，二进制名称为 `file-helper`（Windows 下自动加 `.exe`）
- `archives`：Linux/macOS 打包为 `.tar.gz`，Windows 打包为 `.zip`
- `checksum`：生成 `checksums.txt`，算法 SHA256
- `release`：使用 GitHub Release，`draft: false`（直接发布，不创建草稿）
