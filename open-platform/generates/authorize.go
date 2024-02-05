package generates

import (
	"bytes"
	"context"
	"encoding/base64"
	"strings"

	"github.com/google/uuid"
)

// NewAuthorizeGenerate create to generate the authorize code instance
func NewAuthorizeGenerateIns() AuthorizeGenerate {
	return &AuthorizeGenerateIns{}
}

// AuthorizeGenerate generate the authorize code
type AuthorizeGenerateIns struct{}

// Token based on the UUID generated token
func (ag *AuthorizeGenerateIns) Token(
	ctx context.Context,
	data *GenerateBasic,
) (string, error) {
	buf := bytes.NewBufferString(data.Client.GetID())
	buf.WriteString(data.UserID)
	token := uuid.NewMD5(uuid.Must(uuid.NewRandom()), buf.Bytes())
	code := base64.URLEncoding.EncodeToString([]byte(token.String()))
	code = strings.ToUpper(strings.TrimRight(code, "="))

	return code, nil
}
