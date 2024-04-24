package helper

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = "947df3a2c6aee40d3bef66a24be193dec84fcb0c233fd3a5e899f9300a3e792ac0bea819f39f33cabfed2cae097b0b4ed5acd0009a39b6bdc9beb8e21f19409baadd9757af465dc2e9dfdff2d8d12097f7bd4034ec5396ad9cf3a2ce82c4ec436322d0cdf3a69734432842524a0da78d7a1205134edda560c4ee9020cdc61e52"

func GenerateToken(id uint, email string) string {
	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
		"exp":   time.Now().Add(time.Minute * 10),
	}

	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := parseToken.SignedString([]byte(secretKey))
	if err != nil {
		return err.Error()
	}

	return signedToken

}

func VerifyToken(ctx *gin.Context) (interface{}, error) {
	headerToken := ctx.Request.Header.Get("Authorization")
	bearer := strings.HasPrefix(headerToken, "Bearer")

	if !bearer {
		return nil, errors.New("sign in to proceed")
	}

	stringToken := strings.Split(headerToken, " ")[1]
	token, _ := jwt.Parse(stringToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("sign in to proceed")
		}
		return []byte(secretKey), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, errors.New("sign in to proceed")
	}

	expClaim, exists := claims["exp"]
	if !exists {
		return nil, errors.New("expire claim is missing")
	}

	expStr, ok := expClaim.(string)
	if !ok {
		return nil, errors.New("expire claim is not a valid type")
	}

	expTime, err := time.Parse(time.RFC3339, expStr)
	if err != nil {
		return nil, errors.New("error parsing expiration time")
	}

	if time.Now().After(expTime) {
		return nil, errors.New("token is expired")
	}

	return token.Claims.(jwt.MapClaims), nil

}
