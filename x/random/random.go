package random

import (
	"github.com/google/uuid"
)

func RandomString() string {
	return uuid.NewString()
}
