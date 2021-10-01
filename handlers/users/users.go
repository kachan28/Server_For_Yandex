package users

import (
	"discountDealer/logger"
	"discountDealer/models/userModels"
	"discountDealer/repository/userRepository"
	"discountDealer/validator"
	app_context "discountDealer/x/app-context"
	"discountDealer/x/helpers"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func RegisterHandler(c *fiber.Ctx) error {
	log := logger.New("User Register Handler")
	newUser := new(userModels.User)

	log.Info("body - ", zap.Any("register", c.Body()))
	err := c.BodyParser(newUser)
	if err != nil {
		log.Error("Can't decode body", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := validator.Validate(newUser)
	if errors != nil {
		log.Info("Errors during validation", zap.Any("errors", errors))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errors,
		})
	}

	newUser.GenerateData()

	err = userRepository.Init().Insert(newUser)
	if err != nil {
		log.Error("Can't create user", zap.Any("user", newUser), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	accessToken, refreshToken, err := newUser.GenerateTokens()
	if err != nil {
		log.Error("Can't create user", zap.Any("user", newUser), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func LoginHandler(c *fiber.Ctx) error {
	log := logger.New("User Login Handler")
	loginUser := new(userModels.User)

	err := c.BodyParser(loginUser)
	if err != nil {
		log.Error("Can't decode body", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if !validator.UserExist(loginUser.Username) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Пользователя с таким именем не существует",
		})
	}

	passwordToHash := loginUser.Password
	err = userRepository.Init().GetByFilter("username = ?", loginUser.Username, loginUser)
	if err != nil {
		log.Error("Can't get user", zap.Error(err))
	}

	if loginUser.Password != loginUser.HashPassword(passwordToHash) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Неверный пароль",
		})
	}

	accessToken, refreshToken, err := loginUser.GenerateTokens()
	if err != nil {
		log.Error("Can't generate tokens", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func UpdateTokens(c *fiber.Ctx) error{
	log := logger.New("Update tokens")

	err := helpers.ParseToken(c, userModels.TokenRefresh)
	if err != nil{
		log.Error("Error while updating tokens", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message":err.Error(),
		})
	}

	log.Info("updating tokens")

	id := app_context.ExtractUserID(c)
	user := &userModels.User{ID: &id}

	accessToken, refreshToken, err := user.GenerateTokens()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}