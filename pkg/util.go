package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
)

func DetectJdkVersion(jdkPath string) (string, error) {
	var javaExe string
	if runtime.GOOS == "windows" {
		javaExe = filepath.Join(jdkPath, "bin", "java.exe")
	} else {
		javaExe = filepath.Join(jdkPath, "bin", "java")
	}

	if _, err := os.Stat(javaExe); os.IsNotExist(err) {
		return "", fmt.Errorf("未找到 java 可执行文件: %s", javaExe)
	}

	cmd := exec.Command(javaExe, "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("执行 java -version 失败: %v", err)
	}

	version := parseJdkVersion(string(output))
	if version == "" {
		return "", fmt.Errorf("无法从输出中解析版本号")
	}
	return version, nil
}

func DetectMavenVersion(mavenPath string) (string, error) {
	var mvnExe string
	if runtime.GOOS == "windows" {
		mvnExe = filepath.Join(mavenPath, "bin", "mvn.cmd")
		if _, err := os.Stat(mvnExe); os.IsNotExist(err) {
			mvnExe = filepath.Join(mavenPath, "bin", "mvn")
		}
	} else {
		mvnExe = filepath.Join(mavenPath, "bin", "mvn")
	}

	if _, err := os.Stat(mvnExe); os.IsNotExist(err) {
		return "", fmt.Errorf("未找到 mvn 可执行文件: %s", mvnExe)
	}

	cmd := exec.Command(mvnExe, "-v")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("执行 mvn -v 失败: %v", err)
	}

	version := parseMavenVersion(string(output))
	if version == "" {
		return "", fmt.Errorf("无法从输出中解析版本号")
	}
	return version, nil
}

func parseJdkVersion(output string) string {
	// Match versions like 1.8.0_391, 17.0.9, 21, 11.0.2
	re := regexp.MustCompile(`(\d+\.\d+(?:\.\d+)?(?:_\d+)?)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		return matches[1]
	}

	// Try matching a single major version number like "21"
	re2 := regexp.MustCompile(`version "(\d+)"`)
	matches2 := re2.FindStringSubmatch(output)
	if len(matches2) > 1 {
		return matches2[1]
	}

	return ""
}

func parseMavenVersion(output string) string {
	re := regexp.MustCompile(`Apache Maven\s+(\d+\.\d+(?:\.\d+)?)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}