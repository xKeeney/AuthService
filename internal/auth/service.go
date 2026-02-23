package auth

type authService struct {
	authRepo *authRepository
}

func InitAuthService(authRepo *authRepository) *authService {
	return &authService{
		authRepo: authRepo,
	}
}

func (s *authService) CreateUser(uuid, email, passwordHash, status string) error {
	return s.authRepo.CreateUser(uuid, email, passwordHash, status)
}
