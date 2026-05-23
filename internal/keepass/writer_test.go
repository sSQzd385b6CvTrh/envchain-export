package keepass

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// fakeExecSuccess returns a command that exits 0.
func fakeExecSuccess(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess", "--", "success")
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	return cmd
}

// fakeExecAlreadyExists returns a command that prints "already exists" and exits 1.
func fakeExecAlreadyExists(callCount *int) func(string, ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		*callCount++
		if *callCount == 1 {
			cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess", "--", "already_exists")
			cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
			return cmd
		}
		// Second call (edit) succeeds.
		cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess", "--", "success")
		cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
		return cmd
	}
}

// fakeExecFailure returns a command that exits non-zero with a generic error.
func fakeExecFailure(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess", "--", "failure")
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	return cmd
}

// captureArgs records the arguments passed to the fake command.
func captureArgs(captured *[]string) func(string, ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		*captured = append([]string{name}, args...)
		cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess", "--", "success")
		cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
		return cmd
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	args := os.Args
	for i, a := range args {
		if a == "--" {
			args = args[i+1:]
			break
		}
	}
	switch args[0] {
	case "success":
		os.Exit(0)
	case "already_exists":
		os.Stderr.WriteString("Entry with path already exists")
		os.Exit(1)
	case "failure":
		os.Stderr.WriteString("unexpected error")
		os.Exit(2)
	}
}

func TestWriteSecret_Success(t *testing.T) {
	w := NewWriter("/tmp/test.kdbx", "mygroup")
	w.execCmd = fakeExecSuccess

	if err := w.WriteSecret("myns", "API_KEY", "s3cr3t"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestWriteSecret_AlreadyExists_Edits(t *testing.T) {
	count := 0
	w := NewWriter("/tmp/test.kdbx", "")
	w.execCmd = fakeExecAlreadyExists(&count)

	if err := w.WriteSecret("myns", "TOKEN", "abc"); err != nil {
		t.Fatalf("expected no error on already-exists path, got: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 exec calls (add + edit), got %d", count)
	}
}

func TestWriteSecret_Failure(t *testing.T) {
	w := NewWriter("/tmp/test.kdbx", "")
	w.execCmd = fakeExecFailure

	if err := w.WriteSecret("myns", "BAD", "val"); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWriteSecret_PassesCorrectArgs(t *testing.T) {
	var captured []string
	w := NewWriter("/tmp/db.kdbx", "secrets")
	w.execCmd = captureArgs(&captured)

	_ = w.WriteSecret("myapp", "DB_PASS", "hunter2")

	if len(captured) == 0 {
		t.Fatal("no args captured")
	}
	if captured[0] != "keepassxc-cli" {
		t.Errorf("expected keepassxc-cli, got %s", captured[0])
	}
	entryArg := captured[len(captured)-1]
	if !strings.Contains(entryArg, "secrets/myapp/DB_PASS") {
		t.Errorf("expected entry path to contain 'secrets/myapp/DB_PASS', got %s", entryArg)
	}
}
