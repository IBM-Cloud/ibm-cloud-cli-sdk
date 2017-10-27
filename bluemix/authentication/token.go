package authentication

import (
	"fmt"
	"strings"
)

type Token struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
}

func (t Token) Token() string {
	return strings.TrimSpace(fmt.Sprintf("%s %s", t.TokenType, t.AccessToken))
}
