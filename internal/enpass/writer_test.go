package enpass

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
		os.Stderr.WriteString("enpass error")
		os.Exit(1)
	}
}

func TestWriteSecret_Success(t *testing.T) {
	w := NewWriter("/path/to/vault")
	w.execCmd = fakeExecSuccess
	if err := w.WriteSecret("myapp", "API_KEY", "s3cr3t"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestWriteSecret_Failure(t *testing.T) {
	w := NewWriter("/path/to/vault")
	w.execCmd = fakeExecFailure
	err := w.WriteSecret("myapp", "API_KEY", "s3cr3t")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWriteSecret_PassesCorrectArgs(t *testing.T) {
	var captured []string
	w := NewWriter("/my/vault")
	w.execCmd = captureArgs(&captured)
	_ = w.WriteSecret("ns", "KEY", "val")

	expected := []string{"enpass-cli", "add", "--vault", "/my/vault", "--category", "Login", "--title", "ns/KEY", "--password", "val"}
	if strings.Join(captured, " ") != strings.Join(expected, " ") {
		t.Errorf("args mismatch:\n got  %v\n want %v", captured, expected)
	}
}
