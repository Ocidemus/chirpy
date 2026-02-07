package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error){
	key := make([]byte, 32)
	_,err := rand.Read(key)
	if err != nil {
		return "",err
	}
	encodedStr := hex.EncodeToString(key)
	return encodedStr,nil
}