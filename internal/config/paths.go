package config

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	ContainerWorkspacePath = "/workspace"
	DefaultRootWindows     = "D:\\Desktop\\sandbox\\open-sandbox"
	DefaultRootUnix        = "/workspace"
)

func RootPath() string {
	if value := envPath("SANDBOX_ROOT"); value != "" {
		return value
	}

	if repo := findRepoRoot(); repo != "" {
		return normalizeAbs(repo)
	}

	if runtime.GOOS == "windows" {
		return normalizeAbs(DefaultRootWindows)
	}
	return normalizeAbs(DefaultRootUnix)
}

func WorkspacePath() string {
	if value := envPath("SANDBOX_WORKSPACE"); value != "" {
		return value
	}
	return normalizeAbs(filepath.Join(RootPath(), "workspace"))
}

func CachePath() string {
	if value := envPath("SANDBOX_CACHE_ROOT"); value != "" {
		return value
	}
	return normalizeAbs(filepath.Join(RootPath(), ".cache"))
}

func LogsPath() string {
	if value := envPath("SANDBOX_LOGS_ROOT"); value != "" {
		return value
	}
	return normalizeAbs(filepath.Join(RootPath(), "logs"))
}

func BuildPath() string {
	if value := envPath("SANDBOX_BUILD_ROOT"); value != "" {
		return value
	}
	return normalizeAbs(filepath.Join(RootPath(), "build"))
}

func EnsureWorkspace() error {
	workspace := WorkspacePath()
	if !filepath.IsAbs(workspace) {
		return errors.New("workspace path must be absolute")
	}
	return os.MkdirAll(workspace, 0755)
}

func envPath(key string) string {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return ""
	}
	return normalizeAbs(raw)
}

func normalizeAbs(path string) string {
	normalized := normalizePath(path)
	if filepath.IsAbs(normalized) {
		return normalized
	}
	abs, err := filepath.Abs(normalized)
	if err != nil {
		return normalized
	}
	return abs
}

func normalizePath(path string) string {
	if runtime.GOOS != "windows" {
		return toWSLPath(path)
	}
	return path
}

func toWSLPath(path string) string {
	if len(path) < 3 || path[1] != ':' || (path[2] != '\\' && path[2] != '/') {
		return path
	}
	drive := strings.ToLower(string(path[0]))
	rest := strings.ReplaceAll(path[2:], "\\", "/")
	return "/mnt/" + drive + rest
}

func findRepoRoot() string {
	current, err := os.Getwd()
	if err != nil {
		return ""
	}
	current = normalizeAbs(current)

	for {
		if hasRepoMarker(current) {
			return current
		}
		parent := filepath.Dir(current)
		if parent == current {
			return ""
		}
		current = parent
	}
}

func hasRepoMarker(path string) bool {
	if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
		return true
	}
	if _, err := os.Stat(filepath.Join(path, ".specify")); err == nil {
		return true
	}
	return false
}
