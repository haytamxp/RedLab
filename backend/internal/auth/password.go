package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes a plain password.
func HashPassword(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// CheckPassword compares a hashed password with a plain password.
func CheckPassword(hash, password string) bool {

	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)

	return err == nil
}