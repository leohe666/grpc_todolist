# 1ms-helper

> 毫秒镜像（1ms.run）助手工具 - 一键配置Docker镜像加速

[![Go Version](https://img.shields.io/badge/Go-1.23.4+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows%20%7C%20Synology-lightgrey?style=flat-square)](https://cnb.cool/mliev/1ms.run/1ms-helper)

## 📋 项目简介

1ms-helper 是一个专为 [毫秒镜像（1ms.run）](https://1ms.run) 设计的命令行助手工具，旨在帮助开发者快速配置Docker镜像加速服务。支持一键配置多个主流Docker镜像仓库的加速地址和认证信息。

## ✨ 功能特性

- 🚀 **一键配置** - 快速配置毫秒镜像加速服务
- 🔐 **账号管理** - 安全管理毫秒镜像账号和认证信息
- 🌐 **多仓库支持** - 支持 Docker Hub、GitHub Container Registry、Google Container Registry 等多个镜像仓库
- 🖥️ **跨平台兼容** - 支持 Linux、macOS、Windows 和群晖 NAS 系统
- ⚡ **智能检测** - 自动检测系统环境并应用最适合的配置
- 🛡️ **连接测试** - 提供连接状态检查和问题诊断功能

## 🎯 支持的镜像仓库

| 仓库名称 | 原始地址 | 加速地址 | 状态 |
|---------|---------|---------|------|
| Docker Hub | `docker.io` | `docker.1ms.run` | ✅ |
| GitHub Container Registry | `ghcr.io` | `ghcr.1ms.run` | ✅ |
| Google Container Registry | `gcr.io` | `gcr.1ms.run` | ✅ |
| NVIDIA Container Registry | `nvcr.io` | `nvcr.1ms.run` | ✅ |
| Red Hat Quay | `quay.io` | `quay.1ms.run` | ✅ |
| Elastic Docker Registry | `docker.elastic.co` | `elastic.1ms.run` | ✅ |
| Microsoft Container Registry | `mcr.microsoft.com` | `mcr.1ms.run` | ✅ |
| Kubernetes Container Registry | `registry.k8s.io` | `k8s.1ms.run` | ✅ |

## 📦 安装方式

### 方式一：一键安装（推荐）

```bash
# Linux/macOS
curl -sSL https://static.1ms.run/1ms-helper/install.sh | bash

# 或者使用 wget
wget -qO- https://static.1ms.run/1ms-helper/install.sh | bash
```

### 方式二：手动安装

1. 访问 [Releases 页面](https://cnb.cool/mliev/1ms.run/1ms-helper/-/releases) 下载对应系统的二进制文件
2. 解压并移动到系统 PATH 目录：

```bash
# Linux/macOS 示例
tar -xzf 1ms-helper_Linux_x86_64.tar.gz
sudo mv 1ms-helper /usr/local/bin/
chmod +x /usr/local/bin/1ms-helper
```

### 方式三：源码编译

```bash
# 克隆仓库
git clone https://github.com/mliev/1ms-helper.git
cd 1ms-helper

# 安装依赖
go mod tidy

# 编译
go build -o 1ms-helper main.go

# 运行
./1ms-helper
```

## 🚀 使用指南

### 基本命令

```bash
# 显示帮助信息
1ms-helper --help

# 检查毫秒镜像连接状态
1ms-helper check

# 配置毫秒镜像账号
1ms-helper config:account

# 配置镜像加速
1ms-helper config:mirror

# 一键配置（推荐）
1ms-helper config

# 移除镜像配置
1ms-helper remove:mirror
```

### 详细使用步骤

#### 1. 检查连接状态
```bash
1ms-helper check
```
检查与毫秒镜像的网络连接状态和配置是否正确。

#### 2. 配置账号信息
```bash
1ms-helper config:account
```
按提示输入您在 [毫秒镜像](https://1ms.run) 注册的账号和密码。

#### 3. 配置镜像加速
```bash
1ms-helper config:mirror
```
根据您的系统环境自动配置 Docker daemon 的镜像加速设置。

#### 4. 一键配置（推荐新用户）
```bash
1ms-helper config
```
依次执行镜像配置和账号配置，适合首次使用的用户。

## 📁 项目结构

```
1ms-helper/
├── app/                    # 应用核心代码
│   ├── Command/           # 命令实现
│   │   ├── Check.go      # 连接检查命令
│   │   ├── Config.go     # 一键配置命令
│   │   ├── ConfigAccount.go    # 账号配置命令
│   │   ├── ConfigMirror.go     # 镜像配置命令
│   │   ├── ConfigMirror/       # 镜像配置实现
│   │   │   ├── Linux.go        # Linux系统配置
│   │   │   └── Synology.go     # 群晖系统配置
│   │   ├── RemoveAccount.go    # 账号移除命令
│   │   └── RemoveMirror/       # 镜像移除命令
│   ├── Dto/               # 数据传输对象
│   ├── Interfaces/        # 接口定义
│   ├── Lib/              # 公共库
│   │   ├── Ask.go        # 用户交互
│   │   └── Question/     # 问题类型
│   └── Utils/            # 工具类
├── cmd/                   # 命令行入口
├── config/               # 配置管理
├── main.go              # 程序入口
├── go.mod               # Go模块定义
├── go.sum               # 依赖校验
└── install.sh           # 安装脚本
```

## 🔧 开发相关

### 开发环境要求

- Go 1.23.4+
- Git

### 主要依赖

- [cobra](https://github.com/spf13/cobra) - 命令行框架
- [color](https://github.com/gookit/color) - 彩色输出
- [term](https://golang.org/x/term) - 终端控制

### 本地开发

```bash
# 克隆项目
git clone https://github.com/mliev/1ms-helper.git
cd 1ms-helper

# 安装依赖
go mod tidy

# 运行项目
go run main.go

# 运行测试
go test ./...

# 构建二进制文件
go build -o bin/1ms-helper main.go
```

### 项目打包

```bash
# 使用 goreleaser 进行跨平台打包
goreleaser release --snapshot --clean
```

## 🤝 贡献指南

我们欢迎社区贡献！请遵循以下步骤：

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开一个 Pull Request

## 📞 支持与反馈

- 🐛 **问题反馈**: [GitHub Issues](https://cnb.cool/mliev/1ms.run/1ms-helper/-/issues)
- 💬 **功能建议**: [GitHub Discussions](https://cnb.cool/mliev/1ms.run/1ms-helper/-/issues)
- 📖 **使用文档**: [毫秒镜像文档](https://www.mliev.com/docs/1ms.run)
- 🌐 **官方网站**: [https://1ms.run](https://1ms.run)

## 🙏 致谢

感谢所有为本项目做出贡献的开发者和 [毫秒镜像](https://1ms.run) 团队提供的优质服务。

---

<div align="center">
  <p>如果这个项目对您有帮助，请给我们一个 ⭐ Star！</p>
  <p>Made with ❤️ by <a href="https://github.com/mliev">mliev</a></p>
</div>