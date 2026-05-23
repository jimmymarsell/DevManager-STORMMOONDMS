package symlink

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func CreateSymlink(linkPath string, targetPath string) error {
	if err := os.MkdirAll(filepath.Dir(linkPath), 0755); err != nil {
		return fmt.Errorf("创建软链接目录失败: %v", err)
	}

	_ = RemoveSymlink(linkPath)

	if runtime.GOOS == "windows" {
		return createWindowsSymlink(linkPath, targetPath)
	}

	if err := os.Symlink(targetPath, linkPath); err != nil {
		return fmt.Errorf("创建软链接失败: %v", err)
	}
	return nil
}

func createWindowsSymlink(linkPath string, targetPath string) error {
	targetInfo, err := os.Stat(targetPath)
	if err != nil {
		return fmt.Errorf("目标路径不存在: %v", err)
	}

	if targetInfo.IsDir() {
		if err := createJunctionPoint(linkPath, targetPath); err != nil {
			if err2 := os.Symlink(targetPath, linkPath); err2 != nil {
				return fmt.Errorf("创建软链接失败（Junction Point 和 Symbolic Link 均失败）\n  Junction Point 错误: %v\n  Symbolic Link 错误: %v\n  请尝试以管理员身份运行或开启 Windows 开发者模式", err, err2)
			}
		}
	} else {
		if err := os.Symlink(targetPath, linkPath); err != nil {
			return fmt.Errorf("创建软链接失败: %v（可能需要管理员权限或开发者模式）", err)
		}
	}
	return nil
}

func RemoveSymlink(linkPath string) error {
	_, err := os.Lstat(linkPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("检测软链接失败: %v", err)
	}

	if err := os.Remove(linkPath); err != nil {
		return fmt.Errorf("删除软链接失败: %v", err)
	}
	return nil
}

func CheckSymlink(linkPath string) bool {
	_, err := os.Readlink(linkPath)
	return err == nil
}

func GetSymlinkTarget(linkPath string) (string, error) {
	target, err := os.Readlink(linkPath)
	if err != nil {
		return "", fmt.Errorf("读取软链接目标失败: %v", err)
	}
	return target, nil
}