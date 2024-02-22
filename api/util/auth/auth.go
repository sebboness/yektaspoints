package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"unicode"

	"github.com/sebboness/yektaspoints/util/log"
)

const (
	GrantTypeClientCredentials = "client_credentials"
	GrantTypePassword          = "password"
	GrantTypeRefreshToken      = "refresh_token"
)

var SupportedGrantTypes = map[string]bool{
	GrantTypePassword:     true,
	GrantTypeRefreshToken: true,
}

type AuthController interface {
	Authenticate(ctx context.Context, username, password string) (AuthResult, error)
	ConfirmRegistration(ctx context.Context, username, code string) error
	RefreshToken(ctx context.Context, username, token string) (AuthResult, error)
	Register(ctx context.Context, ur UserRegisterRequest) (UserRegisterResult, error)
	UpdatePassword(ctx context.Context, session, username, password string) error
}

type AuthResult struct {
	Username            string `json:"username"`
	AccessToken         string `json:"access_token"`
	IdToken             string `json:"id_token"`
	ExpiresIn           int32  `json:"expires_in"`
	NewPasswordRequired bool   `json:"new_password_required"`
	Session             string `json:"session"`
}

type UserRegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

type UserRegisterResult struct {
	IsConfirmed        bool   `json:"is_confirmed"`
	ConfirmationType   string `json:"confirmation_type"`
	ConfirmationSentTo string `json:"confirmation_sent_to"`
	Username           string `json:"username"`
}

var logger = log.Get()

// computeSecretHash returns a secret hash string using HMAC_SHA256 algorithm
func computeSecretHash(username, clientID, clientSecret string) string {
	data := []byte(username + clientID)

	// create a new HMAC by defining the hash type and the key
	hmac := hmac.New(sha256.New, []byte(clientSecret))

	// compute the HMAC
	hmac.Write([]byte(data))
	dataHmac := hmac.Sum(nil)

	encodedHash := base64.StdEncoding.EncodeToString(dataHmac)
	return encodedHash
}

var pwSpecialChars = map[rune]bool{
	'^': true, '$': true, '*': true, '.': true,
	'[': true, ']': true, '{': true, '}': true, '(': true, ')': true, '?': true,
	'"': true, '!': true, '@': true, '#': true, '%': true, '&': true, '/': true,
	'\\': true, ',': true, '>': true, '<': true, '\'': true, ':': true, ';': true,
	'|': true, '_': true, '~': true, '`': true, '=': true, '+': true, '-': true,
}

type pwResult struct {
	WithinLength bool
	Number       bool
	Lower        bool
	Upper        bool
	Special      bool
}

// ValidatePassword validates that a password meets the minimum requirements, which are:
//   - Between 8-256 characters
//   - At least one lower case letter
//   - At least one upper case letter
//   - At least one digit
//   - At least one special character (i.e. one or more of: ^ $ * . [ ] { } ( ) ? " ! @ # % & / \ , > < ' : ; | _ ~ ` = + -
func ValidatePassword(s string) pwResult {
	r := pwResult{}
	letters := 0
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			r.Number = true
		case unicode.IsUpper(c):
			r.Upper = true
		case unicode.IsLower(c):
			r.Lower = true
		default:
			if _, ok := pwSpecialChars[c]; ok {
				r.Special = true
			}
		}

		letters++
	}
	r.WithinLength = letters >= 8 && letters <= 256
	return r
}
