package app_context

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

const (
	userIDContextKey = "user-id-key"
)

func SetUserID(c *fiber.Ctx, userID string){
	c.SetUserContext(context.WithValue(c.UserContext(), userIDContextKey, userID))
}

func ExtractUserID(c *fiber.Ctx) string{
	return c.UserContext().Value(userIDContextKey).(string)
}