package conf

import (
	"discountDealer/logger"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"reflect"
)

const envFile = ".env"

type definition struct {
	AuthScheme string `env:"AUTH_SCHEME"`
	JWTKey     []byte `env:"JWT_KEY"`

	UserDBHost     string `env:"USER_DB_HOST"`
	UserDBLogin    string `env:"USER_DB_LOGIN"`
	UserDBPassword string `env:"USER_DB_PASSWORD"`
	UserDBName     string `env:"USER_DB_DBNAME"`

	ProductsDBHost       string `env:"PRODUCTS_DB_HOST"`
	ProductsDBLogin      string `env:"PRODUCTS_DB_LOGIN"`
	ProductsDBPassword   string `env:"PRODUCTS_DB_PASSWORD"`
	ProductsDBName       string `env:"PRODUCTS_DB_DBNAME"`
	ProductsDBCollection string `env:"PRODUCTS_DB_COLLECTION"`

	ManticoreSearchURL string `env:"SEARCH_READ_URL"`
	ManticoreQueryURL  string `env:"SEARCH_SQL_URL"`
	ManticoreIndex     string `env:"SEARCH_INDEX"`
}

var Config *definition

func load() *definition {
	log := logger.New("Config")
	def := &definition{}

	log.Info("loading .env file to OS env")
	if err := godotenv.Overload(envFile); err != nil {
		log.Fatal("Can't load .env file", zap.Error(err))
		return nil
	}

	log.Info("Reading OS env")
	err := env.ParseWithFuncs(def, map[reflect.Type]env.ParserFunc{
		reflect.TypeOf([]byte{0}): func(v string) (interface{}, error) {
			return []byte(v), nil
		},
	})
	if err != nil {
		log.Fatal("Can't read OS env", zap.Error(err))
		return nil
	}

	return def
}

func init() {
	Config = load()
}
