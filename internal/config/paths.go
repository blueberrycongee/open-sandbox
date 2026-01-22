package config

import (
	"errors"
	"os"
	"path/filepath"
)

const (
	HostWorkspacePath      = "D:\\Desktop\\sandbox\\open-sandbox\\workspace"
	ContainerWorkspacePath = "/workspace"
	CacheRoot              = "D:\\Desktop\\sandbox\\open-sandbox\\.cache"
	LogsRoot               = "D:\\Desktop\\sandbox\\open-sandbox\\logs"
	BuildRoot              = "D:\\Desktop\\sandbox\\open-sandbox\\build"
)

func EnsureWorkspace() error {
	if !filepath.IsAbs(HostWorkspacePath) {
		return errors.New("workspace path must be absolute")
	}
	return os.MkdirAll(HostWorkspacePath, 0755)
}
