package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func PasswordHashing(password string) string {
	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
	}

	return string(passwordHashed)
}
