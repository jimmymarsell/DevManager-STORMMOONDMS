# DevManager (STORMMOONDMS)

开发环境版本管理器 - 轻量级跨平台命令行工具，用于管理 JDK、Maven 等多版本开发环境。

## 功能特性

- **多版本管理**：支持 JDK、Maven 等开发环境的多版本安装与切换
- **一键切换**：通过软链接机制实现版本秒级切换，无需修改环境变量
- **跨平台支持**：Windows 10+、Ubuntu 20.04+、macOS 12+
- **零管理员权限**：Windows 使用 Junction Point，普通用户即可运行
- **便携设计**：所有数据集中在用户指定安装目录，不占用 C 盘，可拷贝迁移

## 快速开始

### 安装

1. 下载对应平台的可执行文件
2. 放到指定安装目录（如 `E:\DevManager\`）
3. 首次运行自动配置用户环境变量

### 常用命令

```bash
# 添加本地版本
stormmoondms jdk add E:\Java\jdk-17.0.9

# 切换版本
stormmoondms jdk use 17.0.9

# 列出所有版本
stormmoondms jdk list

# 卸载版本
stormmoondms jdk uninstall 1.8.0_391
```

## 目录结构

```
E:\DevManager\
├── stormmoondms.exe      # 主程序
├── config.json           # 配置文件
├── versions\             # 版本文件目录
│   ├── jdk\
│   └── maven\
├── symlinks\             # 软链接目录
│   ├── jdk -> ...
│   └── maven -> ...
└── logs\                 # 日志目录
```

## 技术栈

- **开发语言**：Golang 1.26
- **CLI 框架**：cobra
- **编译方式**：纯静态编译，单文件无依赖

## License

MIT
