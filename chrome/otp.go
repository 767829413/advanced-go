package chrome

import (
	"time"

	"github.com/pquerna/otp/totp"
)

func GetTotpCode(secret string) (string, error) {
	return totp.GenerateCode(
		secret,
		time.Now(),
	)
}
