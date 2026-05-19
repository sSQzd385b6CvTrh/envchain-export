package azurekeyvault

import (
	"fmt"
	"testing"
)

func fakeExecSuccess(_ string, _ ...string) ([]byte, error) {
	return []byte("secret created"), nil
}

func fakeExecFailure(_ string, _ ...string) ([]byte, error) {
	return []byte("ERROR: vault not found"), fmt.Errorf("exit status 1")
}

func captureArgs(captured *[]string) func(string, ...string) ([]byte, error) {
	return func(name string, args ...string) ([]byte, error) {
		*captured = append([]string{name}, args...)
		return []byte{}, nil
	}
}

func TestWriteSecret_Success(t *testing.T) {
	w := NewWriter("my-vault")
	w.execCmd = fakeExecSuccess

	if err := w.WriteSecret("myapp", "DB_PASSWORD", "s3cr3t"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestWriteSecret_Failure(t *testing.T) {
	w := NewWriter("my-vault")
	w.execCmd = fakeExecFailure

	err := w.WriteSecret("myapp", "DB_PASSWORD", "s3cr3t")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWriteSecret_PassesCorrectArgs(t *testing.T) {
	var got []string
	w := NewWriter("prod-vault")
	w.execCmd = captureArgs(&got)

	_ = w.WriteSecret("myapp", "API_KEY", "abc123")

	expectedName := "myapp-API-KEY"
	if len(got) < 8 {
		t.Fatalf("expected at least 8 args, got %d: %v", len(got), got)
	}
	if got[0] != "az" {
		t.Errorf("expected command 'az', got %q", got[0])
	}
	if got[4] != "prod-vault" {
		t.Errorf("expected vault-name 'prod-vault', got %q", got[4])
	}
	if got[6] != expectedName {
		t.Errorf("expected secret name %q, got %q", expectedName, got[6])
	}
	if got[8] != "abc123" {
		t.Errorf("expected value 'abc123', got %q", got[8])
	}
}

func TestSanitizeName(t *testing.T) {
	cases := []struct {
		input, want string
	}{
		{"myapp_DB_HOST", "myapp-DB-HOST"},
		{"my.app.key", "my-app-key"},
		{"already-fine", "already-fine"},
		{"with space", "with-space"},
	}
	for _, tc := range cases {
		got := sanitizeName(tc.input)
		if got != tc.want {
			t.Errorf("sanitizeName(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
