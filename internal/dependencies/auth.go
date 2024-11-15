package dependencies

import (
	"context"
	auth_types "vitalik_backend/internal/pkg/services/auth_service/types"
)

type IAuthService interface {
	Register(ctx context.Context, args auth_types.AuthArgs) error
	Login(ctx context.Context, args auth_types.AuthArgs) (*auth_types.LoginResponse, error)
	Authenticate(ctx context.Context, tokenString string) (string, error)
}
