package repository

type AuthRepository interface {
}

type authRepositoryImpl struct {
}

func NewAuthAuthRepository() AuthRepository {
	return &authRepositoryImpl{}
}
