package core

import (
	"ShopService/models"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var accessKey = []byte(JWTConfig["ACCESS_KEY"])
var refreshKey = []byte(JWTConfig["REFRESH_KEY"])

type JWTClaim struct {
	User models.User
	jwt.StandardClaims
}

func GenerateJWT(user models.User, refreshToken bool) (tokenString string, err error) {
	var expirationTime time.Time
	if refreshToken {
		expirationTime = time.Now().Add(48 * time.Hour)
	} else {
		expirationTime = time.Now().Add(1 * time.Hour)
	}

	claims := &JWTClaim{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if refreshToken {
		tokenString, err = token.SignedString(refreshKey)
	} else {
		tokenString, err = token.SignedString(accessKey)
	}

	return
}

func ValidateAccessToken(signedToken string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return accessKey, nil
		},
	)
	if err != nil {
		err = errors.New("token is invalid")
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return nil, err
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("access token is expired")
		return nil, err
	}
	return claims, nil
}

func ValidateRefreshToken(signedToken string) (error, *JWTClaim) {
	var err error
	var claims *JWTClaim
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return refreshKey, nil
		},
	)
	if err != nil {
		err = errors.New("token is invalid")

	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("refresh token is expired. Please Login")
		return err, nil
	}
	return nil, claims
}
