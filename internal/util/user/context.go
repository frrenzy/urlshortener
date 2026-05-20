package user

import (
	"context"
	"errors"
)

func GetUserFromContext(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(UserContextKey).(int)
	if !ok {
		return -1, errors.New("can not get user from context")
	}

	return userID, nil
}
