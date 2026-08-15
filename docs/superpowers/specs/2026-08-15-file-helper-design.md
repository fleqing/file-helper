# file-helper 设计规范

日期：2026-08-15

## 概述

`file-helper` 是一个基于阿里云 OSS 的文件上传/下载命令行工具，使用 Go 语言编写，产出单一可执行文件。

## 架构

单 Go 模块，三个源文件：

```
file-helper/
├── main.go    # 入口：解析子命令，路由到 upload/download
├── oss.go     # OSS 客户端封装：初始化、上传、下载
├── go.mod
└── go.sum
```

- `main.go`：用标准库 `os.Args` 解析子命令，负责参数验证和错误输出，不包含业务逻辑
- `oss.go`：封装阿里云 `aliyun-oss-go-sdk`，提供 `Upload(localPath string)` 和 `Download(ossKey string)` 两个函数

## 命令行接口

```bash
# 上传：将本地文件上传到 OSS，对象名 = 文件名（不含目录路径）
file-helper upload <本地文件路径>

# 下载：从 OSS 下载指定对象名的文件到当前目录，保存为同名文件
file-helper download <OSS对象名>

# 示例
file-helper upload ./report.pdf     # OSS 对象名: report.pdf
file-helper upload /tmp/photo.jpg   # OSS 对象名: photo.jpg
file-helper download report.pdf     # 下载到 ./report.pdf
```

**关键约定：**
- 上传时 OSS 对象名始终取文件名部分（`filepath.Base`），忽略本地目录结构
- 下载时直接用对象名作为本地保存文件名，保存到当前工作目录

## 环境变量配置

启动时从以下四个环境变量读取 OSS 配置，缺少任意一个则报错退出：

| 环境变量 | 说明 | 示例值 |
|---|---|---|
| `OSS_ACCESS_KEY_ID` | 阿里云 AccessKey ID | `LTAI5t...` |
| `OSS_ACCESS_KEY_SECRET` | 阿里云 AccessKey Secret | `abc123...` |
| `OSS_BUCKET` | Bucket 名称 | `my-bucket` |
| `OSS_ENDPOINT` | Endpoint 地址 | `oss-cn-hangzhou.aliyuncs.com` |

## 错误处理与输出

- 成功信息输出到 `stdout`，错误输出到 `stderr`
- 退出码：成功 `0`，失败 `1`

成功输出：
```
uploaded: report.pdf
downloaded: report.pdf
```

错误场景：
| 场景 | 错误信息 |
|---|---|
| 缺少环境变量 | `error: OSS_ACCESS_KEY_ID is not set` |
| 参数不足 / 子命令错误 | `usage: file-helper <upload\|download> <path>` |
| 本地文件不存在 | `error: file not found: ./report.pdf` |
| OSS 对象不存在 | `error: object not found: report.pdf` |
| 其他 OSS 错误 | `error: <OSS SDK 返回的错误信息>` |

## 依赖

- `github.com/aliyun/aliyun-oss-go-sdk`：阿里云官方 Go SDK
- Go 标准库：`os`、`path/filepath`、`fmt`
