package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"devmanager/internal/model"
)

var GlobalConfig *model.Config

func GetInstallDir() string {
	dir := os.Getenv("STORMMOON_HOME")
	if dir == "" {
		exePath, err := os.Executable()
		if err != nil {
			return ""
		}
		dir = filepath.Dir(exePath)
	}
	return dir
}

func CheckInstallDirWritable() error {
	installDir := GetInstallDir()
	if installDir == "" {
		return fmt.Errorf("无法确定安装目录")
	}

	testFile := filepath.Join(installDir, ".stormmoon_write_test")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return fmt.Errorf("安装目录无写权限: %s\n请修改目录权限或将程序移到有写权限的目录", installDir)
	}
	os.Remove(testFile)
	return nil
}

func InitDirectories() error {
	installDir := GetInstallDir()
	if installDir == "" {
		return fmt.Errorf("无法确定安装目录")
	}

	if err := CheckInstallDirWritable(); err != nil {
		return err
	}

	dirs := []string{
		filepath.Join(installDir, "versions", "jdk"),
		filepath.Join(installDir, "versions", "maven"),
		filepath.Join(installDir, "symlinks"),
		filepath.Join(installDir, "logs"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录失败 [%s]: %v", dir, err)
		}
	}

	return nil
}

func GetConfigPath() string {
	return filepath.Join(GetInstallDir(), "config.json")
}

func LoadConfig() (*model.Config, error) {
	configPath := GetConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return initConfigFile(configPath)
		}
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var cfg model.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		backupPath := configPath + ".bak"
		_ = os.Rename(configPath, backupPath)
		return initConfigFile(configPath)
	}

	if cfg.JdkConfig.Versions == nil {
		cfg.JdkConfig.Versions = make(map[string]model.VersionInfo)
	}
	if cfg.MavenConfig.Versions == nil {
		cfg.MavenConfig.Versions = make(map[string]model.VersionInfo)
	}

	GlobalConfig = &cfg
	return &cfg, nil
}

func SaveConfig(cfg *model.Config) error {
	configPath := GetConfigPath()

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	GlobalConfig = cfg
	return nil
}

func UpdateJdkSymlinkPath(cfg *model.Config) {
	installDir := GetInstallDir()
	cfg.JdkConfig.SymlinkPath = filepath.Join(installDir, "symlinks", "jdk")
}

func UpdateMavenSymlinkPath(cfg *model.Config) {
	installDir := GetInstallDir()
	cfg.MavenConfig.SymlinkPath = filepath.Join(installDir, "symlinks", "maven")
}

func initConfigFile(configPath string) (*model.Config, error) {
	cfg := model.NewDefaultConfig()

	UpdateJdkSymlinkPath(cfg)
	UpdateMavenSymlinkPath(cfg)

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("序列化默认配置失败: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return nil, fmt.Errorf("写入默认配置失败: %v", err)
	}

	GlobalConfig = cfg
	return cfg, nil
}