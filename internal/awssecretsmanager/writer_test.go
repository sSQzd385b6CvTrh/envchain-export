package awssecretsmanager

import (
	"fmt"
	"strings"
	"testing"
)

func fakeExecSuccess(_ string, _ ...string) ([]byte, error) {
	return []byte(`{"ARN":"arn:aws:secretsmanager:us-east-1:123:secret:ns/key"}`), nil
}

func fakeExecAlreadyExists(_ string, _ ...string) ([]byte, error) {
	return []byte("ResourceExistsException: secret already exists"), fmt.Errorf("exit status 254")
}

func fakeExecFailure(_ string, _ ...string) ([]byte, error) {
	return []byte("An error occurred"), fmt.Errorf("exit status 1")
}

func captureArgs(captured *[]string) func(string, ...string) ([]byte, error) {
	return func(name string, args ...string) ([]byte, error) {
		*captured = append([]string{name}, args...)
		return []byte(`{}`), nil
	}
}

func TestWriteSecret_Success(t *testing.T) {
	w := NewWriter("us-east-1")
	w.execCmd = fakeExecSuccess

	if err := w.WriteSecret("myapp", "API_KEY", "s3cr3t"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestWriteSecret_AlreadyExists_UpdatesSecret(t *testing.T) {
	callCount := 0
	w := NewWriter("")
	w.execCmd = func(name string, args ...string) ([]byte, error) {
		callCount++
		if callCount == 1 {
			return fakeExecAlreadyExists(name, args...)
		}
		return fakeExecSuccess(name, args...)
	}

	if err := w.WriteSecret("myapp", "DB_PASS", "hunter2"); err != nil {
		t.Fatalf("expected no error on update path, got %v", err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 exec calls (create + update), got %d", callCount)
	}
}

func TestWriteSecret_Failure(t *testing.T) {
	w := NewWriter("")
	w.execCmd = fakeExecFailure

	if err := w.WriteSecret("myapp", "KEY", "val"); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWriteSecret_PassesRegion(t *testing.T) {
	var captured []string
	w := NewWriter("eu-west-1")
	w.execCmd = captureArgs(&captured)

	_ = w.WriteSecret("svc", "TOKEN", "abc")

	args := strings.Join(captured, " ")
	if !strings.Contains(args, "--region eu-west-1") {
		t.Errorf("expected --region eu-west-1 in args, got: %s", args)
	}
}

func TestWriteSecret_NoRegion_OmitsFlag(t *testing.T) {
	var captured []string
	w := NewWriter("")
	w.execCmd = captureArgs(&captured)

	_ = w.WriteSecret("svc", "TOKEN", "abc")

	args := strings.Join(captured, " ")
	if strings.Contains(args, "--region") {
		t.Errorf("expected no --region flag when region is empty, got: %s", args)
	}
}
