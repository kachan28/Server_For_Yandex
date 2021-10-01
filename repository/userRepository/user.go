package userRepository

import (
	"context"
	"discountDealer/conf"
	"discountDealer/logger"
	"discountDealer/models/userModels"
	"github.com/go-pg/pg/v10"
	"go.uber.org/zap"
)

var UserRepositoryContext = context.Background()

type userDB struct {
	db *pg.DB
}

func Init() *userDB {
	UserRepositoryContext = logger.InsertLogger(
		UserRepositoryContext,
		logger.UserRepositoryKey,
		logger.New("User Repository"),
	)
	return &userDB{}
}

func (u *userDB) Get(user *userModels.User) error {
	u.connect()
	defer u.disconnect()
	log := logger.ExtractLogger(
		UserRepositoryContext,
		logger.UserRepositoryKey,
	)
	log.Info("Getting user with ID", zap.Any("id", *user.ID))
	err := u.db.Model(user).WherePK().Select()
	if err != nil {
		return err
	}
	return nil
}

func (u *userDB) GetByFilter(filter string, value interface{}, user *userModels.User) error {
	u.connect()
	defer u.disconnect()

	log := logger.ExtractLogger(
		UserRepositoryContext,
		logger.UserRepositoryKey,
	)
	log.Info("Finding user with filter", zap.String("filter", filter), zap.Any("value", value))

	err := u.db.Model(user).Where(filter, value).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			log.Info("No such user", zap.Any("user", user))
		} else {
			log.Error("Got an error", zap.Error(err))
		}
		return err
	}
	log.Info("Found user", zap.Any("user", user))
	return nil
}

func (u *userDB) Insert(user *userModels.User) error {
	u.connect()
	defer u.disconnect()

	logger.ExtractLogger(
		UserRepositoryContext,
		logger.UserRepositoryKey,
	).Info("Inserting user", zap.Any("user", user))
	_, err := u.db.Model(user).Insert()
	if err != nil {
		return err
	}
	return nil
}

func (u *userDB) Update(user *userModels.User) error {
	u.connect()
	defer u.disconnect()

	logger.ExtractLogger(
		UserRepositoryContext,
		logger.UserRepositoryKey,
	).Info("Updating user", zap.Any("user", user))
	_, err := u.db.Model(user).Update()
	if err != nil {
		return err
	}
	return nil
}

func (u *userDB) Delete(user *userModels.User) error {
	u.connect()
	defer u.disconnect()

	logger.ExtractLogger(
		UserRepositoryContext,
		logger.UserRepositoryKey,
	).Info("Deleting user", zap.Any("user", user))
	_, err := u.db.Model(user).Delete()
	if err != nil {
		return err
	}
	return nil
}

func (u *userDB) List(users *[]userModels.User) error {
	u.connect()
	defer u.disconnect()

	logger.ExtractLogger(
		UserRepositoryContext,
		logger.UserRepositoryKey,
	).Info("Getting user list")
	err := u.db.Model(users).Select()
	if err != nil {
		return err
	}
	return nil
}

func (u *userDB) connect() {
	pg := pg.Connect(&pg.Options{
		Addr:     conf.Config.UserDBHost,
		User:     conf.Config.UserDBLogin,
		Password: conf.Config.UserDBPassword,
		Database: conf.Config.UserDBName,
	})

	u.db = pg
}

func (u *userDB) disconnect() {
	err := u.db.Close()
	if err != nil {
		logger.ExtractLogger(
			UserRepositoryContext,
			logger.UserRepositoryKey,
		).Error("Error while PostgreDB disconnecting", zap.Error(err))
	}
}
