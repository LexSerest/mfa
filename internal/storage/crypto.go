package storage

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
	"io"

	"mfa/internal/models"
)

func getKey(password string, salt []byte) ([]byte, error) {
	key, err := pbkdf2.Key(sha256.New, password, salt, 10000, 32)
	if err != nil {
		return nil, fmt.Errorf("failed: %w", err)
	}
	return key, nil
}

func encryptAccounts(data []byte, password string, salt []byte) ([]byte, error) {
	key, err := getKey(password, salt)
	if err != nil {
		return []byte{}, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, fmt.Errorf("failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, fmt.Errorf("failed: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return []byte{}, fmt.Errorf("failed: %w", err)
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func decryptAccounts(data []byte, password string, salt []byte) ([]byte, error) {
	key, err := getKey(password, salt)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, actualCipherText := data[:nonceSize], data[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, actualCipherText, nil)
	if err != nil {
		return nil, errors.New("invalid password")
	}
	return plainText, nil
}

func generateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)

	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("failed: %w", err)
	}

	return salt, nil
}

func encode(acc models.AccountsDecode, password string, salt []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(acc)
	if err != nil {
		return []byte{}, fmt.Errorf("failed: %w", err)
	}

	return encryptAccounts(buf.Bytes(), password, salt)
}

func decode(data []byte, password string, salt []byte) (models.AccountsDecode, error) {
	if len(data) == 0 || len(salt) == 0 {
		return make(models.AccountsDecode), nil
	}
	gb, err := decryptAccounts(data, password, salt)
	if err != nil {
		return models.AccountsDecode{}, err
	}

	var accounts models.AccountsDecode
	reader := bytes.NewReader(gb)

	err = gob.NewDecoder(reader).Decode(&accounts)
	if err != nil {
		return models.AccountsDecode{}, fmt.Errorf("failed: %w", err)
	}

	return accounts, nil
}
