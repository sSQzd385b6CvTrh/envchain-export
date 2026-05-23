package conjur

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func fakeExecSuccess(name string, args ...string) *exec.Cmd {
	return exec.Command(os.Args[0], append([]string{"-test.run=TestHelperProcess", "--", "success"}, args...)...)
}

func fakeExecFailure(name string, args ...string) *exec.Cmd {
	return exec.Command(os.Args[0], append([]string{"-test.run=TestHelperProcess", "--", "failure"}, args...)...)
}

func captureArgs(captured *[]string) func(string, ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		*captured = append([]string{name}, args...)
		return fakeExecSuccess(name, args...)
	}
}

func TestHelperProcess(t *testing.T) {
	args := os.Args
	for i, a := range args {
		if a == "--" {
			args = args[i+1:]
			break
		}
	}
	if len(args) == 0 {
		return
	}
	switch args[0] {
	case "success":
		os.Exit(0)
	case "failure":
		os.Stderr.WriteString("error: permission denied")
		os.Exit(1)
	}
}

func TestWriteSecret_Success(t *testing.T) {
	w := NewWriter("myorg")
	w.execCmd = fakeExecSuccess

	if err := w.WriteSecret("myapp", "DB_PASSWORD", "s3cr3t"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestWriteSecret_Failure(t *testing.T) {
	w := NewWriter("myorg")
	w.execCmd = fakeExecFailure

	err := w.WriteSecret("myapp", "DB_PASSWORD", "s3cr3t")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "conjur: failed to set variable") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestWriteSecret_PassesCorrectArgs(t *testing.T) {
	var captured []string
	w := NewWriter("acme")
	w.execCmd = captureArgs(&captured)

	_ = w.WriteSecret("ci", "API_KEY", "abc123")

	if len(captured) == 0 {
		t.Fatal("no command was captured")
	}
	if captured[0] != "conjur" {
		t.Errorf("expected binary 'conjur', got %q", captured[0])
	}
	wantID := "acme:variable:ci/API_KEY"
	found := false
	for _, a := range captured {
		if a == wantID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected variable ID %q in args %v", wantID, captured)
	}
}
