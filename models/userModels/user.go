package userModels

import (
	"crypto/sha512"
	"discountDealer/conf"
	"discountDealer/x/random"
	"encoding/hex"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const (
	IdClaim         = "id"
	MethodClaim     = "method"
	ExpirationClaim = "exp"

	TokenAccess  = "access"
	TokenRefresh = "refresh"
)

type User struct {
	ID             *string   `json:"id,omitempty"`
	Username       string    `json:"username" validate:"required,max=32,username-unique"`
	Password       string    `json:"password" validate:"required,min=8,max=32"`
	Salt           string    `json:"-"`
	AdvertDisabled bool      `json:"advertDisabled,omitempty"`
	IsPremium      bool      `json:"isPremium,omitempty"`
	ReferalCode    string    `json:"referalCode,omitempty"`
	ReferalCount   int8      `json:"referalCount,omitempty"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

func (u *User) GenerateData() {
	u.setSalt()
	u.makeRefCode()
	u.Password = u.HashPassword(u.Password)
	u.setTimeStamps()
}

func (u *User) GenerateTokens() (*string, *string, error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessClaims := accessToken.Claims.(jwt.MapClaims)
	accessClaims[IdClaim] = u.ID
	accessClaims[MethodClaim] = TokenAccess
	accessClaims[ExpirationClaim] = time.Now().Add(time.Minute * 15).Unix()

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims[IdClaim] = u.ID
	refreshClaims[MethodClaim] = TokenRefresh
	refreshClaims[ExpirationClaim] = time.Now().Add(time.Hour * 72).Unix()

	a, err := accessToken.SignedString(conf.Config.JWTKey)
	if err != nil {
		return nil, nil, err
	}

	r, err := refreshToken.SignedString(conf.Config.JWTKey)
	if err != nil {
		return nil, nil, err
	}

	return &a, &r, nil
}

func (u *User) HashPassword(password string) string {
	hasher := sha512.New()
	hasher.Write([]byte(password))
	hash1 := hex.EncodeToString(hasher.Sum(nil))
	hasher.Write(append([]byte(hash1), []byte(u.Salt)...))
	hash2 := hex.EncodeToString(hasher.Sum(nil))
	return hash2
}

func (u *User) setTimeStamps() {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) setSalt() {
	u.Salt = random.RandomString()[:32]
}

func (u *User) makeRefCode() {
	u.ReferalCode = random.RandomString()[:5]
}
