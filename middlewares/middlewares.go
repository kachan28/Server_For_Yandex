package middlewares

import (
	"discountDealer/logger"
	"discountDealer/models/userModels"
	"discountDealer/x/helpers"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const TokenContextKey = "token-key"

func JWTMiddleware(c *fiber.Ctx) error {
	log := logger.New("JWT-Middleware")

	err := helpers.ParseToken(c, userModels.TokenAccess)
	if err != nil {
		log.Error("Got error", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	log.Info("Good token")

	return c.Next()
}
