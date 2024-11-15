package auth_types

type AuthArgs struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
