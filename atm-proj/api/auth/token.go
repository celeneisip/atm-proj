package auth

import (
	"errors"
	"os"
	"strconv"
	"time"

	model "github.com/celeneisip/atm-proj/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

// create new jwt token
func CreateJWTToken(user *model.User) (string, int64, error) {
	//Token expires after 10 minutes
	exp := time.Now().Add(time.Minute * 10).Unix()
	claims := &jwt.StandardClaims{
		ExpiresAt: exp,
		Issuer:    strconv.Itoa(int(user.ID)),
	}
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := jwt.SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		return "", 0, err
	}
	return t, exp, nil
}

func ValidateJWTToken(c *fiber.Ctx) (*jwt.StandardClaims, error) {
	cookie := c.Cookies("jwt")
	if cookie == "" {
		// Means user has not yet login
		return &jwt.StandardClaims{}, errors.New("unauthenticated[auth]")
	}
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("API_SECRET")), nil
	})

	if err != nil {
		// Means user has not yet login
		return &jwt.StandardClaims{}, errors.New("unauthenticated[auth]")
	}
	claims := token.Claims.(*jwt.StandardClaims)
	return claims, nil
}
