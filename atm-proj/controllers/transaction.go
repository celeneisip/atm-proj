package controllers

import (
	model "github.com/celeneisip/atm-proj/models"
	"gorm.io/gorm"
)

type UpdateBalanceRequest struct {
	AccountType     string  `json:"account_type"`
	TransactionType string  `json:"transaction_type"`
	Amount          float64 `json:"amount"`
}

func GetBalanceByType(db *gorm.DB, userId string, accountType string) (float64, error) {
	return model.GetAccountBalance(db, userId, accountType)
}

func UpdateBalanceByType(db *gorm.DB, userID string, ubpr *UpdateBalanceRequest) (float64, error) {
	return model.UpdateAccountBalance(db, userID, ubpr.AccountType, ubpr.TransactionType, ubpr.Amount)
}
