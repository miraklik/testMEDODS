package tokens

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your-secret-key")

type Claims struct {
	UserID string `json:"user_id"`
	IP     string `json:"ip"`
	jwt.StandardClaims
}

func GenerateAccessToken(userID, ip string) (string, error) {
	claims := &Claims{
		UserID: userID,
		IP:     ip,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(jwtSecret)
}

func GenerateRefreshToken() (string, error) {
	token := make([]byte, 64)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	hash := sha512.Sum512(token)
	return base64.URLEncoding.EncodeToString(hash[:]), nil
}

func HashAndSalt(token string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareHash(hashedToken, token string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(token))
}
