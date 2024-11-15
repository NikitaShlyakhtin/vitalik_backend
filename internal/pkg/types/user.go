package types

type User struct {
	UserID         string `json:"user_id"`
	HashedPassword string `json:"hashed_password"`
}
