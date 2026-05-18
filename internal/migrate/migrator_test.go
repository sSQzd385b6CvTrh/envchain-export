package migrate

import (
	"errors"
	"testing"
)

// --- fakes ---

type fakeReader struct {
	namespaces []string
	secrets    map[string]map[string]string
	listErr    error
	readErr    error
}

func (f *fakeReader) ListNamespaces() ([]string, error) {
	return f.namespaces, f.listErr
}

func (f *fakeReader) ReadSecrets(ns string) (map[string]string, error) {
	if f.readErr != nil {
		return nil, f.readErr
	}
	return f.secrets[ns], nil
}

type fakeWriter struct {
	written map[string]map[string]string
	writeErr error
}

func (f *fakeWriter) EnsureItem(ns string, secrets map[string]string) error {
	if f.writeErr != nil {
		return f.writeErr
	}
	if f.written == nil {
		f.written = make(map[string]map[string]string)
	}
	f.written[ns] = secrets
	return nil
}

// --- tests ---

func TestMigrateAll_Success(t *testing.T) {
	reader := &fakeReader{
		namespaces: []string{"app1", "app2"},
		secrets: map[string]map[string]string{
			"app1": {"KEY": "val1"},
			"app2": {"TOKEN": "tok"},
		},
	}
	writer := &fakeWriter{}
	m := New(reader, writer)

	results, err := m.MigrateAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Err != nil {
			t.Errorf("namespace %s failed: %v", r.Namespace, r.Err)
		}
	}
	if writer.written["app1"]["KEY"] != "val1" {
		t.Error("expected app1 KEY to be val1")
	}
}

func TestMigrateAll_ListError(t *testing.T) {
	reader := &fakeReader{listErr: errors.New("envchain not found")}
	m := New(reader, &fakeWriter{})
	_, err := m.MigrateAll()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMigrateAll_WriteError(t *testing.T) {
	reader := &fakeReader{
		namespaces: []string{"app1"},
		secrets:    map[string]map[string]string{"app1": {"K": "v"}},
	}
	writer := &fakeWriter{writeErr: errors.New("op error")}
	m := New(reader, writer)

	results, err := m.MigrateAll()
	if err != nil {
		t.Fatalf("unexpected top-level error: %v", err)
	}
	if results[0].Err == nil {
		t.Error("expected per-namespace error, got nil")
	}
}
