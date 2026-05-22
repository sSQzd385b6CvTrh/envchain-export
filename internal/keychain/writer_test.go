package keychain

import (
	"os"
	"os/exec"
	"testing"
)

func fakeExecSuccess(name string, args ...string) *exec.Cmd {
	return exec.Command(os.Args[0], "-test.run=TestHelperProcess", "--", "success")
}

func fakeExecFailure(name string, args ...string) *exec.Cmd {
	return exec.Command(os.Args[0], "-test.run=TestHelperProcess", "--", "failure")
}

func captureArgs(captured *[]string) func(string, ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		*captured = append([]string{name}, args...)
		return exec.Command(os.Args[0], "-test.run=TestHelperProcess", "--", "success")
	}
}

func TestHelperProcess(t *testing.T) {
	args := os.Args
	for i, a := range args {
		if a == "--" {
			if i+1 < len(args) && args[i+1] == "failure" {
				os.Exit(1)
			}
			os.Exit(0)
		}
	}
}

func TestWriteSecret_Success(t *testing.T) {
	w := &Writer{execCmd: fakeExecSuccess}
	if err := w.WriteSecret("my-app", "API_KEY", "supersecret"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestWriteSecret_Failure(t *testing.T) {
	w := &Writer{execCmd: fakeExecFailure}
	err := w.WriteSecret("my-app", "API_KEY", "supersecret")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWriteSecret_PassesCorrectArgs(t *testing.T) {
	var captured []string
	w := &Writer{execCmd: captureArgs(&captured)}

	_ = w.WriteSecret("svc", "TOKEN", "abc123")

	expected := []string{
		"security",
		"add-generic-password",
		"-s", "svc",
		"-a", "TOKEN",
		"-w", "abc123",
		"-U",
	}
	if len(captured) != len(expected) {
		t.Fatalf("expected %d args, got %d: %v", len(expected), len(captured), captured)
	}
	for i, v := range expected {
		if captured[i] != v {
			t.Errorf("arg[%d]: expected %q, got %q", i, v, captured[i])
		}
	}
}
