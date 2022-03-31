package router

import (

	// "github.com/celeneisip/atm-proj/api/auth"
	// model "github.com/celeneisip/atm-proj/models"

	"time"

	"github.com/celeneisip/atm-proj/api/auth"
	"github.com/celeneisip/atm-proj/controllers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var dB *gorm.DB

// api router
func New(app *fiber.App, db *gorm.DB) {
	dB = db
	// route maps
	app.Get("/health", health)
	app.Post("/login", login)
	app.Get("account/:type/balance", getBalance)
	app.Put("account/transaction", updateBalance)
}

// health
func health(c *fiber.Ctx) error {
	return c.Status(200).SendString("200 OK")
}

// login
func login(c *fiber.Ctx) error {
	var req controllers.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	res, err := controllers.Login(dB, &req)
	if err != nil {
		return err
	}

	cookie := fiber.Cookie{
		Name:    "jwt",
		Value:   res.Token,
		Expires: time.Unix(res.Exp, 0),
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{"user": res.User, "success": res.Success})

}

func getBalance(c *fiber.Ctx) error {
	claims, err := auth.ValidateJWTToken(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	res, err := controllers.GetBalanceByType(dB, claims.Issuer, c.Params("type"))
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"balance": res, "success": true})

}

func updateBalance(c *fiber.Ctx) error {
	claims, err := auth.ValidateJWTToken(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	var payload *controllers.UpdateBalanceRequest

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	res, err := controllers.UpdateBalanceByType(dB, claims.Issuer, payload)

	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"newBalance": res, "success": true})
}
