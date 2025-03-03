package xhash

import "golang.org/x/crypto/bcrypt"

const (
	MinCost     = bcrypt.MinCost     // the minimum allowable cost as passed in to GenerateFromPassword
	MaxCost     = bcrypt.MaxCost     // the maximum allowable cost as passed in to GenerateFromPassword
	DefaultCost = bcrypt.DefaultCost // the cost that will actually be set if a cost below MinCost is passed into GenerateFromPassword
)

// MinCost: 4
// MaxCost: 31
// DefaultCost: 10

func HashPassword(password string, cost int) (string, error) {
	hashedByte, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hashedByte), nil
}

func ComparePassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
