# DevManager (STORMMOON)

开发环境版本管理器 — 轻量级跨平台命令行工具，支持 JDK、Maven 等多版本开发环境的添加、切换与卸载。

---

## 目录

- [功能特性](#功能特性)
- [工作原理](#工作原理)
- [快速开始](#快速开始)
- [命令详解](#命令详解)
- [配置文件](#配置文件)
- [环境变量](#环境变量)
- [安装目录结构](#安装目录结构)
- [编译构建](#编译构建)
- [技术栈](#技术栈)
- [项目结构](#项目结构)
- [架构设计](#架构设计)
- [扩展新运行环境](#扩展新运行环境)
- [常见问题](#常见问题)
- [ License](#license)

---

## 功能特性

| 特性 | 说明 |
|:---|:---|
| 多版本管理 | 支持 JDK、Maven 多版本安装与切换 |
| 一键切换 | 通过软链接机制实现版本秒级切换，无需手动改环境变量 |
| 一次配置 | 首次运行自动配置用户环境变量，后续切换版本仅修改软链接 |
| 跨平台支持 | Windows / Linux / macOS 三平台统一实现 |
| 零管理员权限 | Windows 使用 Junction Point，普通用户即可运行 |
| 便携设计 | 所有数据集中在安装目录，不占用 C 盘，可拷贝迁移 |
| 版本自动识别 | 添加 JDK/Maven 时自动解析版本号，也支持手动指定 |
| 彩色终端输出 | 成功绿色、错误红色、警告黄色，一目了然 |
| 日志文件记录 | 所有操作同时输出到 `logs/stormmoondms.log` |

---

## 工作原理

DevManager 采用 **「固定环境变量 + 全局统一软链接」** 方案：

```
环境变量 JAVA_HOME  ──→  E:\DevManager\symlinks\jdk  ──→  实际JDK版本目录
                                                        ↑
                                    切换版本时仅修改软链接指向
                                    环境变量 JAVA_HOME 永远不变
```

1. 首次运行自动配置 `JAVA_HOME`、`MAVEN_HOME`、`STORMMOON_HOME` 等环境变量，指向 `symlinks` 目录下的软链接
2. 切换版本时，程序仅删除旧软链接并创建新软链接指向目标版本目录
3. 环境变量始终指向 `symlinks` 目录，无需修改，无需重启电脑（新开终端即可生效）

**Windows 平台**：使用 Junction Point（目录联接），普通用户权限即可创建和删除。

**Linux/macOS 平台**：使用系统原生 symlink。

---

## 快速开始

### 1. 安装

1. 下载对应平台的可执行文件
2. 放到任意安装目录（如 `E:\DevManager\`），目录路径不含中文和空格
3. 在终端中运行程序

### 2. 首次运行

首次运行时程序会自动检测环境变量，提示是否自动配置：

```
检测到您尚未配置 DevManager 环境变量，是否自动配置？[Y/n]

自动配置将完成以下操作：
1. 设置 STORMMOON_HOME 环境变量
2. 设置 JAVA_HOME 环境变量
3. 设置 MAVEN_HOME 环境变量
4. 将相关路径添加到 PATH

配置后，新终端窗口将自动生效。
请输入 (Y/n):
```

输入 `Y` 或直接回车即可自动配置，输入 `n` 则输出手动配置指引。

### 3. 添加版本

```bash
# 添加已安装的 JDK（自动识别版本号）
stormmoondms jdk add --path D:\Java\jdk-17.0.9

# 添加已安装的 Maven（自动识别版本号）
stormmoondms maven add --path D:\Maven\apache-maven-3.8.8

# 如果版本号识别失败，可手动指定
stormmoondms jdk add --path D:\Java\jdk-17.0.9 --version 17.0.9
```

### 4. 切换版本

```bash
stormmoondms jdk use 17.0.9
stormmoondms maven use 3.8.8
```

切换成功后重新打开终端窗口，即可使用新版本。

---

## 命令详解

### 全局命令

| 命令 | 说明 |
|:---|:---|
| `stormmoondms` | 显示欢迎信息和帮助 |
| `stormmoondms version` | 显示版本号 |
| `stormmoondms env setup` | 手动触发环境变量配置 |

### JDK 版本管理

| 命令 | 说明 |
|:---|:---|
| `stormmoondms jdk add --path <路径>` | 添加 JDK 版本，自动识别版本号 |
| `stormmoondms jdk add --path <路径> --version <版本>` | 手动指定版本号 |
| `stormmoondms jdk list` | 列出所有已管理的 JDK 版本 |
| `stormmoondms jdk use <版本号>` | 切换到指定 JDK 版本 |
| `stormmoondms jdk uninstall <版本号>` | 卸载指定 JDK 版本（当前使用版本禁止卸载） |

### Maven 版本管理

| 命令 | 说明 |
|:---|:---|
| `stormmoondms maven add --path <路径>` | 添加 Maven 版本，自动识别版本号 |
| `stormmoondms maven add --path <路径> --version <版本>` | 手动指定版本号 |
| `stormmoondms maven list` | 列出所有已管理的 Maven 版本 |
| `stormmoondms maven use <版本号>` | 切换到指定 Maven 版本 |
| `stormmoondms maven uninstall <版本号>` | 卸载指定 Maven 版本（当前使用版本禁止卸载） |

### 输出示例

```
# jdk list
  1.8.0_391
* 17.0.9           ← * 标记当前使用版本

# jdk use 17.0.9
✔ JDK 版本切换成功：17.0.9
  请重新打开终端窗口使环境变量生效

# jdk add --path D:\Java\jdk-17.0.9
✔ JDK 版本添加成功：17.0.9
  路径：D:\Java\jdk-17.0.9

# 错误提示示例
✘ JDK 版本 11.0.1 不存在，可用版本：
  1.8.0_391
  17.0.9
```

---

## 配置文件

程序首次运行时在安装目录下自动生成 `config.json`：

```json
{
  "jdk": {
    "current": "17.0.9",
    "symlinkPath": "E:\\DevManager\\symlinks\\jdk",
    "versions": {
      "1.8.0_391": {
        "version": "1.8.0_391",
        "path": "D:\\Java\\jdk1.8.0_391"
      },
      "17.0.9": {
        "version": "17.0.9",
        "path": "D:\\Java\\jdk-17.0.9"
      }
    }
  },
  "maven": {
    "current": "3.8.8",
    "symlinkPath": "E:\\DevManager\\symlinks\\maven",
    "versions": {
      "3.6.3": {
        "version": "3.6.3",
        "path": "D:\\Maven\\apache-maven-3.6.3"
      },
      "3.8.8": {
        "version": "3.8.8",
        "path": "D:\\Maven\\apache-maven-3.8.8"
      }
    }
  }
}
```

- 配置文件损坏时自动备份为 `config.json.bak` 并重置为空配置
- 每次修改立即落盘写入，无需手动保存

---

## 环境变量

程序首次运行检测并引导配置以下用户环境变量：

| 环境变量 | 值 | 说明 |
|:---|:---|:---|
| `STORMMOON_HOME` | 安装目录 | 程序安装路径 |
| `JAVA_HOME` | 安装目录/symlinks/jdk | 指向当前 JDK 版本软链接 |
| `MAVEN_HOME` | 安装目录/symlinks/maven | 指向当前 Maven 版本软链接 |
| `PATH` | 追加 `%STORMMOON_HOME%;%JAVA_HOME%\bin;%MAVEN_HOME%\bin` | 命令行全局可用 |

### 各平台实现方式

| 平台 | 实现方式 |
|:---|:---|
| Windows | 写入注册表 `HKEY_CURRENT_USER\Environment`，广播 `WM_SETTINGCHANGE` |
| Linux/macOS | 写入 Shell 配置文件（`.bashrc` / `.zshrc` / `config.fish`） |

- 仅修改**用户**环境变量，不修改系统环境变量，无需管理员权限
- PATH 中已存在的路径不会重复添加
- 配置后需重新打开终端窗口使环境变量生效

---

## 安装目录结构

程序首次运行自动创建以下目录结构：

```
E:\DevManager\                          # 安装目录（示例）
├── stormmoondms.exe                      # 主程序
├── config.json                         # 配置文件（自动生成）
├── versions\                           # 版本文件目录
│   ├── jdk\                            # （预留给未来版本复制模式）
│   └── maven\                          # （预留给未来版本复制模式）
├── symlinks\                           # 软链接目录（Junction Point）
│   ├── jdk -> D:\Java\jdk-17.0.9       # 当前使用的 JDK 版本
│   └── maven -> D:\Maven\apache-maven-3.8.8  # 当前使用的 Maven 版本
└── logs\                               # 日志目录
    └── stormmoondms.log                # 日志文件
```

- `symlinks/jdk` 和 `symlinks/maven` 是软链接，始终指向当前使用的版本
- 切换版本时仅修改软链接指向，`JAVA_HOME` 和 `MAVEN_HOME` 路径不变
- 完全便携，将整个安装目录拷贝到另一台电脑即可使用

---

## 编译构建

### 前置要求

- Go 1.26 及以上

### 编译命令

```bash
# 当前平台编译（推荐，最简单）
go build -ldflags "-s -w -X devmanager/cmd.version=V0.1.0" -o stormmoondms.exe

# 交叉编译其他平台
# Linux 64位
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -ldflags "-s -w -X devmanager/cmd.version=V0.1.0" -o stormmoondms-linux

# macOS 64位
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -ldflags "-s -w -X devmanager/cmd.version=V0.1.0" -o stormmoondms-macos
```

### 编译优化参数说明

| 参数 | 说明 |
|:---|:---|
| `-s` | 去除符号表，减小体积 |
| `-w` | 去除 DWARF 调试信息，减小体积 |
| `-X devmanager/cmd.version=V0.1.0` | 编译时注入版本号 |

### 编译产物

| 平台 | 产物名称 | 体积 |
|:---|:---|:---|
| Windows | `stormmoondms.exe` | ~3.3 MB |
| Linux | `stormmoondms-linux` | ~3.0 MB |
| macOS | `stormmoondms-macos` | ~3.1 MB |

所有产物均为纯静态编译，无外部依赖，单文件即可运行。

---

## 技术栈

| 分类 | 选型 | 说明 |
|:---|:---|:---|
| 开发语言 | Go 1.26 | 静态编译、跨平台、原生系统 API |
| CLI 框架 | github.com/spf13/cobra | 子命令、参数解析、帮助文档 |
| 颜色输出 | github.com/fatih/color | 终端彩色提示 |
| 环境变量（Windows） | golang.org/x/sys/windows/registry | 操作注册表 `HKEY_CURRENT_USER\Environment` |
| 环境变量（Linux/macOS） | Go 标准库 os、bufio | 修改 Shell 配置文件 |
| 配置文件 | JSON（标准库 encoding/json） | 轻量、无需额外依赖 |
| 软链接 | 系统原生 API | Windows Junction Point / Unix symlink |
| 编译方式 | CGO_ENABLED=0 静态编译 | 单文件、无运行时依赖 |

---

## 项目结构

```
devmanager/
├── main.go                              # 程序入口
├── go.mod                               # Go Module 定义
├── go.sum                               # 依赖校验
│
├── cmd/                                 # 命令层（cobra）
│   ├── root.go                          # 根命令 + env setup + 初始化
│   ├── jdk.go                           # JDK 子命令组
│   └── maven.go                         # Maven 子命令组
│
├── internal/                            # 内部私有代码
│   ├── config/
│   │   └── config.go                    # 配置读写、目录初始化
│   ├── model/
│   │   └── model.go                     # Config/EnvConfig/VersionInfo 结构体
│   ├── symlink/
│   │   ├── symlink.go                   # 软链接统一接口
│   │   ├── symlink_windows.go           # Windows Junction Point 实现
│   │   └── symlink_unix.go              # Unix symlink 桩函数
│   ├── envcheck/
│   │   └── env.go                       # 环境变量检测 + 引导文案
│   ├── envsetup/
│   │   ├── env_setup.go                 # 环境变量配置路由 + Unix 实现
│   │   ├── env_setup_windows.go         # Windows 注册表操作 + WM_SETTINGCHANGE
│   │   └── env_setup_unix.go            # Unix 桩函数
│   └── service/
│       ├── jdk_service.go               # JDK 业务逻辑
│       └── maven_service.go             # Maven 业务逻辑
│
├── pkg/                                 # 公共工具包
│   ├── logger.go                        # 彩色终端输出 + 文件日志
│   └── util.go                          # 版本识别、文件判断
│
└── docs/                                # 项目文档
    └── 开发计划V1.2/                     # 开发计划与验收清单
```

---

## 架构设计

```
┌─────────────────────────────────────────────┐
│                  入口层 main.go               │
├─────────────────────────────────────────────┤
│                命令层 cmd/                    │
│   root.go  │  jdk.go  │  maven.go          │
├─────────────────────────────────────────────┤
│              业务服务层 service/              │
│   jdk_service.go  │  maven_service.go       │
├───────┬───────┬───────────┬─────────────────┤
│config │symlink│ envcheck  │  envsetup        │
│ 配置  │软链接 │ 环境检测  │  环境变量配置      │
├───────┴───────┴───────────┴─────────────────┤
│              模型层 model/                   │
├─────────────────────────────────────────────┤
│              公共工具层 pkg/                 │
│         logger  │  util                     │
└─────────────────────────────────────────────┘
```

**核心原则**：

- 命令层只做参数解析和调用，不含业务逻辑
- 业务服务层实现核心逻辑，调用底层模块
- 核心通用层永久不变、业务无关
- 上层依赖下层，下层禁止依赖上层

---

## 扩展新运行环境

按架构设计，新增 Python / Node 版本管理只需：

1. `cmd/python.go` — 注册 cobra 子命令（add/list/use/uninstall）
2. `internal/service/python_service.go` — 实现 Add/List/Use/Uninstall 四大方法

无需修改 config、symlink、env 等底层代码，完全插拔式扩展。

---

## 常见问题

### 切换版本后终端没有生效？

需要**重新打开终端窗口**。程序修改的是用户环境变量（Windows 注册表 / Unix Shell 配置），已打开的终端不会自动刷新。新开的终端会自动读取最新环境变量。

### Windows 提示创建软链接失败？

- 程序优先使用 Junction Point（普通用户权限即可创建）
- 如果 Junction Point 失败，会自动降级为 Symbolic Link
- 如果两者都失败，请以管理员身份运行或开启 Windows 开发者模式（设置 → 更新和安全 → 开发者选项）
- 部分杀毒软件可能拦截软链接创建，请将程序加入白名单

### 如何手动配置环境变量？

运行 `stormmoondms env setup` 即可自动配置。如果自动配置失败，程序会输出手动配置指引，包括：

- **Windows**：系统属性 → 环境变量 → 用户变量
- **Linux/macOS**：在 Shell 配置文件中添加 export 语句

### 如何卸载版本？

```bash
stormmoondms jdk uninstall <版本号>
```

> 注意：当前正在使用的版本禁止卸载，请先切换到其他版本。

卸载仅从管理列表中移除记录，不会删除本地原始安装目录。

### 如何完全卸载 DevManager？

1. 删除安装目录（如 `E:\DevManager`）
2. 删除用户环境变量：`STORMMOON_HOME`、`JAVA_HOME`、`MAVEN_HOME`
3. 从 PATH 中移除相关路径

---

## License

MIT