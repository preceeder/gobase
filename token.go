package gobase

import (
	"encoding/base64"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/golang-jwt/jwt/v5"
	"github.com/preceeder/gobase/utils"
	"time"
)

const TokenSignKey = "tjly(Ap]@m$T0^x"
const TokenAesKey = "tjly(Ap]@m$T0^23"

type CustomClaims struct {
	jwt.RegisteredClaims
	UserData map[string]string
}

func TokenGenerateUsingHs256(mclaim map[string]string) (string, error) {
	claim := CustomClaims{
		UserData: mclaim,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "token",                                    // 签发者
			Subject:  "token",                                    // 签发对象
			Audience: jwt.ClaimStrings{"Android_APP", "IOS_APP"}, //签发受众
			//ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)), //过期时间
			//NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)), //最早使用时间
			IssuedAt: jwt.NewNumericDate(time.Now()), //签发时间
			ID:       utils.RandStr(10),              // wt ID, 类似于盐值
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(TokenSignKey))
	token = base64.StdEncoding.EncodeToString(cryptor.AesCbcEncrypt([]byte(token), []byte(TokenAesKey)))
	return token, err
}

func TokenParseHs256(tokenSecrete string) (*CustomClaims, error) {
	defer func() {
		if err := recover(); err != nil {
			panic(utils.BaseHttpError{ErrorCode: 401, Message: "toekn error"})
		}
	}()
	tokenbyte, _ := base64.StdEncoding.DecodeString(tokenSecrete)
	tokenString := string(cryptor.AesCbcDecrypt(tokenbyte, []byte(TokenAesKey)))
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(TokenSignKey), nil //返回签名密钥
	})

	if err != nil {
		return nil, utils.BaseHttpError{ErrorCode: 401, Message: "toekn error 1"}
	}

	if !token.Valid {
		return nil, utils.BaseHttpError{ErrorCode: 401, Message: "toekn error 2"}
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, utils.BaseHttpError{ErrorCode: 401, Message: "toekn error 3"}
	}
	return claims, nil
}
