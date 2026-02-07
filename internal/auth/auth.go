package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	hash,err := argon2id.CreateHash(password,argon2id.DefaultParams)
	if err != nil {
		return "",err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error){
	check,err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false,err
	} 
	return check,nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error){
	claims := jwt.RegisteredClaims{
		
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn).UTC()),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		Issuer:    "chirpy",
		Subject:   userID.String(),
	}
	signingkey := []byte(tokenSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	tokenString, err := token.SignedString(signingkey)
	if err != nil {
		return "",err
	}
	return tokenString,nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error){

	claims := &jwt.RegisteredClaims{}

	token,err := jwt.ParseWithClaims(tokenString,claims,func(token *jwt.Token) (interface{}, error){
		tokenbyte := []byte(tokenSecret)
		return tokenbyte,nil
	})
	if err != nil {
    	return uuid.UUID{}, err
	}
	if !token.Valid {
		return uuid.UUID{}, errors.New("token is invalid")
	}
	sub,err := uuid.Parse(claims.Subject)
	if err != nil{
		return uuid.UUID{}, errors.New("cant parse thorugh uuid")
	}
	return sub,nil
}

func GetBearerToken(headers http.Header) (string, error){
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "",errors.New("header is empty")
	}
	parts := strings.Fields(authHeader)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "",errors.New("Invalid authorization header")
	}
	token := parts[1]
	return token,nil
}