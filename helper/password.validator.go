package helper

import "golang.org/x/crypto/bcrypt"

// hashPassword: password received from db
// testPassword: password inputed from user, this value that we want to check

func VerifyPassword(hashPassword, testPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(testPassword))
}
