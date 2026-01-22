package config

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	HostWorkspacePath      = "D:\\Desktop\\sandbox\\open-sandbox\\workspace"
	ContainerWorkspacePath = "/workspace"
	CacheRoot              = "D:\\Desktop\\sandbox\\open-sandbox\\.cache"
	LogsRoot               = "D:\\Desktop\\sandbox\\open-sandbox\\logs"
	BuildRoot              = "D:\\Desktop\\sandbox\\open-sandbox\\build"
)

func WorkspacePath() string {
	if runtime.GOOS == "windows" {
		return HostWorkspacePath
	}
	return toWSLPath(HostWorkspacePath)
}

func EnsureWorkspace() error {
	workspace := WorkspacePath()
	if !filepath.IsAbs(workspace) {
		return errors.New("workspace path must be absolute")
	}
	return os.MkdirAll(workspace, 0755)
}

func toWSLPath(path string) string {
	if len(path) < 3 || path[1] != ':' || (path[2] != '\\' && path[2] != '/') {
		return path
	}
	drive := strings.ToLower(string(path[0]))
	rest := strings.ReplaceAll(path[2:], "\\", "/")
	return "/mnt/" + drive + rest
}
