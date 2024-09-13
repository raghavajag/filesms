package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type Maker interface {
	CreateToken(userID uuid.UUID, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}

type JWTMaker struct {
	secretKey string
}

type Payload struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}

func NewJWTMaker(secretKey string) Maker {
	return &JWTMaker{secretKey: secretKey}
}

func (maker *JWTMaker) CreateToken(userID uuid.UUID, duration time.Duration) (string, error) {
	payload := &Payload{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		return nil, err
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return payload, nil
}
