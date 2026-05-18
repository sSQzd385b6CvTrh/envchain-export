package envchain

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// buildFakeEnvchain compiles a small fake envchain binary used in tests.
func buildFakeEnvchain(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "envchain")

	src := `package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) == 1 && args[0] == "--list" {
		fmt.Println("myapp")
		fmt.Println("database")
		return
	}
	if len(args) == 2 && args[1] == "env" {
		switch args[0] {
		case "myapp":
			fmt.Println("API_KEY=secret123")
			fmt.Println("APP_TOKEN=tok_abc")
		case "database":
			fmt.Println("DB_PASSWORD=hunter2")
		}
	}
}
`
	srcFile := filepath.Join(tmpDir, "fake_envchain.go")
	if err := os.WriteFile(srcFile, []byte(src), 0600); err != nil {
		t.Fatalf("writing fake envchain source: %v", err)
	}
	cmd := exec.Command("go", "build", "-o", binPath, srcFile)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("building fake envchain: %v\n%s", err, out)
	}
	return binPath
}

func TestListNamespaces(t *testing.T) {
	bin := buildFakeEnvchain(t)
	r := NewReader(bin)

	ns, err := r.ListNamespaces()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ns) != 2 {
		t.Fatalf("expected 2 namespaces, got %d", len(ns))
	}
	if ns[0] != "myapp" || ns[1] != "database" {
		t.Errorf("unexpected namespaces: %v", ns)
	}
}

func TestReadSecrets(t *testing.T) {
	bin := buildFakeEnvchain(t)
	r := NewReader(bin)

	secrets, err := r.ReadSecrets("myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(secrets) != 2 {
		t.Fatalf("expected 2 secrets, got %d", len(secrets))
	}
	expected := []Secret{
		{Namespace: "myapp", Key: "API_KEY", Value: "secret123"},
		{Namespace: "myapp", Key: "APP_TOKEN", Value: "tok_abc"},
	}
	for i, s := range secrets {
		if s != expected[i] {
			t.Errorf("secret[%d]: got %+v, want %+v", i, s, expected[i])
		}
	}
}
