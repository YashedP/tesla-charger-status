package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

type AESCipher struct {
	aead cipher.AEAD
}

func NewAESCipher(rawKey []byte) (*AESCipher, error) {
	if len(rawKey) != 32 {
		return nil, fmt.Errorf("encryption key must be 32 bytes, got %d", len(rawKey))
	}

	block, err := aes.NewCipher(rawKey)
	if err != nil {
		return nil, fmt.Errorf("create aes cipher: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create gcm cipher: %w", err)
	}

	return &AESCipher{aead: aead}, nil
}

func LoadKeyFromFile(path string) ([]byte, error) {
	if err := ensurePrivateKeyFilePermissions(path); err != nil {
		return nil, err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read key file: %w", err)
	}

	trimmed := strings.TrimSpace(string(content))
	if trimmed == "" {
		return nil, errors.New("key file is empty")
	}

	raw, err := base64.StdEncoding.DecodeString(trimmed)
	if err != nil {
		return nil, fmt.Errorf("decode key file base64: %w", err)
	}
	if len(raw) != 32 {
		return nil, fmt.Errorf("decoded key length must be 32 bytes, got %d", len(raw))
	}

	return raw, nil
}

func ensurePrivateKeyFilePermissions(path string) error {
	if runtime.GOOS == "windows" {
		// POSIX-style mode bits are not reliable on Windows.
		return nil
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat key file: %w", err)
	}

	// Reject keys readable/writable/executable by group or others.
	if info.Mode().Perm()&0o077 != 0 {
		return fmt.Errorf("key file permissions too open (%#o); expected owner-only permissions like 0o600", info.Mode().Perm())
	}

	return nil
}

func (c *AESCipher) EncryptString(plaintext string) (string, error) {
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generate nonce: %w", err)
	}

	ciphertext := c.aead.Seal(nil, nonce, []byte(plaintext), nil)
	blob := append(nonce, ciphertext...)

	return base64.StdEncoding.EncodeToString(blob), nil
}

func (c *AESCipher) DecryptString(encoded string) (string, error) {
	blob, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("decode ciphertext base64: %w", err)
	}

	nonceSize := c.aead.NonceSize()
	if len(blob) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := blob[:nonceSize], blob[nonceSize:]
	plaintext, err := c.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt ciphertext: %w", err)
	}

	return string(plaintext), nil
}
