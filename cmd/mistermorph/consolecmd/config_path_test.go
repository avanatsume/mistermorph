package consolecmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestResolveConsoleConfigPath_ExplicitFlagWins(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	restoreConsoleConfigKey(t)
	viper.Set("config", "~/custom.yaml")

	got, err := resolveConsoleConfigPath()
	if err != nil {
		t.Fatalf("resolveConsoleConfigPath() error = %v", err)
	}
	want := filepath.Join(home, "custom.yaml")
	if got != filepath.Clean(want) {
		t.Fatalf("resolveConsoleConfigPath() path = %q, want %q", got, filepath.Clean(want))
	}
}

func TestResolveConsoleConfigPath_DefaultIgnoresCWD(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	restoreConsoleConfigKey(t)
	viper.Set("config", "")

	wd := t.TempDir()
	restoreConsoleWD(t, wd)
	if err := os.WriteFile("config.yaml", []byte("llm:\n  provider: openai\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(config.yaml) error = %v", err)
	}

	got, err := resolveConsoleConfigPath()
	if err != nil {
		t.Fatalf("resolveConsoleConfigPath() error = %v", err)
	}
	want := filepath.Join(home, ".morph", "config.yaml")
	if got != filepath.Clean(want) {
		t.Fatalf("resolveConsoleConfigPath() path = %q, want %q", got, filepath.Clean(want))
	}
}

func restoreConsoleConfigKey(t *testing.T) {
	t.Helper()
	prev, had := viper.Get("config"), viper.IsSet("config")
	t.Cleanup(func() {
		if had {
			viper.Set("config", prev)
			return
		}
		viper.Set("config", nil)
	})
}

func restoreConsoleWD(t *testing.T, wd string) {
	t.Helper()
	prevWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	if err := os.Chdir(wd); err != nil {
		t.Fatalf("Chdir(%q) error = %v", wd, err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(prevWD)
	})
}
