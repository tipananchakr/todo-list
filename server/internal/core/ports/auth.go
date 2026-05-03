package ports

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash string, password string) error
}

type TokenManager interface {
	Generate(userID string) (string, error)
	Validate(token string) (string, error)
}
