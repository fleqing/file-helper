# file-helper

基于阿里云 OSS 的文件上传/下载命令行工具。

## 安装

从 [Releases](https://github.com/fleqing/file-helper/releases) 页面下载对应平台的压缩包，解压后将 `file-helper` 可执行文件放到 PATH 目录中。

或从源码构建：

```bash
go install file-helper@latest
```

## 配置

工具通过以下 4 个环境变量读取阿里云 OSS 配置，使用前必须全部设置：

| 环境变量 | 说明 | 示例 |
|---|---|---|
| `OSS_ACCESS_KEY_ID` | AccessKey ID | `LTAI5t...` |
| `OSS_ACCESS_KEY_SECRET` | AccessKey Secret | `abc123...` |
| `OSS_BUCKET` | Bucket 名称 | `my-bucket` |
| `OSS_ENDPOINT` | Endpoint 地址 | `oss-cn-hangzhou.aliyuncs.com` |

```bash
export OSS_ACCESS_KEY_ID=your-access-key-id
export OSS_ACCESS_KEY_SECRET=your-access-key-secret
export OSS_BUCKET=your-bucket-name
export OSS_ENDPOINT=oss-cn-hangzhou.aliyuncs.com
```

## 使用方法

### 上传文件

将本地文件上传到 OSS，OSS 对象名自动取文件名（不含目录路径）。

```bash
file-helper upload <本地文件路径>
```

示例：

```bash
file-helper upload ./report.pdf
# 输出: uploaded: report.pdf
# OSS 对象名: report.pdf

file-helper upload /tmp/images/photo.jpg
# 输出: uploaded: photo.jpg
# OSS 对象名: photo.jpg
```

### 下载文件

从 OSS 下载指定对象名的文件，保存到当前目录。

```bash
file-helper download <OSS 对象名>
```

示例：

```bash
file-helper download report.pdf
# 输出: downloaded: report.pdf
# 保存到: ./report.pdf
```

## 发布新版本

### 版本号规则（语义化版本 SemVer）

| 情况 | 操作 | 示例 |
|---|---|---|
| 不兼容的变更（破坏性更新） | 升级主版本号 | `v1.x.x → v2.0.0` |
| 新增功能（向后兼容） | 升级次版本号 | `v1.0.x → v1.1.0` |
| Bug 修复 | 升级修订号 | `v1.0.0 → v1.0.1` |

### 发布流程

提交代码后，打 tag 并推送即可自动触发 GitHub Actions 构建并发布 Release：

```bash
git tag v1.1.0
git push origin v1.1.0
```

GitHub Actions 自动完成：
- 交叉编译 5 个平台的二进制（Linux/macOS/Windows）
- 打包为 `.tar.gz`（Linux/macOS）和 `.zip`（Windows）
- 生成 `checksums.txt`（SHA256）
- 创建 GitHub Release 并上传所有产物
