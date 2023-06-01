package encrypt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var CustomSecret = []byte("k8s-test")

var (
	ErrToKenInvalid = errors.New("token验证失败")
	ErrTokenExpire  = errors.New("token过期")
)

const TokenExpireDuration = time.Hour * 24

type MyClaim struct {
	Username string
	Id       uint32
	Uid      string
	jwt.RegisteredClaims
}

func GenToken(username, uid string, id uint32) (string, error) {
	myclaim := MyClaim{
		username,
		id,
		uid,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    "k8s-test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, myclaim)
	return token.SignedString(CustomSecret)
}

func VerifyToken(tokenString string) (*MyClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaim{}, func(token *jwt.Token) (interface{}, error) {
		return CustomSecret, nil
	})
	if err != nil {
		fmt.Println(err)
		return nil, ErrToKenInvalid
	}
	if myclaim, ok := token.Claims.(*MyClaim); ok && token.Valid {
		return myclaim, nil
	}
	return nil, ErrTokenExpire
}
