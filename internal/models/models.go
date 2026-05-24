package models

type Account struct {
	Secret    string
	Digits    int
	Period    int
	Algorithm string
}

type AccountsDecode map[string]Account

type Storage struct {
	Data     []byte
	Accounts []string
	Salt     []byte
}

func NewAccount(secret string, digits int, period int, algorithm string) Account {
	if digits == 0 {
		digits = 6
	}
	if period == 0 {
		period = 30
	}
	if algorithm == "" {
		algorithm = "SHA1"
	}
	return Account{
		Secret:    secret,
		Digits:    digits,
		Period:    period,
		Algorithm: algorithm,
	}
}
