package helper

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserClaimsData struct {
	ID       string `json:"nik"`
	Username string `json:"name"`
	RoleID   string `json:"role_id"`
}

type JWTClaims struct {
	jwt.StandardClaims
	User UserClaimsData `json:"user"`
}

var (
	secretKey       = os.Getenv("JWT_SECRET_KEY")
	sessionDuration = os.Getenv("JWT_SESSION_DURATION")
)

func GenerateJWT(issuer string, user UserClaimsData) (string, error) {
	session, _ := strconv.Atoi(sessionDuration)
	claims := &JWTClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			Issuer: issuer,
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(session)).Unix(),
		},
	}

	tokens := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	return tokens.SignedString([]byte(secretKey))
}

func ParseJWT(tokenString string) (claims JWTClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return claims, err
	}

	return claims, err
}