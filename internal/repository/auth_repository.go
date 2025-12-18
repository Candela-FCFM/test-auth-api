package repository

type IAuthRepository interface {
}

type AuthRepository struct {
}

func NewAuthRepository() IAuthRepository {
	return &AuthRepository{}
}
