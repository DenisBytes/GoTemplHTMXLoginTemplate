package view

import (
	"context"
	"fmt"

	"github.com/DenisBytes/GoTemplHTMXLoginTemplate/types"
)

func AuthenticatedUser(ctx context.Context) types.User {
	user, ok := ctx.Value("user").(types.User)
	if !ok {
		return types.User{}
	}
	fmt.Printf("USER %+v\n\n", user)
	return user
}