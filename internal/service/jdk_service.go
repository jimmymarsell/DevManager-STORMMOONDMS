package service

import (
	"devmanager/internal/config"
	"devmanager/internal/model"
	"devmanager/internal/symlink"
	"devmanager/pkg"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func JdkAdd(path string, manualVersion string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("路径解析失败: %v", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("路径不存在: %s", absPath)
	}

	var javaExe string
	if runtime.GOOS == "windows" {
		javaExe = filepath.Join(absPath, "bin", "java.exe")
	} else {
		javaExe = filepath.Join(absPath, "bin", "java")
	}

	if _, err := os.Stat(javaExe); os.IsNotExist(err) {
		return fmt.Errorf("路径不是合法的 JDK 目录（未找到 %s）", javaExe)
	}

	version := manualVersion
	if version == "" {
		version, err = pkg.DetectJdkVersion(absPath)
		if err != nil {
			return fmt.Errorf("无法自动识别版本号，请使用 --version 手动指定: %v", err)
		}
	}

	cfg := config.GlobalConfig
	if cfg == nil {
		return fmt.Errorf("配置未加载")
	}

	if _, exists := cfg.JdkConfig.Versions[version]; exists {
		return fmt.Errorf("JDK 版本 %s 已存在", version)
	}

	cfg.JdkConfig.Versions[version] = model.VersionInfo{
		Version: version,
		Path:    absPath,
	}

	config.UpdateJdkSymlinkPath(cfg)

	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}

	pkg.PrintSuccess("✔ JDK 版本添加成功：%s", version)
	pkg.PrintInfo("  路径：%s", absPath)
	return nil
}

func JdkList() error {
	cfg := config.GlobalConfig
	if cfg == nil {
		return fmt.Errorf("配置未加载")
	}

	if len(cfg.JdkConfig.Versions) == 0 {
		pkg.PrintInfo("暂无已管理的 JDK 版本")
		return nil
	}

	for ver, info := range cfg.JdkConfig.Versions {
		prefix := "  "
		if ver == cfg.JdkConfig.Current {
			prefix = "* "
		}

		valid := pkg.IsDir(info.Path)
		if valid {
			fmt.Printf("%s%s\n", prefix, ver)
		} else {
			pkg.PrintWarning("%s%s (路径无效)", prefix, ver)
		}
	}

	return nil
}

func JdkUse(version string) error {
	cfg := config.GlobalConfig
	if cfg == nil {
		return fmt.Errorf("配置未加载")
	}

	info, exists := cfg.JdkConfig.Versions[version]
	if !exists {
		pkg.PrintError("✘ JDK 版本 %s 不存在，可用版本：", version)
		for ver := range cfg.JdkConfig.Versions {
			fmt.Printf("  %s\n", ver)
		}
		return nil
	}

	if !pkg.IsDir(info.Path) {
		return fmt.Errorf("JDK 版本 %s 的路径不存在: %s", version, info.Path)
	}

	installDir := config.GetInstallDir()
	linkPath := filepath.Join(installDir, "symlinks", "jdk")

	if symlink.CheckSymlink(linkPath) {
		if err := symlink.RemoveSymlink(linkPath); err != nil {
			return fmt.Errorf("删除旧软链接失败: %v", err)
		}
	}

	if err := symlink.CreateSymlink(linkPath, info.Path); err != nil {
		return fmt.Errorf("创建软链接失败: %v", err)
	}

	cfg.JdkConfig.Current = version
	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}

	pkg.PrintSuccess("✔ JDK 版本切换成功：%s", version)
	pkg.PrintInfo("  请重新打开终端窗口使环境变量生效")
	return nil
}

func JdkUninstall(version string) error {
	cfg := config.GlobalConfig
	if cfg == nil {
		return fmt.Errorf("配置未加载")
	}

	if version == cfg.JdkConfig.Current {
		return fmt.Errorf("无法卸载当前正在使用的 JDK 版本 %s，请先切换到其他版本", version)
	}

	_, exists := cfg.JdkConfig.Versions[version]
	if !exists {
		return fmt.Errorf("JDK 版本 %s 不存在", version)
	}

	delete(cfg.JdkConfig.Versions, version)

	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}

	pkg.PrintSuccess("✔ JDK 版本 %s 已卸载", version)
	return nil
}