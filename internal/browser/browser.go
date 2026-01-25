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

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/target"
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
	NavigateTimeout        time.Duration
	ScreenshotTimeout      time.Duration
	Headless               bool
	DownloadDir            string
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

	tabs      []tabHandle
	activeTab int

	downloadsMu sync.Mutex
	downloads   map[string]*DownloadInfo
}

type versionInfo struct {
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
}

type tabHandle struct {
	ctx      context.Context
	cancel   context.CancelFunc
	targetID target.ID
}

type DownloadInfo struct {
	GUID          string    `json:"guid"`
	URL           string    `json:"url"`
	Filename      string    `json:"filename"`
	Path          string    `json:"path"`
	State         string    `json:"state"`
	ReceivedBytes float64   `json:"received_bytes"`
	TotalBytes    float64   `json:"total_bytes"`
	StartedAt     time.Time `json:"started_at"`
	FinishedAt    time.Time `json:"finished_at"`
}

func DefaultConfig() Config {
	return Config{
		RemoteDebuggingHost: "127.0.0.1",
		RemoteDebuggingPort: 9222,
		StartupTimeout:      15 * time.Second,
		NavigateTimeout:     15 * time.Second,
		ScreenshotTimeout:   15 * time.Second,
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
	if config.NavigateTimeout == 0 {
		config.NavigateTimeout = 15 * time.Second
	}
	if config.ScreenshotTimeout == 0 {
		config.ScreenshotTimeout = 15 * time.Second
	}
	return &Service{
		config:    config,
		downloads: make(map[string]*DownloadInfo),
	}
}

func (service *Service) Start() error {
	return service.ensureStarted()
}

func (service *Service) Close() {
	service.mu.Lock()
	defer service.mu.Unlock()

	service.resetLocked()
}

func (service *Service) Info() (string, error) {
	if err := service.ensureStarted(); err != nil {
		return "", err
	}
	return service.cdpURL, nil
}

func (service *Service) Navigate(url string) error {
	return service.runTabAction(service.config.NavigateTimeout, func(ctx context.Context) error {
		return chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
			_, _, _, err := page.Navigate(url).Do(ctx)
			return err
		}))
	})
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
	var buf []byte
	if err := service.runTabAction(service.config.ScreenshotTimeout, func(ctx context.Context) error {
		return chromedp.Run(ctx, chromedp.CaptureScreenshot(&buf))
	}); err != nil {
		return nil, err
	}
	return buf, nil
}

func (service *Service) Click(x, y float64) error {
	return service.runTabAction(service.config.NavigateTimeout, func(ctx context.Context) error {
		return chromedp.Run(ctx, chromedp.MouseClickXY(x, y))
	})
}

func (service *Service) FormInputFill(selector string, value string) error {
	return service.runTabAction(service.config.NavigateTimeout, func(ctx context.Context) error {
		return chromedp.Run(ctx,
			chromedp.Focus(selector, chromedp.ByQuery),
			chromedp.SetValue(selector, value, chromedp.ByQuery),
		)
	})
}

func (service *Service) ElementSelect(selector string, value string) error {
	return service.runTabAction(service.config.NavigateTimeout, func(ctx context.Context) error {
		return chromedp.Run(ctx, chromedp.SetValue(selector, value, chromedp.ByQuery))
	})
}

func (service *Service) Scroll(x float64, y float64) error {
	expr := fmt.Sprintf("window.scrollBy(%f, %f)", x, y)
	return service.runTabAction(service.config.NavigateTimeout, func(ctx context.Context) error {
		return chromedp.Run(ctx, chromedp.Evaluate(expr, nil))
	})
}

func (service *Service) Evaluate(expression string) (any, error) {
	var result any
	if err := service.runTabAction(service.config.NavigateTimeout, func(ctx context.Context) error {
		return chromedp.Run(ctx, chromedp.Evaluate(expression, &result))
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func (service *Service) PressKey(keys string) error {
	return service.runTabAction(service.config.NavigateTimeout, func(ctx context.Context) error {
		return chromedp.Run(ctx, chromedp.KeyEvent(keys))
	})
}

func (service *Service) NewTab(url string) (int, error) {
	service.mu.Lock()
	defer service.mu.Unlock()

	if err := service.ensureStartedLocked(); err != nil {
		return -1, err
	}

	handle, err := service.createTabLocked()
	if err != nil {
		return -1, err
	}
	service.tabs = append(service.tabs, handle)
	service.activeTab = len(service.tabs) - 1
	service.tabCtx = handle.ctx
	service.tabCancel = handle.cancel

	if url != "" {
		if err := runWithTimeout(handle.ctx, service.config.NavigateTimeout, func(ctx context.Context) error {
			_, _, _, err := page.Navigate(url).Do(ctx)
			return err
		}); err != nil {
			return -1, err
		}
	}
	return service.activeTab, nil
}

func (service *Service) SwitchTab(index int) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if err := service.ensureStartedLocked(); err != nil {
		return err
	}
	if index < 0 || index >= len(service.tabs) {
		return errors.New("invalid tab index")
	}
	handle := service.tabs[index]
	if handle.ctx == nil || handle.ctx.Err() != nil {
		return errors.New("tab unavailable")
	}
	service.activeTab = index
	service.tabCtx = handle.ctx
	service.tabCancel = handle.cancel
	return nil
}

func (service *Service) CloseTab(index int) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if err := service.ensureStartedLocked(); err != nil {
		return err
	}
	if index < 0 || index >= len(service.tabs) {
		return errors.New("invalid tab index")
	}
	handle := service.tabs[index]
	if handle.ctx != nil {
		_ = target.CloseTarget(handle.targetID).Do(handle.ctx)
		handle.cancel()
	}
	service.tabs = append(service.tabs[:index], service.tabs[index+1:]...)
	if len(service.tabs) == 0 {
		newHandle, err := service.createTabLocked()
		if err != nil {
			return err
		}
		service.tabs = append(service.tabs, newHandle)
	}
	if service.activeTab >= len(service.tabs) {
		service.activeTab = len(service.tabs) - 1
	}
	active := service.tabs[service.activeTab]
	service.tabCtx = active.ctx
	service.tabCancel = active.cancel
	return nil
}

type TabInfo struct {
	Index int    `json:"index"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

func (service *Service) TabList() ([]TabInfo, error) {
	service.mu.Lock()
	defer service.mu.Unlock()

	if err := service.ensureStartedLocked(); err != nil {
		return nil, err
	}

	results := make([]TabInfo, 0, len(service.tabs))
	for i, tab := range service.tabs {
		if tab.ctx == nil || tab.ctx.Err() != nil {
			results = append(results, TabInfo{Index: i})
			continue
		}
		var title string
		var url string
		_ = runWithTimeout(tab.ctx, service.config.NavigateTimeout, func(ctx context.Context) error {
			return chromedp.Run(ctx,
				chromedp.Evaluate("document.title", &title),
				chromedp.Evaluate("location.href", &url),
			)
		})
		results = append(results, TabInfo{Index: i, Title: title, URL: url})
	}
	return results, nil
}

func (service *Service) DownloadList() []DownloadInfo {
	service.downloadsMu.Lock()
	defer service.downloadsMu.Unlock()

	results := make([]DownloadInfo, 0, len(service.downloads))
	for _, info := range service.downloads {
		results = append(results, *info)
	}
	return results
}

func (service *Service) ensureStarted() error {
	service.mu.Lock()
	defer service.mu.Unlock()

	return service.ensureStartedLocked()
}

func (service *Service) ensureStartedLocked() error {
	if service.isTabHealthyLocked() {
		return nil
	}

	service.resetLocked()

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
	service.cdpURL = wsURL
	service.downloads = make(map[string]*DownloadInfo)

	handle, err := service.setupTabLocked(tabCtx, tabCancel)
	if err != nil {
		tabCancel()
		allocCancel()
		return err
	}
	service.tabs = []tabHandle{handle}
	service.activeTab = 0
	service.tabCtx = handle.ctx
	service.tabCancel = handle.cancel
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

func (service *Service) runTabAction(timeout time.Duration, action func(ctx context.Context) error) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	if err := service.ensureStartedLocked(); err != nil {
		return err
	}

	err := runWithTimeout(service.tabCtx, timeout, action)
	if err == nil {
		return nil
	}
	if !isContextErr(err) {
		return err
	}

	service.resetLocked()
	if err := service.ensureStartedLocked(); err != nil {
		return err
	}
	return runWithTimeout(service.tabCtx, timeout, action)
}

func runWithTimeout(parent context.Context, timeout time.Duration, action func(ctx context.Context) error) error {
	if timeout <= 0 {
		return action(parent)
	}
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()
	return action(ctx)
}

func isContextErr(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

func (service *Service) isTabHealthyLocked() bool {
	if service.tabCtx == nil {
		return false
	}
	if err := service.tabCtx.Err(); err != nil {
		return false
	}
	if service.cmd != nil && service.cmd.ProcessState != nil && service.cmd.ProcessState.Exited() {
		return false
	}
	return true
}

func (service *Service) resetLocked() {
	for _, tab := range service.tabs {
		if tab.cancel != nil {
			tab.cancel()
		}
	}
	service.tabs = nil
	service.activeTab = 0
	service.tabCtx = nil
	service.tabCancel = nil
	if service.allocCancel != nil {
		service.allocCancel()
		service.allocCancel = nil
	}
	service.allocCtx = nil
	service.cdpURL = ""
	service.stopChromeProcess()
	service.downloads = make(map[string]*DownloadInfo)
}

func (service *Service) ensureDownloadDirLocked() (string, error) {
	if service.config.DownloadDir == "" {
		service.config.DownloadDir = filepath.Join(os.TempDir(), "open-sandbox-downloads")
	}
	if err := os.MkdirAll(service.config.DownloadDir, 0755); err != nil {
		return "", err
	}
	return service.config.DownloadDir, nil
}

func (service *Service) setupTabLocked(ctx context.Context, cancel context.CancelFunc) (tabHandle, error) {
	downloadDir, err := service.ensureDownloadDirLocked()
	if err != nil {
		return tabHandle{}, err
	}
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch e := ev.(type) {
		case *browser.EventDownloadWillBegin:
			filename := e.SuggestedFilename
			if filename == "" {
				filename = e.GUID
			}
			info := &DownloadInfo{
				GUID:      e.GUID,
				URL:       e.URL,
				Filename:  filename,
				Path:      filepath.Join(downloadDir, filename),
				State:     browser.DownloadProgressStateInProgress.String(),
				StartedAt: time.Now(),
			}
			service.downloadsMu.Lock()
			service.downloads[e.GUID] = info
			service.downloadsMu.Unlock()
		case *browser.EventDownloadProgress:
			service.downloadsMu.Lock()
			info := service.downloads[e.GUID]
			if info == nil {
				info = &DownloadInfo{GUID: e.GUID, StartedAt: time.Now()}
				service.downloads[e.GUID] = info
			}
			info.TotalBytes = e.TotalBytes
			info.ReceivedBytes = e.ReceivedBytes
			info.State = e.State.String()
			if e.State == browser.DownloadProgressStateCompleted || e.State == browser.DownloadProgressStateCanceled {
				info.FinishedAt = time.Now()
			}
			service.downloadsMu.Unlock()
		}
	})

	if err := chromedp.Run(ctx,
		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllow).
			WithDownloadPath(downloadDir).
			WithEventsEnabled(true),
	); err != nil {
		return tabHandle{}, err
	}

	targetID := target.ID("")
	if c := chromedp.FromContext(ctx); c != nil && c.Target != nil {
		targetID = c.Target.TargetID
	}
	return tabHandle{ctx: ctx, cancel: cancel, targetID: targetID}, nil
}

func (service *Service) createTabLocked() (tabHandle, error) {
	tabCtx, tabCancel := chromedp.NewContext(service.allocCtx)
	if err := chromedp.Run(tabCtx); err != nil {
		tabCancel()
		return tabHandle{}, err
	}
	return service.setupTabLocked(tabCtx, tabCancel)
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
