package cmd

import (
	"devmanager/internal/service"

	"github.com/spf13/cobra"
)

var jdkCmd = &cobra.Command{
	Use:   "jdk",
	Short: "JDK 版本管理",
	Long:  "管理本地 JDK 版本：添加、列表、切换、卸载",
}

var jdkAddCmd = &cobra.Command{
	Use:   "add",
	Short: "添加 JDK 版本",
	Long:  "添加本地 JDK 版本到管理列表，支持自动识别版本号",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		ver, _ := cmd.Flags().GetString("version")
		if err := service.JdkAdd(path, ver); err != nil {
			cmd.PrintErrln(err)
		}
	},
}

var jdkListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有 JDK 版本",
	Long:  "列出已管理的所有 JDK 版本，当前使用版本标记 *",
	Run: func(cmd *cobra.Command, args []string) {
		if err := service.JdkList(); err != nil {
			cmd.PrintErrln(err)
		}
	},
}

var jdkUseCmd = &cobra.Command{
	Use:   "use [version]",
	Short: "切换 JDK 版本",
	Long:  "切换当前使用的 JDK 版本",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := service.JdkUse(args[0]); err != nil {
			cmd.PrintErrln(err)
		}
	},
}

var jdkUninstallCmd = &cobra.Command{
	Use:   "uninstall [version]",
	Short: "卸载 JDK 版本",
	Long:  "从管理列表中移除指定 JDK 版本",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := service.JdkUninstall(args[0]); err != nil {
			cmd.PrintErrln(err)
		}
	},
}

func init() {
	jdkAddCmd.Flags().StringP("path", "p", "", "JDK 本地路径（必填）")
	jdkAddCmd.Flags().StringP("version", "v", "", "手动指定版本号（可选，默认自动识别）")
	_ = jdkAddCmd.MarkFlagRequired("path")

	jdkCmd.AddCommand(jdkAddCmd)
	jdkCmd.AddCommand(jdkListCmd)
	jdkCmd.AddCommand(jdkUseCmd)
	jdkCmd.AddCommand(jdkUninstallCmd)
	rootCmd.AddCommand(jdkCmd)
}