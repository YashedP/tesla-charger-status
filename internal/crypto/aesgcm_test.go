package crypto

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestEncryptDecryptRoundTrip(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i + 1)
	}

	cipher, err := NewAESCipher(key)
	if err != nil {
		t.Fatalf("new cipher: %v", err)
	}

	enc, err := cipher.EncryptString("hello-world")
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	got, err := cipher.DecryptString(enc)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if got != "hello-world" {
		t.Fatalf("unexpected plaintext: %q", got)
	}
}

func TestLoadKeyFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token_enc_key.b64")

	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}
	encoded := base64.StdEncoding.EncodeToString(key)
	if err := os.WriteFile(path, []byte(encoded+"\n"), 0o600); err != nil {
		t.Fatalf("write key file: %v", err)
	}

	loaded, err := LoadKeyFromFile(path)
	if err != nil {
		t.Fatalf("load key: %v", err)
	}
	if len(loaded) != 32 {
		t.Fatalf("unexpected key length: %d", len(loaded))
	}
}

func TestLoadKeyFromFileRejectsOpenPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("permission mode checks are skipped on windows")
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "token_enc_key.b64")

	key := make([]byte, 32)
	encoded := base64.StdEncoding.EncodeToString(key)
	if err := os.WriteFile(path, []byte(encoded+"\n"), 0o600); err != nil {
		t.Fatalf("write key file: %v", err)
	}
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatalf("chmod key file: %v", err)
	}

	_, err := LoadKeyFromFile(path)
	if err == nil {
		t.Fatalf("expected permission error, got nil")
	}
	if !strings.Contains(err.Error(), "permissions too open") {
		t.Fatalf("unexpected error: %v", err)
	}
}
