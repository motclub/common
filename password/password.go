package password

import "golang.org/x/crypto/bcrypt"

func GeneratePasswordHash(plaintextPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), bcrypt.DefaultCost)
	return string(hash), err
}

func ValidatePassword(hashedPassword, plaintextPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plaintextPassword))
	return err
}
