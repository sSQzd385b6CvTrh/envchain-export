package gopass

import (
	"fmt"
	"strings"
	"testing"
)

func fakeExecSuccess(_ string, _ ...string) ([]byte, error) {
	return []byte("ok"), nil
}

func fakeExecFailure(_ string, _ ...string) ([]byte, error) {
	return []byte("error: something went wrong"), fmt.Errorf("exit status 1")
}

func captureArgs(captured *[]string) func(string, ...string) ([]byte, error) {
	return func(name string, args ...string) ([]byte, error) {
		*captured = append([]string{name}, args...)
		return []byte("ok"), nil
	}
}

func TestWriteSecret_Success(t *testing.T) {
	w := NewWriter("")
	w.execCmd = fakeExecSuccess

	if err := w.WriteSecret("myapp", "API_KEY", "s3cr3t"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestWriteSecret_Failure(t *testing.T) {
	w := NewWriter("")
	w.execCmd = fakeExecFailure

	err := w.WriteSecret("myapp", "API_KEY", "s3cr3t")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "gopass insert failed") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestWriteSecret_PathWithoutStore(t *testing.T) {
	w := NewWriter("")
	var captured []string
	w.execCmd = captureArgs(&captured)

	_ = w.WriteSecret("myapp", "DB_PASS", "hunter2")

	expectedPath := "myapp/DB_PASS"
	if len(captured) < 4 || captured[3] != expectedPath {
		t.Errorf("expected path %q, got args: %v", expectedPath, captured)
	}
}

func TestWriteSecret_PathWithStore(t *testing.T) {
	w := NewWriter("work")
	var captured []string
	w.execCmd = captureArgs(&captured)

	_ = w.WriteSecret("myapp", "DB_PASS", "hunter2")

	expectedPath := "work/myapp/DB_PASS"
	if len(captured) < 4 || captured[3] != expectedPath {
		t.Errorf("expected path %q, got args: %v", expectedPath, captured)
	}
}

func TestBuildPath_NoStore(t *testing.T) {
	w := NewWriter("")
	got := w.buildPath("ns", "KEY")
	if got != "ns/KEY" {
		t.Errorf("expected 'ns/KEY', got %q", got)
	}
}

func TestBuildPath_WithStore(t *testing.T) {
	w := NewWriter("personal")
	got := w.buildPath("ns", "KEY")
	if got != "personal/ns/KEY" {
		t.Errorf("expected 'personal/ns/KEY', got %q", got)
	}
}
