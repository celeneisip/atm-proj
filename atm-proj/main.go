package main

import (
	"fmt"
	"log"
	"os"

	model "github.com/celeneisip/atm-proj/models"
	"github.com/celeneisip/atm-proj/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// loads .env values
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print(".env file found")
	}
}

// initialize db
func InitializeDB() {
	var err error

	dns := fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local`, os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	db, err = gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to Database")
	}
	var card model.Card
	var user model.User
	var account model.Account
	db.AutoMigrate(&card, &user, &account)
}

// main
func main() {
	InitializeDB()

	app := fiber.New()
	app.Use(cors.New(cors.Config{AllowCredentials: true}))

	router.New(app, db)
	app.Listen(":3001")
}
