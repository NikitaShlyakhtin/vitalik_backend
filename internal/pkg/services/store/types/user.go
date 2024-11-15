package store_types

import "vitalik_backend/internal/pkg/types"

type User struct {
	UserID         string `db:"users.user_id"`
	HashedPassword string `db:"users.hashed_password"`
}

func MapToUserStore(user *types.User) *User {
	return &User{
		UserID:         user.UserID,
		HashedPassword: user.HashedPassword,
	}
}

func MapFromUserStore(user *User) *types.User {
	return &types.User{
		UserID:         user.UserID,
		HashedPassword: user.HashedPassword,
	}
}
