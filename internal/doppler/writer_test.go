package doppler

import (
	"errors"
	"strings"
	"testing"
)

func fakeExecSuccess(_ string, _ ...string) ([]byte, error) {
	return []byte("Updated secrets."), nil
}

func fakeExecFailure(_ string, _ ...string) ([]byte, error) {
	return []byte("error: invalid credentials"), errors.New("exit status 1")
}

func captureArgs() (func(string, ...string) ([]byte, error), *[]string) {
	var captured []string
	fn := func(_ string, args ...string) ([]byte, error) {
		captured = append(captured, args...)
		return []byte("Updated secrets."), nil
	}
	return fn, &captured
}

func TestWriteSecret_Success(t *testing.T) {
	w := NewWriter("myproject", "dev")
	w.execCmd = fakeExecSuccess

	if err := w.WriteSecret("API_KEY", "supersecret"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestWriteSecret_PassesCorrectArgs(t *testing.T) {
	w := NewWriter("myproject", "dev")
	fn, captured := captureArgs()
	w.execCmd = fn

	_ = w.WriteSecret("DB_PASS", "hunter2")

	args := strings.Join(*captured, " ")
	for _, want := range []string{"--project", "myproject", "--config", "dev", "DB_PASS=hunter2"} {
		if !strings.Contains(args, want) {
			t.Errorf("expected args to contain %q, got: %s", want, args)
		}
	}
}

func TestWriteSecret_Failure(t *testing.T) {
	w := NewWriter("myproject", "dev")
	w.execCmd = fakeExecFailure

	err := w.WriteSecret("API_KEY", "val")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "API_KEY") {
		t.Errorf("expected error to mention key name, got: %v", err)
	}
}

func TestWriteSecrets_PropagatesError(t *testing.T) {
	w := NewWriter("proj", "prd")
	w.execCmd = fakeExecFailure

	secrets := map[string]string{"FOO": "bar"}
	err := w.WriteSecrets("mynamespace", secrets)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "mynamespace") {
		t.Errorf("expected error to mention namespace, got: %v", err)
	}
}
