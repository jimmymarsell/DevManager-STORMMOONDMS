package envcheck

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type EnvStatus struct {
	StormmoonHome bool
	JavaHome      bool
	MavenHome     bool
	PathConfig    bool
	AllConfiged   bool
}

func CheckEnvVars(installDir string) *EnvStatus {
	status := &EnvStatus{}

	stormmoonHome := os.Getenv("STORMMOON_HOME")
	status.StormmoonHome = stormmoonHome != "" && stormmoonHome == installDir

	javaHome := os.Getenv("JAVA_HOME")
	expectedJavaHome := filepath.Join(installDir, "symlinks", "jdk")
	status.JavaHome = javaHome != "" && strings.EqualFold(javaHome, expectedJavaHome)

	mavenHome := os.Getenv("MAVEN_HOME")
	expectedMavenHome := filepath.Join(installDir, "symlinks", "maven")
	status.MavenHome = mavenHome != "" && strings.EqualFold(mavenHome, expectedMavenHome)

	pathEnv := os.Getenv("PATH")
	pathLower := strings.ToLower(pathEnv)
	installDirLower := strings.ToLower(installDir)
	status.PathConfig = strings.Contains(pathLower, installDirLower) ||
		strings.Contains(pathLower, "stormmoon")

	status.AllConfiged = status.StormmoonHome && status.JavaHome && status.MavenHome

	return status
}

func GenerateGuideText(status *EnvStatus, installDir string) string {
	if status.AllConfiged {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("检测到您尚未配置 DevManager 环境变量，是否自动配置？[Y/n]\n\n")
	sb.WriteString("自动配置将完成以下操作：\n")
	sb.WriteString("1. 设置 STORMMOON_HOME 环境变量\n")
	sb.WriteString("2. 设置 JAVA_HOME 环境变量\n")
	sb.WriteString("3. 设置 MAVEN_HOME 环境变量\n")
	sb.WriteString("4. 将相关路径添加到 PATH\n\n")

	if runtime.GOOS == "windows" {
		sb.WriteString("配置后，新终端窗口将自动生效。\n")
	} else {
		sb.WriteString("配置后，请重新打开终端或执行 source 命令使配置生效。\n")
	}

	return sb.String()
}

func GenerateManualGuide(installDir string) string {
	var sb strings.Builder

	sb.WriteString("=== 手动配置环境变量指引 ===\n\n")

	if runtime.GOOS == "windows" {
		sb.WriteString(fmt.Sprintf("1. 设置 STORMMOON_HOME = %s\n", installDir))
		sb.WriteString(fmt.Sprintf("2. 设置 JAVA_HOME = %s\n", filepath.Join(installDir, "symlinks", "jdk")))
		sb.WriteString(fmt.Sprintf("3. 设置 MAVEN_HOME = %s\n", filepath.Join(installDir, "symlinks", "maven")))
		sb.WriteString("4. 将以下路径添加到 PATH:\n")
		sb.WriteString("   %STORMMOON_HOME%\n")
		sb.WriteString("   %JAVA_HOME%\\bin\n")
		sb.WriteString("   %MAVEN_HOME%\\bin\n")
		sb.WriteString("\n操作方式：系统属性 → 环境变量 → 用户变量\n")
	} else {
shell := detectShell()
	var configFile string
	switch shell {
		case "zsh":
			configFile = "~/.zshrc"
		case "fish":
			configFile = "~/.config/fish/config.fish"
		default:
			configFile = "~/.bashrc"
		}

		sb.WriteString(fmt.Sprintf("1. 在配置文件 %s 中添加以下内容：\n\n", configFile))
		sb.WriteString(fmt.Sprintf("   export STORMMOON_HOME=%s\n", installDir))
		sb.WriteString(fmt.Sprintf("   export JAVA_HOME=%s\n", filepath.Join(installDir, "symlinks", "jdk")))
		sb.WriteString(fmt.Sprintf("   export MAVEN_HOME=%s\n", filepath.Join(installDir, "symlinks", "maven")))

		if shell == "fish" {
			sb.WriteString("   set -gx PATH $STORMMOON_HOME $JAVA_HOME/bin $MAVEN_HOME/bin $PATH\n")
		} else {
			sb.WriteString("   export PATH=$STORMMOON_HOME:$JAVA_HOME/bin:$MAVEN_HOME/bin:$PATH\n")
		}
		sb.WriteString(fmt.Sprintf("\n2. 执行 source %s 使配置生效\n", configFile))
	}

	return sb.String()
}

func detectShell() string {
	shell := os.Getenv("SHELL")
	if strings.Contains(shell, "zsh") {
		return "zsh"
	}
	if strings.Contains(shell, "fish") {
		return "fish"
	}

	home, _ := os.UserHomeDir()
	if _, err := os.Stat(filepath.Join(home, ".zshrc")); err == nil {
		return "zsh"
	}

	return "bash"
}