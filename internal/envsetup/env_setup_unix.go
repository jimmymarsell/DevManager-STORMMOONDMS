//go:build !windows

package envsetup

import "fmt"

func SetupWindowsEnv(installDir string) error {
	return fmt.Errorf("此功能仅在 Windows 平台可用")
}

func broadcastSettingChange() {}