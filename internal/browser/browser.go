package browser

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

var ErrBrowserUnavailable = errors.New("browser unavailable")

type Config struct {
	BinaryPath             string
	UserDataDir            string
	RemoteDebuggingHost    string
	RemoteDebuggingPort    int
	ExistingWebSocketDebug string
	StartupTimeout         time.Duration
	Headless               bool
}

type Service struct {
	config Config

	mu          sync.Mutex
	allocCtx    context.Context
	allocCancel context.CancelFunc
	tabCtx      context.Context
	tabCancel   context.CancelFunc
	cdpURL      string
	cmd         *exec.Cmd
}

type versionInfo struct {
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
}

func DefaultConfig() Config {
	return Config{
		RemoteDebuggingHost: "127.0.0.1",
		RemoteDebuggingPort: 9222,
		StartupTimeout:      15 * time.Second,
		Headless:            false,
	}
}

func NewService(config Config) *Service {
	if config.RemoteDebuggingHost == "" {
		config.RemoteDebuggingHost = "127.0.0.1"
	}
	if config.RemoteDebuggingPort == 0 {
		config.RemoteDebuggingPort = 9222
	}
	if config.StartupTimeout == 0 {
		config.StartupTimeout = 15 * time.Second
	}
	return &Service{config: config}
}

func (service *Service) Start() error {
	return service.ensureStarted()
}

func (service *Service) Close() {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.tabCancel != nil {
		service.tabCancel()
		service.tabCancel = nil
	}
	if service.allocCancel != nil {
		service.allocCancel()
		service.allocCancel = nil
	}
	service.allocCtx = nil
	service.tabCtx = nil
	service.cdpURL = ""
	if service.cmd != nil && service.cmd.Process != nil {
		_ = service.cmd.Process.Kill()
		_, _ = service.cmd.Process.Wait()
		service.cmd = nil
	}
}

func (service *Service) Info() (string, error) {
	if err := service.ensureStarted(); err != nil {
		return "", err
	}
	return service.cdpURL, nil
}

func (service *Service) Navigate(url string) error {
	if err := service.ensureStarted(); err != nil {
		return err
	}
	service.mu.Lock()
	defer service.mu.Unlock()
	return chromedp.Run(service.tabCtx, chromedp.Navigate(url))
}

func (service *Service) Screenshot(path string) error {
	data, err := service.ScreenshotPNG()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (service *Service) ScreenshotPNG() ([]byte, error) {
	if err := service.ensureStarted(); err != nil {
		return nil, err
	}

	service.mu.Lock()
	defer service.mu.Unlock()

	var buf []byte
	if err := chromedp.Run(service.tabCtx, chromedp.CaptureScreenshot(&buf)); err != nil {
		return nil, err
	}
	return buf, nil
}

func (service *Service) Click(x, y float64) error {
	if err := service.ensureStarted(); err != nil {
		return err
	}
	service.mu.Lock()
	defer service.mu.Unlock()
	return chromedp.Run(service.tabCtx, chromedp.MouseClickXY(x, y))
}

func (service *Service) ensureStarted() error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.tabCtx != nil {
		return nil
	}

	if service.config.ExistingWebSocketDebug != "" {
		return service.connectRemote(service.config.ExistingWebSocketDebug)
	}

	return service.launchBrowser()
}

func (service *Service) connectRemote(wsURL string) error {
	allocCtx, allocCancel := chromedp.NewRemoteAllocator(context.Background(), wsURL)
	tabCtx, tabCancel := chromedp.NewContext(allocCtx)
	if err := chromedp.Run(tabCtx); err != nil {
		allocCancel()
		tabCancel()
		return err
	}

	service.allocCtx = allocCtx
	service.allocCancel = allocCancel
	service.tabCtx = tabCtx
	service.tabCancel = tabCancel
	service.cdpURL = wsURL
	return nil
}

func (service *Service) launchBrowser() error {
	binary := service.config.BinaryPath
	if binary == "" {
		binary = detectChromeBinary()
	}
	if binary == "" {
		return ErrBrowserUnavailable
	}

	if service.config.RemoteDebuggingPort == 0 {
		port, err := pickFreePort()
		if err != nil {
			return err
		}
		service.config.RemoteDebuggingPort = port
	}

	userDataDir := service.config.UserDataDir
	if userDataDir == "" {
		userDataDir = filepath.Join(os.TempDir(), "open-sandbox-chrome")
	}
	if err := os.MkdirAll(userDataDir, 0755); err != nil {
		return err
	}

	if err := service.startChromeProcess(binary, userDataDir); err != nil {
		return err
	}

	cdpURL, err := fetchWebSocketURL(service.config.RemoteDebuggingHost, service.config.RemoteDebuggingPort, service.config.StartupTimeout)
	if err != nil {
		service.stopChromeProcess()
		return err
	}

	if err := service.connectRemote(cdpURL); err != nil {
		service.stopChromeProcess()
		return err
	}

	return nil
}

func (service *Service) startChromeProcess(binary string, userDataDir string) error {
	args := []string{
		"--no-first-run",
		"--no-default-browser-check",
		"--disable-background-networking",
		"--disable-client-side-phishing-detection",
		"--disable-component-update",
		"--disable-default-apps",
		"--disable-sync",
		"--disable-translate",
		"--disable-popup-blocking",
		"--remote-debugging-address=" + service.config.RemoteDebuggingHost,
		"--remote-debugging-port=" + fmt.Sprintf("%d", service.config.RemoteDebuggingPort),
		"--user-data-dir=" + userDataDir,
		"--disable-gpu",
		"about:blank",
	}
	if service.config.Headless {
		args = append(args, "--headless=new")
	}

	cmd := exec.Command(binary, args...)
	cmd.Stdout = &bytes.Buffer{}
	cmd.Stderr = cmd.Stdout
	if err := cmd.Start(); err != nil {
		return err
	}
	service.cmd = cmd
	return nil
}

func (service *Service) stopChromeProcess() {
	if service.cmd == nil || service.cmd.Process == nil {
		return
	}
	_ = service.cmd.Process.Kill()
	_, _ = service.cmd.Process.Wait()
	service.cmd = nil
}

func fetchWebSocketURL(host string, port int, timeout time.Duration) (string, error) {
	url := fmt.Sprintf("http://%s:%d/json/version", host, port)
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: minDuration(2*time.Second, timeout)}
	var lastErr error

	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err != nil {
			lastErr = err
			time.Sleep(200 * time.Millisecond)
			continue
		}

		var info versionInfo
		if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
			resp.Body.Close()
			lastErr = err
			time.Sleep(200 * time.Millisecond)
			continue
		}
		resp.Body.Close()

		if info.WebSocketDebuggerURL == "" {
			lastErr = errors.New("missing websocket debugger url")
			time.Sleep(200 * time.Millisecond)
			continue
		}
		return info.WebSocketDebuggerURL, nil
	}
	if lastErr == nil {
		lastErr = errors.New("timed out waiting for browser websocket")
	}
	return "", lastErr
}

func minDuration(a time.Duration, b time.Duration) time.Duration {
	if a <= b {
		return a
	}
	return b
}

func detectChromeBinary() string {
	candidates := []string{
		`C:\Program Files\Google\Chrome\Application\chrome.exe`,
		`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
		`C:\Program Files\Chromium\Application\chrome.exe`,
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return ""
}

func pickFreePort() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}
