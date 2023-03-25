package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

func PasswordHashing(password string) string {
	randInt, error := rand.Int(rand.Reader, big.NewInt(10))
	if error != nil {
		panic(error)
	}
	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(password), int(randInt.Int64()))
	if err != nil {
		fmt.Println(err)
	}

	return string(passwordHashed)
}

func PasswordVerify(hashedPassword string, password string) bool {
	error := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	fmt.Println(error)

	return error == nil
}
