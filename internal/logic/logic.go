package logic

import (
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"

	"summer/practice2022/internal/config"
	st "summer/practice2022/internal/structures"
)

var (
	ErrRefreshTokenIsAlreadyUsed = errors.New("refresh token is already used")
	ErrPasswordIsIncorrect       = errors.New("password is incorrect")
)

type DB interface {
	GetHashedPassword(login string) (st.HashedPassword, error)
	SaveRefreshTokens(refresh, login string) error
	GetLoginByRefreshToken(refresh string) (string, bool, error)
	LogAlreadyUsedRefreshToken(refresh string) error
	MarkTokenAsUsed(refresh string) error
}

type Logic struct {
	cfg config.Config
	db  DB
}

func NewLogic(cfg config.Config, db DB) *Logic {
	return &Logic{
		cfg: cfg,
		db:  db,
	}
}

func (l *Logic) GetTokensByLoginAndPassword(
	login string,
	password string,
) (
	st.Tokens,
	error,
) {
	hashedPassword, err := l.db.GetHashedPassword(login)
	if err != nil {
		err = errors.Wrap(err, "get hashed password by login")
		return st.Tokens{}, err
	}

	if !hashedPassword.Compare(password) {
		err := ErrPasswordIsIncorrect
		return st.Tokens{}, err
	}

	return l.generateTokens(login)
}

func (l *Logic) GetTokensByRefreshToken(refresh string) (st.Tokens, error) {
	login, isUsed, err := l.db.GetLoginByRefreshToken(refresh)
	if err != nil {
		err = errors.Wrap(err, "get tokens by refresh token")
		return st.Tokens{}, err
	}

	if isUsed {
		if err := l.db.LogAlreadyUsedRefreshToken(refresh); err != nil {
			err = errors.Wrap(err, "log already used refresh token")
			return st.Tokens{}, err
		}

		return st.Tokens{}, ErrRefreshTokenIsAlreadyUsed
	}

	return l.generateTokens(login)
}

func (l *Logic) generateTokens(login string) (st.Tokens, error) {
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(l.cfg.AccessTokenExpiration)
	claims := jwt.RegisteredClaims{
		Issuer:    "auth-server",
		Subject:   login,
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokens := st.Tokens{}

	var err error

	tokens.Access, err = token.SignedString(l.cfg.JWTSecret)
	if err != nil {
		err = errors.Wrap(err, "sign access token")
		return st.Tokens{}, err
	}

	tokens.Refresh = randRefreshToken(l.cfg.RefreshTokenLength)

	if err := l.db.SaveRefreshTokens(tokens.Refresh, login); err != nil {
		err = errors.Wrap(err, "save refresh token")
		return st.Tokens{}, err
	}

	return tokens, nil
}

var letterRunes = []rune(
	"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
)

func randRefreshToken(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
