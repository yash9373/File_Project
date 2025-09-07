package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/scrypt"
)

const (
	saltSize  = 16 // 16 bytes salt
	nonceSize = 12 // 12 bytes nonce for GCM
	keyLen    = 32 // 32 bytes key for AES-256
)

// DeriveKey derives a key from password+salt using scrypt
func DeriveKey(password string, salt []byte) ([]byte, error) {
	return scrypt.Key([]byte(password), salt, 32768, 8, 1, keyLen)
}

// EncryptBytes encrypts plaintext using AES-256-GCM with scrypt-derived key.
// Output format: [salt(16)][nonce(12)][ciphertext]
func EncryptBytes(plaintext []byte, password string) ([]byte, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	key, err := DeriveKey(password, salt)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	out := make([]byte, 0, len(salt)+len(nonce)+len(ciphertext))
	out = append(out, salt...)
	out = append(out, nonce...)
	out = append(out, ciphertext...)
	return out, nil
}

// DecryptBytes expects input format [salt(16)][nonce(12)][ciphertext]
func DecryptBytes(input []byte, password string) ([]byte, error) {
	if len(input) < saltSize+nonceSize+1 {
		return nil, errors.New("ciphertext too short")
	}
	salt := input[:saltSize]
	nonce := input[saltSize : saltSize+nonceSize]
	ciphertext := input[saltSize+nonceSize:]
	key, err := DeriveKey(password, salt)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return gcm.Open(nil, nonce, ciphertext, nil)
}
