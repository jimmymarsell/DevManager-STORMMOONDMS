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

func MavenAdd(path string, manualVersion string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("路径解析失败: %v", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("路径不存在: %s", absPath)
	}

	var mvnExe string
	if runtime.GOOS == "windows" {
		mvnExe = filepath.Join(absPath, "bin", "mvn.cmd")
		if _, err := os.Stat(mvnExe); os.IsNotExist(err) {
			mvnExe = filepath.Join(absPath, "bin", "mvn")
		}
	} else {
		mvnExe = filepath.Join(absPath, "bin", "mvn")
	}

	if _, err := os.Stat(mvnExe); os.IsNotExist(err) {
		return fmt.Errorf("路径不是合法的 Maven 目录（未找到 %s）", mvnExe)
	}

	version := manualVersion
	if version == "" {
		version, err = pkg.DetectMavenVersion(absPath)
		if err != nil {
			return fmt.Errorf("无法自动识别版本号，请使用 --version 手动指定: %v", err)
		}
	}

	cfg := config.GlobalConfig
	if cfg == nil {
		return fmt.Errorf("配置未加载")
	}

	if _, exists := cfg.MavenConfig.Versions[version]; exists {
		return fmt.Errorf("Maven 版本 %s 已存在", version)
	}

	cfg.MavenConfig.Versions[version] = model.VersionInfo{
		Version: version,
		Path:    absPath,
	}

	config.UpdateMavenSymlinkPath(cfg)

	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}

	pkg.PrintSuccess("✔ Maven 版本添加成功：%s", version)
	pkg.PrintInfo("  路径：%s", absPath)
	return nil
}

func MavenList() error {
	cfg := config.GlobalConfig
	if cfg == nil {
		return fmt.Errorf("配置未加载")
	}

	if len(cfg.MavenConfig.Versions) == 0 {
		pkg.PrintInfo("暂无已管理的 Maven 版本")
		return nil
	}

	for ver, info := range cfg.MavenConfig.Versions {
		prefix := "  "
		if ver == cfg.MavenConfig.Current {
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

func MavenUse(version string) error {
	cfg := config.GlobalConfig
	if cfg == nil {
		return fmt.Errorf("配置未加载")
	}

	info, exists := cfg.MavenConfig.Versions[version]
	if !exists {
		pkg.PrintError("✘ Maven 版本 %s 不存在，可用版本：", version)
		for ver := range cfg.MavenConfig.Versions {
			fmt.Printf("  %s\n", ver)
		}
		return nil
	}

	if !pkg.IsDir(info.Path) {
		return fmt.Errorf("Maven 版本 %s 的路径不存在: %s", version, info.Path)
	}

	installDir := config.GetInstallDir()
	linkPath := filepath.Join(installDir, "symlinks", "maven")

	if symlink.CheckSymlink(linkPath) {
		if err := symlink.RemoveSymlink(linkPath); err != nil {
			return fmt.Errorf("删除旧软链接失败: %v", err)
		}
	}

	if err := symlink.CreateSymlink(linkPath, info.Path); err != nil {
		return fmt.Errorf("创建软链接失败: %v", err)
	}

	cfg.MavenConfig.Current = version
	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}

	pkg.PrintSuccess("✔ Maven 版本切换成功：%s", version)
	pkg.PrintInfo("  请重新打开终端窗口使环境变量生效")
	return nil
}

func MavenUninstall(version string) error {
	cfg := config.GlobalConfig
	if cfg == nil {
		return fmt.Errorf("配置未加载")
	}

	if version == cfg.MavenConfig.Current {
		return fmt.Errorf("无法卸载当前正在使用的 Maven 版本 %s，请先切换到其他版本", version)
	}

	_, exists := cfg.MavenConfig.Versions[version]
	if !exists {
		return fmt.Errorf("Maven 版本 %s 不存在", version)
	}

	delete(cfg.MavenConfig.Versions, version)

	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}

	pkg.PrintSuccess("✔ Maven 版本 %s 已卸载", version)
	return nil
}