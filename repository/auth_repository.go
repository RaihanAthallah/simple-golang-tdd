package repository

type AuthRepository interface {
}

type authRepositoryImpl struct {
}

func NewAuthRepository() AuthRepository {
	return &authRepositoryImpl{}
}
