package symlink

import (
	"fmt"
	"os"
	"os/exec"
)

func createJunctionPoint(linkPath string, targetPath string) error {
	// First try os.Symlink - in Go 1.26+ with developer mode this works for directories
	if err := os.Symlink(targetPath, linkPath); err == nil {
		return nil
	}

	// Fallback: use mklink /J which creates a Junction Point without admin privileges
	cmd := exec.Command("cmd", "/c", "mklink", "/J", linkPath, targetPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Junction Point 创建失败: %v\n输出: %s", err, string(output))
	}

	// Verify the junction was created
	if _, err := os.Readlink(linkPath); err != nil {
		return fmt.Errorf("Junction Point 创建后校验失败: %v", err)
	}

	return nil
}