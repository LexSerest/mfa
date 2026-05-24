package totp

import (
	"encoding/base32"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"mfa/internal/models"
)

func parseAlgo(a string) otp.Algorithm {
	switch strings.ToUpper(a) {
	case "SHA256":
		return otp.AlgorithmSHA256
	case "SHA512":
		return otp.AlgorithmSHA512
	default:
		return otp.AlgorithmSHA1
	}
}

func genCode(now time.Time, account models.Account) (string, error) {
	return totp.GenerateCodeCustom(account.Secret, now, totp.ValidateOpts{
		Period:    uint(account.Period),
		Digits:    otp.Digits(account.Digits),
		Algorithm: parseAlgo(account.Algorithm),
	})
}

func GenerateCode(account models.Account) (string, string, int) {
	now := time.Now()
	nextTime := now.Add(time.Duration(account.Period) * time.Second)
	currentCode, _ := genCode(now, account)
	nextCode, _ := genCode(nextTime, account)
	timeLeft := int64(account.Period) - (now.Unix() % int64(account.Period))
	return currentCode, nextCode, int(timeLeft)
}

func ValidSecret(secret string) bool {
	cleanSecret := strings.TrimSpace(secret)
	cleanSecret = strings.ToUpper(cleanSecret)
	dbuf := make([]byte, base32.StdEncoding.DecodedLen(len(cleanSecret)))
	_, err := base32.StdEncoding.Decode(dbuf, []byte(cleanSecret))
	return err == nil
}

func UriGen(label string, acc models.Account) string {
	return fmt.Sprintf("otpauth://totp/%s?secret=%s&issuer=mfa&digits=%d&period=%d&algorithm=%s",
		url.PathEscape(label), acc.Secret, acc.Digits, acc.Period, acc.Algorithm)
}
