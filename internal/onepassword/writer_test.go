package onepassword

import (
	"errors"
	"strings"
	"testing"
)

func fakeExecSuccess(_ string, _ ...string) ([]byte, error) {
	return []byte("ok"), nil
}

func fakeExecAlreadyExists(_ string, _ ...string) ([]byte, error) {
	return []byte("already exists"), errors.New("exit status 1")
}

func fakeExecFailure(_ string, _ ...string) ([]byte, error) {
	return []byte("unexpected error"), errors.New("exit status 2")
}

func TestEnsureItem_Success(t *testing.T) {
	w := &Writer{vault: "TestVault", execCmd: fakeExecSuccess}
	err := w.EnsureItem("myapp", map[string]string{"API_KEY": "secret"})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestEnsureItem_AlreadyExists(t *testing.T) {
	// createItem returns already-exists; setField succeeds.
	callCount := 0
	w := &Writer{
		vault: "TestVault",
		execCmd: func(name string, args ...string) ([]byte, error) {
			callCount++
			if callCount == 1 {
				return fakeExecAlreadyExists(name, args...)
			}
			return fakeExecSuccess(name, args...)
		},
	}
	err := w.EnsureItem("myapp", map[string]string{"TOKEN": "abc"})
	if err != nil {
		t.Fatalf("expected no error on already-exists, got: %v", err)
	}
}

func TestEnsureItem_CreateFailure(t *testing.T) {
	w := &Writer{vault: "TestVault", execCmd: fakeExecFailure}
	err := w.EnsureItem("myapp", map[string]string{"KEY": "val"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "creating item") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestNewWriter(t *testing.T) {
	w := NewWriter("MyVault")
	if w.vault != "MyVault" {
		t.Errorf("expected vault MyVault, got %s", w.vault)
	}
	if w.execCmd == nil {
		t.Error("expected execCmd to be set")
	}
}
