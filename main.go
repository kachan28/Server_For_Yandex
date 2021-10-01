package main

import (
	"context"
	"discountDealer/handlers/products"
	"discountDealer/handlers/users"
	"discountDealer/logger"
	"discountDealer/middlewares"
	"discountDealer/x"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func main() {
	log := logger.New("Discont Dealer")
	ctx := logger.InsertLogger(context.Background(), logger.MainLoggerKey, log)

	app := fiber.New()

	app.Post("/registration", users.RegisterHandler)
	app.Post("/login", users.LoginHandler)
	app.Post("/tokens", users.UpdateTokens)
	app.Get(fmt.Sprintf("/product/:%v", x.ArticulParameter), products.GetProduct)
	app.Get("products/find", products.Find)

	app.Use(middlewares.JWTMiddleware)

	logger.ExtractLogger(ctx, logger.MainLoggerKey).Info("Server has been started")

	app.Listen(":9000")
}
