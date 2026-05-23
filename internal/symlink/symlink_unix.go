//go:build !windows

package symlink

import "fmt"

func createJunctionPoint(linkPath string, targetPath string) error {
	return fmt.Errorf("Junction Point 仅在 Windows 平台可用")
}