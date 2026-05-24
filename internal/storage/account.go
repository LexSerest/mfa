package storage

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"mfa/internal/models"
	"mfa/internal/totp"
)

func AccountList() []string {
	s := loadStore()
	return s.Accounts
}

func AccountGet(label, password string) (models.Account, error) {
	s := loadStore()
	accounts, err := decode(s.Data, password, s.Salt)
	if err != nil {
		return models.Account{}, err
	}

	account, ok := accounts[label]
	if !ok {
		return models.Account{}, fmt.Errorf("'%s' does not exist", label)
	}

	return account, nil
}

func AccountAdd(label, password string, account models.Account) error {
	s := loadStore()
	accounts, err := decode(s.Data, password, s.Salt)
	if err != nil {
		return err
	}

	_, ok := accounts[label]
	if ok {
		return fmt.Errorf("'%s' label exists", label)
	}

	if !totp.ValidSecret(account.Secret) {
		return fmt.Errorf("2fa code not valid")
	}

	accounts[label] = account
	return saveStore(s, accounts, password)
}

func AccountRename(label, newLabel, password string) error {
	s := loadStore()
	accounts, err := decode(s.Data, password, s.Salt)
	if err != nil {
		return err
	}

	account, ok := accounts[label]
	if !ok {
		return fmt.Errorf("'%s' does not exist", label)
	}

	if _, exists := accounts[newLabel]; exists {
		return fmt.Errorf("'%s' already exists", newLabel)
	}

	accounts[newLabel] = account
	delete(accounts, label)
	return saveStore(s, accounts, password)
}

func AccountDelete(label, password string) error {
	s := loadStore()
	accounts, err := decode(s.Data, password, s.Salt)
	if err != nil {
		return err
	}

	_, ok := accounts[label]
	if !ok {
		return fmt.Errorf("'%s' does not exist", label)
	}

	delete(accounts, label)
	return saveStore(s, accounts, password)
}

func AccountGenCode(label, password string) (string, string, int, error) {
	account, err := AccountGet(label, password)

	if err != nil {
		return "", "", 0, err
	}

	code, nextCode, timeLeft := totp.GenerateCode(account)
	return code, nextCode, timeLeft, nil
}

func AccountGenURI(label, password string) (string, error) {
	account, err := AccountGet(label, password)

	if err != nil {
		return "", err
	}

	return totp.UriGen(label, account) + "\n", nil
}

func AccountImport(path, password string) ([]string, error) {
	list := []string{}
	s := loadStore()
	accounts, err := decode(s.Data, password, s.Salt)
	if err != nil {
		return list, err
	}

	file, err := os.Open(path)
	if err != nil {
		return list, fmt.Errorf("Error open file %v\n", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "otpauth://") {
			continue
		}

		u, err := url.Parse(line)
		if err != nil {
			continue
		}

		label := strings.TrimPrefix(u.Path, "/")
		if parts := strings.Split(label, ":"); len(parts) > 1 {
			oldLabel := label
			label = strings.TrimSpace(parts[0])
			fmt.Fprintf(os.Stderr, "Warning: label '%s' was shortened to '%s'\n", oldLabel, label)
		}

		if _, exists := accounts[label]; exists {
			fmt.Fprintf(os.Stderr, "Skip: account '%s' already exists. Rename it manually before import.\n", label)
			continue
		}

		q := u.Query()
		secret := q.Get("secret")
		if secret == "" {
			continue
		}

		digits, _ := strconv.Atoi(q.Get("digits"))
		period, _ := strconv.Atoi(q.Get("period"))
		algo := q.Get("algorithm")

		accounts[label] = models.NewAccount(secret, digits, period, algo)
		list = append(list, label)
	}

	if err := scanner.Err(); err != nil {
		return list, err
	}

	saveStore(s, accounts, password)
	return list, nil
}

func AccountExport(path, password string) error {
	s := loadStore()
	accounts, err := decode(s.Data, password, s.Salt)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open file for export: %v", err)
	}
	defer file.Close()

	for label, account := range accounts {
		uri := totp.UriGen(label, account)
		file.WriteString(uri + "\n")
	}

	return nil
}
