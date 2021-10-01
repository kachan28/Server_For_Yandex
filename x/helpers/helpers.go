package helpers

import (
	"discountDealer/conf"
	"discountDealer/models/userModels"
	app_context "discountDealer/x/app-context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

func ParseToken(c *fiber.Ctx, tokenMethod string) error {
	auth := c.Get("Authorization")
	l := len(conf.Config.AuthScheme)
	if len(auth) < l+1 || !strings.EqualFold(auth[:l], conf.Config.AuthScheme) {
		return fmt.Errorf("Missing or malformed JWT")
	}
	auth = auth[l+1:]
	token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Missing or malformed JWT")
		}
		return conf.Config.JWTKey, nil
	})
	if err != nil {
		return err
	}

	tokenClaims := token.Claims.(jwt.MapClaims)
	if tokenClaims[userModels.MethodClaim] != tokenMethod {
		return fmt.Errorf("Token not for %v", tokenMethod)
	}

	if actual := tokenClaims.VerifyExpiresAt(time.Now().Unix(), true); actual == false {
		return fmt.Errorf("Please, login again")
	}

	app_context.SetUserID(c, tokenClaims[userModels.IdClaim].(string))

	return nil
}
