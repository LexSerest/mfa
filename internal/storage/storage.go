package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"

	"mfa/internal/models"
)

func GetDbPath() string {
	if envPath := os.Getenv("MFA_DB_PATH"); envPath != "" {
		if fi, err := os.Stat(envPath); err == nil {
			if fi.IsDir() {
				fmt.Fprintf(os.Stderr, "Error: MFA_DB_PATH points to a directory, it must be a file path: %s\n", envPath)
				os.Exit(1)
			}
			return envPath
		}
		dir := filepath.Dir(envPath)
		if err := os.MkdirAll(dir, 0700); err != nil {
			fmt.Fprintf(os.Stderr, "Error: cannot create directory: %v\n", err)
			os.Exit(1)
		}
		return envPath
	}

	baseConfigDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot find user config directory: %v\n", err)
		os.Exit(1)
	}
	appConfigDir := filepath.Join(baseConfigDir, "mfa")

	if err := os.MkdirAll(appConfigDir, 0700); err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot create storage directory: %v\n", err)
		os.Exit(1)
	}

	return filepath.Join(appConfigDir, "storage.db")
}

func loadStore() models.Storage {
	s := models.Storage{Accounts: []string{}}
	data, err := os.ReadFile(GetDbPath())

	if err != nil {
		salt, err := generateSalt(32)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		s.Salt = salt
		return s
	}

	reader := bytes.NewReader(data)
	err = gob.NewDecoder(reader).Decode(&s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading database: %v\n", err)
		os.Exit(1)
	}

	return s
}

func saveStore(s models.Storage, accounts models.AccountsDecode, password string) error {
	data, err := encode(accounts, password, s.Salt)
	if err != nil {
		return err
	}

	s.Data = data
	s.Accounts = slices.Collect(maps.Keys(accounts))
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(s); err != nil {
		return err
	}

	if err := os.WriteFile(GetDbPath(), buf.Bytes(), 0600); err != nil {
		return err
	}

	return nil
}
