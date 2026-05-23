package envsetup

import (
	"devmanager/internal/envcheck"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func SetupEnv(installDir string) error {
	if runtime.GOOS == "windows" {
		return SetupWindowsEnv(installDir)
	}
	return SetupUnixEnv(installDir)
}

func SetupUnixEnv(installDir string) error {
	shell := detectShell()
	var configFile string
	home, _ := os.UserHomeDir()

	switch shell {
	case "zsh":
		configFile = filepath.Join(home, ".zshrc")
	case "fish":
		configFile = filepath.Join(home, ".config", "fish", "config.fish")
	default:
		configFile = filepath.Join(home, ".bashrc")
	}

	configDir := filepath.Dir(configFile)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	existing, err := os.ReadFile(configFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	content := string(existing)
	if strings.Contains(content, "STORMMOON_HOME") {
		fmt.Println("环境变量已存在于 Shell 配置文件中，跳过写入。")
		return nil
	}

	var sb strings.Builder
	sb.WriteString("\n# DevManager(STORMMOON) 环境变量配置\n")
	sb.WriteString(fmt.Sprintf("export STORMMOON_HOME=%s\n", installDir))
	sb.WriteString(fmt.Sprintf("export JAVA_HOME=%s\n", filepath.Join(installDir, "symlinks", "jdk")))
	sb.WriteString(fmt.Sprintf("export MAVEN_HOME=%s\n", filepath.Join(installDir, "symlinks", "maven")))

	if shell == "fish" {
		sb.WriteString("set -gx PATH $STORMMOON_HOME $JAVA_HOME/bin $MAVEN_HOME/bin $PATH\n")
	} else {
		sb.WriteString("export PATH=$STORMMOON_HOME:$JAVA_HOME/bin:$MAVEN_HOME/bin:$PATH\n")
	}
	sb.WriteString("# End DevManager config\n")

	f, err := os.OpenFile(configFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开配置文件失败: %v", err)
	}
	defer f.Close()

	if _, err := f.WriteString(sb.String()); err != nil {
		return fmt.Errorf("写入配置内容失败: %v", err)
	}

	fmt.Printf("✔ 环境变量已写入 %s\n", configFile)
	fmt.Printf("  请执行 source %s 或重新打开终端使配置生效\n", configFile)
	return nil
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

func GenerateManualGuide(installDir string) string {
	return envcheck.GenerateManualGuide(installDir)
}