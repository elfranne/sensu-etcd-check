package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sensu/sensu-plugin-sdk/sensu"
)

func resetPlugin() {
	plugin.URL = []string{"http://127.0.0.1:2379"}
	plugin.Size = 1_500_000_000
	plugin.CertFile = ""
	plugin.KeyFile = ""
	plugin.TrustedCAFile = ""
	plugin.Timeout = 5
}

func writeTempFile(t *testing.T, name string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte("dummy"), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestCheckArgs(t *testing.T) {
	existing := writeTempFile(t, "file.pem")
	missing := filepath.Join(t.TempDir(), "missing.pem")

	tests := []struct {
		name      string
		cert      string
		key       string
		ca        string
		wantState int
		wantErr   bool
	}{
		{"no TLS flags", "", "", "", sensu.CheckStateOK, false},
		{"all set and present", existing, existing, existing, sensu.CheckStateOK, false},
		{"only cert set", existing, "", "", sensu.CheckStateUnknown, true},
		{"only key set", "", existing, "", sensu.CheckStateUnknown, true},
		{"cert and key without CA", existing, existing, "", sensu.CheckStateUnknown, true},
		{"cert file missing", missing, existing, existing, sensu.CheckStateUnknown, true},
		{"key file missing", existing, missing, existing, sensu.CheckStateUnknown, true},
		{"CA file missing", existing, existing, missing, sensu.CheckStateUnknown, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetPlugin()
			plugin.CertFile = tt.cert
			plugin.KeyFile = tt.key
			plugin.TrustedCAFile = tt.ca

			state, err := checkArgs(nil)
			if state != tt.wantState {
				t.Errorf("checkArgs() state = %d, want %d", state, tt.wantState)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("checkArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExecuteCheckInvalidTLSFiles(t *testing.T) {
	resetPlugin()
	invalid := writeTempFile(t, "invalid.pem")
	plugin.CertFile = invalid
	plugin.KeyFile = invalid
	plugin.TrustedCAFile = invalid

	state, err := executeCheck(nil)
	if err != nil {
		t.Fatalf("executeCheck() error = %v, want nil", err)
	}
	if state != sensu.CheckStateUnknown {
		t.Errorf("executeCheck() state = %d, want %d", state, sensu.CheckStateUnknown)
	}
}

func TestExecuteCheckUnreachableEndpoint(t *testing.T) {
	resetPlugin()
	plugin.URL = []string{"http://127.0.0.1:1"}
	plugin.Timeout = 1

	state, err := executeCheck(nil)
	if err != nil {
		t.Fatalf("executeCheck() error = %v, want nil", err)
	}
	if state != sensu.CheckStateCritical {
		t.Errorf("executeCheck() state = %d, want %d", state, sensu.CheckStateCritical)
	}
}
