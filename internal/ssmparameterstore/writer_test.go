package ssmparameterstore

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func fakeExecSuccess(name string, args ...string) *exec.Cmd {
	return exec.Command("go", "test", "-run", "TestHelperProcess_SSM", "-test.v", "--", "success")
}

func fakeExecFailure(name string, args ...string) *exec.Cmd {
	return exec.Command("go", "test", "-run", "TestHelperProcess_SSM", "-test.v", "--", "failure")
}

var capturedArgs []string

func captureArgs(name string, args ...string) *exec.Cmd {
	capturedArgs = append([]string{name}, args...)
	return exec.Command("true")
}

func TestWriteSecret_Success(t *testing.T) {
	w := NewWriter("/myapp", "")
	w.execCmd = func(name string, args ...string) *exec.Cmd {
		return exec.Command("true")
	}
	if err := w.WriteSecret("prod", "DB_PASS", "secret123"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestWriteSecret_Failure(t *testing.T) {
	w := NewWriter("/myapp", "")
	w.execCmd = func(name string, args ...string) *exec.Cmd {
		return exec.Command("false")
	}
	if err := w.WriteSecret("prod", "DB_PASS", "secret123"); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWriteSecret_PassesCorrectArgs(t *testing.T) {
	capturedArgs = nil
	w := NewWriter("/myapp", "")
	w.execCmd = captureArgs
	_ = w.WriteSecret("staging", "API_KEY", "val")

	full := strings.Join(capturedArgs, " ")
	if !strings.Contains(full, "/myapp/staging/API_KEY") {
		t.Errorf("expected param path in args, got: %s", full)
	}
	if !strings.Contains(full, "SecureString") {
		t.Errorf("expected SecureString type in args, got: %s", full)
	}
	if !strings.Contains(full, "--overwrite") {
		t.Errorf("expected --overwrite flag in args, got: %s", full)
	}
}

func TestWriteSecret_WithKMSKey(t *testing.T) {
	capturedArgs = nil
	w := NewWriter("/secure", "alias/my-key")
	w.execCmd = captureArgs
	_ = w.WriteSecret("ns", "TOKEN", "abc")

	full := strings.Join(capturedArgs, " ")
	if !strings.Contains(full, "alias/my-key") {
		t.Errorf("expected kms key id in args, got: %s", full)
	}
}

func TestWriteSecret_NoKMSKey(t *testing.T) {
	capturedArgs = nil
	w := NewWriter("/app", "")
	w.execCmd = captureArgs
	_ = w.WriteSecret("ns", "KEY", "val")

	full := strings.Join(capturedArgs, " ")
	if strings.Contains(full, "--key-id") {
		t.Errorf("expected no --key-id flag when kmsKeyID is empty, got: %s", full)
	}
}

func TestHelperProcess_SSM(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(0)
}
