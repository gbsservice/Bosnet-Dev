package jwt_auth

import (
	"errors"
	"api_kino/config/app"
	"api_kino/config/constant"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Auth struct {
	UserID   string `json:"user_id,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	ClientId string `json:"client_id,omitempty"`
	Menu     string `json:"menu,omitempty"`
}

type Claim struct {
	Auth
	jwt.StandardClaims
}

type RefreshClaim struct {
	KeyId string `json:"key_id,omitempty"`
	jwt.StandardClaims
}

type AccessTokenMenuClaim struct {
	KeyId string `json:"key_id,omitempty"`
	Menu  string `json:"menu,omitempty"`
	jwt.StandardClaims
}

func CreateToken(auth Auth) (string, string, string) {
	token, expiredAt, _ := GenerateToken(auth)
	refreshToken, _, _ := GenerateRefreshToken(auth.UserID)
	return token, refreshToken, expiredAt
}

func GenerateToken(auth Auth) (string, string, error) {
	var err error
	signature := app.Config().Key
	expiredAt := time.Now().Add(time.Hour * 3)
	//expiredAt := time.Now().Add(time.Second * 5)
	jwtClaim := Claim{
		Auth: auth,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredAt.Unix(),
			Issuer:    auth.UserID,
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaim)
	token, err := at.SignedString([]byte(signature))
	if err != nil {
		return "", "", err
	}
	return token, expiredAt.Format("2006-01-02T15:04:05.00Z"), nil
}

func ValidateToken(tokenString string) (*Claim, error) {
	if tokenString == "" {
		return nil, errors.New(constant.ErrorLogin)
	}
	tokens := strings.Split(tokenString, "Bearer ")
	signature := app.Config().Key
	if len(tokens) < 2 {
		return nil, errors.New(constant.ErrorLogin)
	}
	tokenString = tokens[1]
	token, err := jwt.ParseWithClaims(tokenString, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signature), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claim); ok && token.Valid {
		if claims.UserID == "" && claims.Menu == "" {
			return nil, errors.New(constant.ErrorLogin)
		}
		return claims, nil
	}
	return nil, err
}

func GenerateRefreshToken(keyId string) (string, string, error) {
	var err error
	signature := app.Config().Key
	//expiredAt := time.Now().AddDate(0, 0, 7)
	expiredAt := time.Now().Add(time.Hour * 5)
	jwtClaim := RefreshClaim{
		KeyId: keyId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredAt.Unix(),
			Issuer:    keyId,
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaim)
	token, err := at.SignedString([]byte(signature))
	if err != nil {
		return "", "", err
	}
	return token, expiredAt.Format("2006-01-02T15:04:05.00Z"), nil
}

func ValidateRefreshToken(tokenString string) (*RefreshClaim, error) {
	if tokenString == "" {
		return nil, errors.New(constant.ErrorLogin)
	}
	signature := app.Config().Key
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signature), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*RefreshClaim); ok && token.Valid {
		if claims.KeyId == "" {
			return nil, errors.New(constant.ErrorLogin)
		}
		return claims, nil
	}
	if claims, ok := token.Claims.(*RefreshClaim); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
