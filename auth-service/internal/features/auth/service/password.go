package usersservice

import "golang.org/x/crypto/bcrypt"

// hashPassword creates a bcrypt hash of the password.
// Uses bcrypt.DefaultCost (10 rounds).
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	
	return string(bytes), err
}

// checkPassword compares a password with its bcrypt hash.
func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}