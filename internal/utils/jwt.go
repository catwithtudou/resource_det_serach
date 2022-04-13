package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type LoginClaims struct {
	Uid     uint  `json:"uid"`
	LoginTs int64 `json:"login_ts"`
	jwt.StandardClaims
}

var (
	TokenSecret         = []byte("my-name-is-catwithtudou")
	TokenExpireDuration = 48 * time.Hour
)

func GenJwtToken(uid uint) (string, error) {
	ts := time.Now()
	claim := &LoginClaims{
		Uid:     uid,
		LoginTs: ts.Unix(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: ts.Add(TokenExpireDuration).Unix(),
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	result, err := token.SignedString(TokenSecret)
	if err != nil {
		return "", err
	}
	if result == "" {
		return "", errors.New("token is nil")
	}
	return result, nil
}

func ParseJwtToken(tk string) (*LoginClaims, error) {
	token, err := jwt.ParseWithClaims(tk, &LoginClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return TokenSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*LoginClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
