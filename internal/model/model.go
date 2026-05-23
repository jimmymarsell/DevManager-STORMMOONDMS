package model

type Config struct {
	JdkConfig   EnvConfig `json:"jdk"`
	MavenConfig EnvConfig `json:"maven"`
}

type EnvConfig struct {
	Current     string                 `json:"current"`
	SymlinkPath string                 `json:"symlinkPath"`
	Versions    map[string]VersionInfo `json:"versions"`
}

type VersionInfo struct {
	Version string `json:"version"`
	Path    string `json:"path"`
}

func NewDefaultConfig() *Config {
	return &Config{
		JdkConfig: EnvConfig{
			Current:     "",
			SymlinkPath: "",
			Versions:    make(map[string]VersionInfo),
		},
		MavenConfig: EnvConfig{
			Current:     "",
			SymlinkPath: "",
			Versions:    make(map[string]VersionInfo),
		},
	}
}