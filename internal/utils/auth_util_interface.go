package utils

type AuthUtil interface {
	CheckPassword(hashed, plain string) bool
	GenerateAccessToken(role, secret string) (string, error)
}

type DefaultAuthUtil struct{}

func (d DefaultAuthUtil) CheckPassword(hashed, plain string) bool {
	return CheckPassword(hashed, plain)
}

func (d DefaultAuthUtil) GenerateAccessToken(role string, secret string) (string, error) {
	return GenerateAccessToken(role, secret)
}
