package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// User has many account(eg Joint account/personal/etc), UserID is the foreign key
// May not be relevant to takehome API, but necessary to represent user to account relationship(eg user could have more 1 account like account/personal, etc)
type User struct {
	gorm.Model
	ID         uint      `gorm:"primaryKey;auto_increment"`
	FirstName  string    `gorm:"size:255;not null;" json:"first_name" db:"first_name"`
	LastName   string    `gorm:"size:255;not null;" json:"last_name" db:"last_name"`
	DOB        time.Time `gorm:"not null;" json:"birth_date"`
	Address    string    `gorm:"size:255;not null;" json:"address"`
	City       string    `gorm:"size:80;not null;" json:"city"`
	Province   string    `gorm:"size:80;not null;" json:"province"`
	PostalCode string    `gorm:"size:10;not null;" json:"postalcode"`
	Phone      string    `gorm:"size:15;not null;" json:"phone_number"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Account    []Account `gorm:"foreignKey:UserID"`
}

// Account Model struct
type Account struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey;auto_increment"`
	Number    int       `gorm:"size:80;unique;not null;" json:"account_number"`
	Type      string    `gorm:"size:100; not null;" json:"account_type"`
	Balance   float64   `gorm:"not null;" json:"balance"`
	Active    uint      `gorm:"default:1;not null" json:"active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	UserID    uint
}

// Card Model struct
type Card struct {
	gorm.Model
	Number uint `gorm:"primaryKey;auto_increment"`
	UserID uint
	Pin    string `gorm:"size:100;not null;" json:"pin"`
}

type AccountsDetails struct {
	Card *Card
	User *User
}

// get accountDetails info via card
func GetAccountDetail(db *gorm.DB, cn int) (*AccountsDetails, error) {
	c := Card{}
	u := User{}
	ad := AccountsDetails{}
	invalidId := int64(1)

	carRes := db.Where("number = ?", cn).Find(&c)
	if carRes.RowsAffected < invalidId {
		return &ad, errors.New("can't get user[models]")
	}

	acRes := db.Model(&u).Preload("Account").Find(&u, &c.UserID)
	if acRes.RowsAffected < invalidId {
		return &ad, errors.New("can't get user[models]")
	}

	return &AccountsDetails{
		Card: &c,
		User: &u,
	}, nil
}

// get account balance info via user id and type
func GetAccountBalance(db *gorm.DB, userId string, accountType string) (float64, error) {
	a := Account{}
	invalidId := int64(1)

	acRes := db.Where("user_id = ? AND type = ? AND active = 1", userId, accountType).Find(&a)
	if acRes.RowsAffected < invalidId {
		return 0, errors.New("can't get account balance[models]")
	}
	return a.Balance, nil
}

// get accountDetails info via card
func UpdateAccountBalance(db *gorm.DB, userId string, accountType string, transactionType string, amount float64) (float64, error) {
	if amount < 0 {
		return 0, errors.New("negative amount is an invalid amount for transaction[models]")
	}

	currentAmount, err := GetAccountBalance(db, userId, accountType)
	if err != nil {
		return 0, err
	}

	invalidId := int64(1)
	newAmount := currentAmount
	a := Account{}

	//Assumption only widthraw and deposits are the supported transaction type
	transaction := strings.ToLower(transactionType)
	switch {
	case transaction == "withdrawal":
		newAmount = currentAmount - amount
		//do not allow withdrawal if there is an insufficient funds
		if newAmount < 0 {
			return 0, errors.New("insufficient funds[models]")
		}
	case transaction == "deposit":
		newAmount = currentAmount + amount
	default:
		return 0, errors.New("unknown transaction type[models]")
	}

	acRes := db.Model(&a).Where("user_id = ? AND type = ? AND active = 1", userId, accountType).Update("balance", newAmount)
	if acRes.RowsAffected < invalidId {
		return 0, errors.New("can't update account balance[models]")
	}
	return a.Balance, nil
}
