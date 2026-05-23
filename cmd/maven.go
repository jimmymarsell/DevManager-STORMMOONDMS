package cmd

import (
	"devmanager/internal/service"

	"github.com/spf13/cobra"
)

var mavenCmd = &cobra.Command{
	Use:   "maven",
	Short: "Maven 版本管理",
	Long:  "管理本地 Maven 版本：添加、列表、切换、卸载",
}

var mavenAddCmd = &cobra.Command{
	Use:   "add",
	Short: "添加 Maven 版本",
	Long:  "添加本地 Maven 版本到管理列表，支持自动识别版本号",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		ver, _ := cmd.Flags().GetString("version")
		if err := service.MavenAdd(path, ver); err != nil {
			cmd.PrintErrln(err)
		}
	},
}

var mavenListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有 Maven 版本",
	Long:  "列出已管理的所有 Maven 版本，当前使用版本标记 *",
	Run: func(cmd *cobra.Command, args []string) {
		if err := service.MavenList(); err != nil {
			cmd.PrintErrln(err)
		}
	},
}

var mavenUseCmd = &cobra.Command{
	Use:   "use [version]",
	Short: "切换 Maven 版本",
	Long:  "切换当前使用的 Maven 版本",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := service.MavenUse(args[0]); err != nil {
			cmd.PrintErrln(err)
		}
	},
}

var mavenUninstallCmd = &cobra.Command{
	Use:   "uninstall [version]",
	Short: "卸载 Maven 版本",
	Long:  "从管理列表中移除指定 Maven 版本",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := service.MavenUninstall(args[0]); err != nil {
			cmd.PrintErrln(err)
		}
	},
}

func init() {
	mavenAddCmd.Flags().StringP("path", "p", "", "Maven 本地路径（必填）")
	mavenAddCmd.Flags().StringP("version", "v", "", "手动指定版本号（可选，默认自动识别）")
	_ = mavenAddCmd.MarkFlagRequired("path")

	mavenCmd.AddCommand(mavenAddCmd)
	mavenCmd.AddCommand(mavenListCmd)
	mavenCmd.AddCommand(mavenUseCmd)
	mavenCmd.AddCommand(mavenUninstallCmd)
	rootCmd.AddCommand(mavenCmd)
}