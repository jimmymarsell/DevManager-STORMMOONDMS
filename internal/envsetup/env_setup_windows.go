package envsetup

import (
	"fmt"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

func SetupWindowsEnv(installDir string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, "Environment", registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("打开注册表失败，可能需要管理员权限: %v", err)
	}
	defer k.Close()

	jdkPath := filepath.Join(installDir, "symlinks", "jdk")
	mavenPath := filepath.Join(installDir, "symlinks", "maven")

	if err := k.SetStringValue("STORMMOON_HOME", installDir); err != nil {
		return fmt.Errorf("设置 STORMMOON_HOME 失败: %v", err)
	}
	if err := k.SetStringValue("JAVA_HOME", jdkPath); err != nil {
		return fmt.Errorf("设置 JAVA_HOME 失败: %v", err)
	}
	if err := k.SetStringValue("MAVEN_HOME", mavenPath); err != nil {
		return fmt.Errorf("设置 MAVEN_HOME 失败: %v", err)
	}

	currentPath, _, _ := k.GetStringValue("PATH")
	newPaths := []string{
		"%STORMMOON_HOME%",
		"%JAVA_HOME%\\bin",
		"%MAVEN_HOME%\\bin",
	}

	var pathParts []string
	if currentPath != "" {
		pathParts = strings.Split(currentPath, ";")
	}
	for _, np := range newPaths {
		found := false
		for _, existing := range pathParts {
			if strings.EqualFold(existing, np) {
				found = true
				break
			}
		}
		if !found {
			pathParts = append(pathParts, np)
		}
	}
	newPathValue := strings.Join(pathParts, ";")

	if err := k.SetExpandStringValue("PATH", newPathValue); err != nil {
		return fmt.Errorf("设置 PATH 失败: %v", err)
	}

	broadcastSettingChange()
	return nil
}

func broadcastSettingChange() {
	user32 := syscall.NewLazyDLL("user32.dll")
	sendMessageTimeout := user32.NewProc("SendMessageTimeoutW")

	envPtr, _ := syscall.UTF16PtrFromString("Environment")
	sendMessageTimeout.Call(
		0xFFFF, // HWND_BROADCAST
		WM_SETTINGCHANGE,
		0,
		uintptr(unsafe.Pointer(envPtr)),
		SMTO_ABORTIFHUNG,
		5000,
		0,
	)
}

const (
	WM_SETTINGCHANGE = 0x001A
	SMTO_ABORTIFHUNG = 0x0002
)