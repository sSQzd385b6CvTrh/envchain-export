package lastpass

import (
	"errors"
	"strings"
	"testing"
)

func fakeExecSuccess(_ string, _ ...string) ([]byte, error) {
	return []byte("Success"), nil
}

func fakeExecFailure(_ string, _ ...string) ([]byte, error) {
	return []byte("Error: Could not find specified account(s)"), errors.New("exit status 1")
}

func captureArgs(captured *[]string) func(string, ...string) ([]byte, error) {
	return func(name string, args ...string) ([]byte, error) {
		*captured = append([]string{name}, args...)
		return []byte("Success"), nil
	}
}

func TestWriteSecret_Success(t *testing.T) {
	w := NewWriter("envchain")
	w.execCmd = fakeExecSuccess

	if err := w.WriteSecret("myapp", "API_KEY", "supersecret"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestWriteSecret_Failure(t *testing.T) {
	w := NewWriter("envchain")
	w.execCmd = fakeExecFailure

	err := w.WriteSecret("myapp", "API_KEY", "supersecret")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "lpass add failed") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestWriteSecret_PassesCorrectArgs(t *testing.T) {
	var captured []string
	w := NewWriter("envchain")
	w.execCmd = captureArgs(&captured)

	_ = w.WriteSecret("myapp", "DB_PASS", "s3cr3t")

	if len(captured) == 0 {
		t.Fatal("no command captured")
	}
	if captured[0] != "lpass" {
		t.Errorf("expected lpass, got %s", captured[0])
	}
	foundName := false
	for _, arg := range captured {
		if strings.Contains(arg, "envchain/myapp/DB_PASS") {
			foundName = true
		}
	}
	if !foundName {
		t.Errorf("expected item name to contain 'envchain/myapp/DB_PASS', args: %v", captured)
	}
}

func TestWriteSecret_NoGroup(t *testing.T) {
	var captured []string
	w := NewWriter("")
	w.execCmd = captureArgs(&captured)

	_ = w.WriteSecret("myapp", "TOKEN", "abc123")

	for _, arg := range captured {
		if strings.Contains(arg, "myapp/TOKEN") {
			return
		}
	}
	t.Errorf("expected item name 'myapp/TOKEN' without group prefix, args: %v", captured)
}
