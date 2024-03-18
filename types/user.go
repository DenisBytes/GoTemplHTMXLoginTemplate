package types

import (
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	Email       string
	LoggedIn    bool
	AccessToken string

	Account
}
