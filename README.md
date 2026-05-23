# DevManager (STORMMOON)

开发环境版本管理器 — 轻量级跨平台命令行工具，支持 JDK、Maven 等多版本开发环境的添加、切换与卸载。

## 功能特性

- **多版本管理**：支持 JDK、Maven 多版本安装与切换
- **一键切换**：通过软链接机制实现版本秒级切换，无需手动改环境变量
- **一次配置**：首次运行自动配置用户环境变量，后续切换版本仅修改软链接
- **跨平台支持**：Windows / Linux / macOS 三平台统一实现
- **零管理员权限**：Windows 使用 Junction Point，普通用户即可运行
- **便携设计**：所有数据集中在安装目录，不占用 C 盘，可拷贝迁移
- **版本自动识别**：添加 JDK/Maven 时自动解析版本号，也支持手动指定

## 快速开始

### 安装

1. 下载对应平台的可执行文件
2. 放到指定安装目录（如 `E:\DevManager\`）
3. 首次运行按提示完成环境变量配置

### 编译

```bash
# Windows 64位
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -X devmanager/cmd.version=V0.1.0" -o stormmoondms-windows.exe

# Linux 64位
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X devmanager/cmd.version=V0.1.0" -o stormmoondms-linux

# macOS 64位
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w -X devmanager/cmd.version=V0.1.0" -o stormmoondms-macos
```

### 常用命令

```bash
# 查看版本
stormmoondms version

# 配置环境变量（首次运行或手动触发）
stormmoondms env setup

# JDK 版本管理
stormmoondms jdk add --path <JDK路径>              # 自动识别版本号
stormmoondms jdk add --path <JDK路径> --version 17  # 手动指定版本号
stormmoondms jdk list                                # 列出所有版本（* 标记当前版本）
stormmoondms jdk use <版本号>                         # 切换版本
stormmoondms jdk uninstall <版本号>                   # 卸载版本（当前版本禁止卸载）

# Maven 版本管理
stormmoondms maven add --path <Maven路径>
stormmoondms maven list
stormmoondms maven use <版本号>
stormmoondms maven uninstall <版本号>
```

## 安装目录结构

程序首次运行自动创建以下目录结构：

```
E:\DevManager\                          # 安装目录（示例）
├── stormmoondms.exe                      # 主程序
├── config.json                         # 配置文件（自动生成）
├── versions\                           # 版本文件目录（当前为引用模式）
│   ├── jdk\
│   └── maven\
├── symlinks\                           # 软链接目录
│   ├── jdk -> 实际JDK路径               # 当前使用的JDK版本
│   └── maven -> 实际Maven路径            # 当前使用的Maven版本
└── logs\                               # 日志目录
    └── stormmoondms.log
```

## 环境变量

程序首次运行检测并引导配置以下用户环境变量：

| 环境变量 | 值 | 说明 |
|:---|:---|:---|
| `STORMMOON_HOME` | 安装目录 | 程序安装路径 |
| `JAVA_HOME` | 安装目录/symlinks/jdk | 指向当前JDK版本软链接 |
| `MAVEN_HOME` | 安装目录/symlinks/maven | 指向当前Maven版本软链接 |
| `PATH` | 追加相关路径 | 使命令行全局可用 |

- **Windows**：写入注册表 `HKEY_CURRENT_USER\Environment`，广播 `WM_SETTINGCHANGE`
- **Linux/macOS**：写入 Shell 配置文件（.bashrc/.zshrc/config.fish）

## 技术栈

| 分类 | 选型 |
|:---|:---|
| 开发语言 | Go 1.26 |
| CLI 框架 | github.com/spf13/cobra |
| 颜色输出 | github.com/fatih/color |
| 环境变量配置（Windows） | golang.org/x/sys/windows/registry |
| 配置文件 | JSON（标准库） |
| 软链接 | 系统原生 API（Windows Junction Point / Unix symlink） |
| 编译方式 | CGO_ENABLED=0 静态编译，单文件无依赖 |

## 项目结构

```
devmanager/
├── main.go                  # 程序入口
├── go.mod
├── go.sum
├── cmd/                     # 命令层
│   ├── root.go              # 根命令 + env setup
│   ├── jdk.go               # JDK 子命令组
│   └── maven.go             # Maven 子命令组
├── internal/                # 内部私有代码
│   ├── config/config.go     # 配置管理
│   ├── model/model.go       # 数据模型
│   ├── symlink/             # 软链接（跨平台）
│   │   ├── symlink.go
│   │   ├── symlink_windows.go
│   │   └── symlink_unix.go
│   ├── envcheck/env.go      # 环境检测
│   ├── envsetup/            # 环境变量配置
│   │   ├── env_setup.go
│   │   ├── env_setup_windows.go
│   │   └── env_setup_unix.go
│   └── service/              # 业务服务
│       ├── jdk_service.go
│       └── maven_service.go
└── pkg/                     # 公共工具
    ├── logger.go            # 彩色日志
    └── util.go              # 版本识别+文件工具
```

## 扩展新运行环境

按架构设计，新增 Python / Node 版本管理只需：

1. `cmd/python.go` — 注册 cobra 子命令
2. `internal/service/python_service.go` — 实现 Add/List/Use/Uninstall

无需修改 config、symlink、env 等底层代码。

## License

MIT