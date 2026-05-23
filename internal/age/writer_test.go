package age

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func fakeExecSuccess(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess", "--")
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1", "FAKE_EXIT=0")
	return cmd
}

func fakeExecFailure(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess", "--")
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1", "FAKE_EXIT=1")
	return cmd
}

func captureArgs(captured *[]string) func(string, ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		*captured = append([]string{name}, args...)
		return fakeExecSuccess(name, args...)
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	if os.Getenv("FAKE_EXIT") == "1" {
		fmt.Fprintln(os.Stderr, "age: error encrypting")
		os.Exit(1)
	}
	os.Exit(0)
}

func TestWriteSecret_Success(t *testing.T) {
	w := NewWriter("age1public...", "/tmp/secrets.age")
	w.execCmd = fakeExecSuccess

	if err := w.WriteSecret("myapp", "API_KEY", "supersecret"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestWriteSecret_Failure(t *testing.T) {
	w := NewWriter("age1public...", "/tmp/secrets.age")
	w.execCmd = fakeExecFailure

	err := w.WriteSecret("myapp", "API_KEY", "supersecret")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWriteSecret_PassesCorrectArgs(t *testing.T) {
	var captured []string
	w := NewWriter("age1abc123", "/tmp/out.age")
	w.execCmd = captureArgs(&captured)

	_ = w.WriteSecret("ns", "KEY", "val")

	if len(captured) < 4 {
		t.Fatalf("expected at least 4 args, got %d: %v", len(captured), captured)
	}
	if captured[0] != "age" {
		t.Errorf("expected command 'age', got %q", captured[0])
	}
	if captured[1] != "--recipient" || captured[2] != "age1abc123" {
		t.Errorf("expected --recipient age1abc123, got %v", captured[1:3])
	}
	if captured[3] != "--output" || captured[4] != "/tmp/out.age" {
		t.Errorf("expected --output /tmp/out.age, got %v", captured[3:5])
	}
}
