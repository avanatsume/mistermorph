//go:build wailsdesktop

package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/pkg/browser"
	"github.com/wailsapp/wails/v3/pkg/application"
)

const desktopOpenURLMessagePrefix = "mistermorph:open-url:"

type App struct {
	wailsApp   *application.App
	restartMu  sync.Mutex
	restarting bool
}

func NewApp() *App {
	return &App{}
}

func (a *App) Attach(wailsApp *application.App) {
	a.wailsApp = wailsApp
}

func (a *App) HandleRawMessage(message string) {
	if !strings.HasPrefix(message, desktopOpenURLMessagePrefix) {
		return
	}

	if err := a.OpenExternalURL(message[len(desktopOpenURLMessagePrefix):]); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "open external URL failed: %v\n", err)
	}
}

func (a *App) OpenExternalURL(rawURL string) error {
	target, err := normalizeExternalBrowserURL(rawURL)
	if err != nil {
		return err
	}
	if err := browser.OpenURL(target); err != nil {
		return fmt.Errorf("open URL in browser: %w", err)
	}
	return nil
}

func normalizeExternalBrowserURL(rawURL string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", fmt.Errorf("empty URL")
	}
	for i, r := range rawURL {
		if r < 32 || r == 127 {
			return "", fmt.Errorf("control character at position %d not allowed", i)
		}
	}
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse URL: %w", err)
	}
	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme != "http" && scheme != "https" {
		return "", fmt.Errorf("unsupported URL scheme %q", parsedURL.Scheme)
	}
	if parsedURL.Host == "" {
		return "", fmt.Errorf("missing URL host")
	}
	return parsedURL.String(), nil
}

// RestartApp relaunches the current executable and quits the current process.
func (a *App) RestartApp() error {
	a.restartMu.Lock()
	if a.restarting {
		a.restartMu.Unlock()
		return nil
	}
	a.restarting = true
	a.restartMu.Unlock()

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve executable path: %w", err)
	}

	cmd := exec.Command(exePath, os.Args[1:]...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if wd, wdErr := os.Getwd(); wdErr == nil {
		cmd.Dir = wd
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start new app process: %w", err)
	}

	if a.wailsApp != nil {
		a.wailsApp.Quit()
	}
	return nil
}
