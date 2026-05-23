package cmd

import (
	"devmanager/internal/config"
	"devmanager/internal/envcheck"
	"devmanager/internal/envsetup"
	"devmanager/pkg"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var version = "V0.1.0"

var rootCmd = &cobra.Command{
	Use:   "stormmoondms",
	Short: "DevManager(STORMMOON) - 开发环境版本管理器",
	Long:  "DevManager(STORMMOON) 是一款轻量级跨平台开发环境版本管理器，支持 JDK、Maven 版本管理",
	Run: func(cmd *cobra.Command, args []string) {
		pkg.PrintInfo("欢迎使用 DevManager(STORMMOON) 开发环境版本管理器")
		fmt.Println("使用 --help 查看可用命令")

		installDir := config.GetInstallDir()
		status := envcheck.CheckEnvVars(installDir)
		if !status.AllConfiged {
			handleEnvSetup(status, installDir)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本号",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("DevManager(STORMMOON) %s\n", version)
	},
}

var envSetupCmd = &cobra.Command{
	Use:   "env",
	Short: "环境变量管理",
}

var envSetupSubCmd = &cobra.Command{
	Use:   "setup",
	Short: "配置环境变量",
	Long:  "自动配置 DevManager 所需的环境变量（STORMMOON_HOME、JAVA_HOME、MAVEN_HOME、PATH）",
	Run: func(cmd *cobra.Command, args []string) {
		installDir := config.GetInstallDir()
		if err := envsetup.SetupEnv(installDir); err != nil {
			pkg.PrintError("环境变量配置失败: %v", err)
			fmt.Println()
			fmt.Print(envcheck.GenerateManualGuide(installDir))
			os.Exit(1)
		}
		pkg.PrintSuccess("✔ 环境变量配置成功！")
		pkg.PrintInfo("  请重新打开终端窗口使配置生效")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	envSetupCmd.AddCommand(envSetupSubCmd)
	rootCmd.AddCommand(envSetupCmd)
}

func Execute() {
	installDir := config.GetInstallDir()
	if installDir == "" {
		pkg.PrintError("无法确定安装目录，请设置 STORMMOON_HOME 环境变量")
		os.Exit(1)
	}

	if err := config.InitDirectories(); err != nil {
		pkg.PrintError("初始化目录失败: %v", err)
		os.Exit(1)
	}

	if err := pkg.InitLogger(installDir); err != nil {
		pkg.PrintWarning("日志文件初始化失败: %v", err)
	}
	defer pkg.CloseLogger()

	if _, err := config.LoadConfig(); err != nil {
		pkg.PrintError("加载配置失败: %v", err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func handleEnvSetup(status *envcheck.EnvStatus, installDir string) {
	fmt.Print(envcheck.GenerateGuideText(status, installDir))

	var input string
	fmt.Print("请输入 (Y/n): ")
	fmt.Scanln(&input)
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "" || input == "y" || input == "yes" {
		if err := envsetup.SetupEnv(installDir); err != nil {
			pkg.PrintError("环境变量配置失败: %v", err)
			fmt.Println()
			fmt.Print(envcheck.GenerateManualGuide(installDir))
			return
		}
		pkg.PrintSuccess("✔ 环境变量配置成功！")
		pkg.PrintInfo("  请重新打开终端窗口使配置生效")
	} else {
		fmt.Println()
		fmt.Print(envcheck.GenerateManualGuide(installDir))
	}
}