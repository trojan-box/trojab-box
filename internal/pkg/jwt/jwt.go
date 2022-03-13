package jwt

import (
	"errors"
	"github.com/aresprotocols/trojan-box/internal/app"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type ClientClaims struct {
	Address string `json:"address"`
	jwt.StandardClaims
}

func GenToken(address string, signSecret []byte) (string, error) {

	c := ClientClaims{
		address,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(app.Conf.JwtTimeout) * time.Second).Unix(),
			Issuer:    "trojan-box",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	return token.SignedString(signSecret)
}

func ParseToken(tokenString string, signSecret []byte) (*ClientClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &ClientClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return signSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*ClientClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
