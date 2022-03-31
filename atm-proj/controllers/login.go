package controllers

import (
	"strconv"

	"github.com/celeneisip/atm-proj/api/auth"
	model "github.com/celeneisip/atm-proj/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequest struct {
	CardNumber int `json:"card_number"`
	Pin        int `json:"pin"`
}

type LoginResponse struct {
	Token   string
	Exp     int64
	User    *model.User
	Success bool
}

// verify pin. parse pin to string, since cannot convert pin (type int) to type []byte
func verifyPin(hashedPassword string, pin int) error {
	//convert pin to string to convert to byte
	parsedPin := strconv.Itoa(pin)

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(parsedPin))
}

func Login(db *gorm.DB, req *LoginRequest) (*LoginResponse, error) {
	var res *LoginResponse

	ad, err := model.GetAccountDetail(db, req.CardNumber)
	if err != nil {
		return res, err
	}

	// verify if pin matched system record
	err = verifyPin(ad.Card.Pin, req.Pin)
	if err != nil {
		return res, err
	}

	token, exp, err := auth.CreateJWTToken(ad.User)

	if err != nil {
		return res, err
	}

	res = &LoginResponse{
		Token:   token,
		Exp:     exp,
		User:    ad.User,
		Success: true,
	}
	return res, nil

}
