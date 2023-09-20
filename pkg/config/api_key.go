package config

import (
	"errors"
	"fmt"
	"math/rand"
	"os"

	"golang.org/x/crypto/bcrypt"
)

const (
	API_KEY_LENGHT = 32
)

func CheckAPIKey() ([]byte, error) {
	hashedKey, err := readAPIKey()
	if err != nil {
		return nil, err
	}

	if len(hashedKey) != 0 {
		return hashedKey, nil
	}

	hashedKey, key, err := createNewAPIKey()
	if err != nil {
		return nil, err
	}

	fmt.Println("Your new API key is:", key)
	return hashedKey, nil
}

func createNewAPIKey() ([]byte, string, error) {
	hashedKey, keyStr := generateHashedAPIKey()

	err := storeAPIKey(hashedKey)
	if err != nil {
		return nil, "", err
	}

	return hashedKey, keyStr, nil
}

func storeAPIKey(key []byte) error {
	return os.WriteFile("data", key, 0644)
}

func generateHashedAPIKey() ([]byte, string) {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, API_KEY_LENGHT)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}

	keyStr := string(s)

	hashedKey, _ := bcrypt.GenerateFromPassword([]byte(keyStr), bcrypt.DefaultCost)

	return hashedKey, keyStr
}

func readAPIKey() ([]byte, error) {
	err := createDataFile()
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile("data")
	if err != nil {
		return nil, err
	}

	return content, nil
}

func createDataFile() error {
	if _, err := os.Stat("data"); errors.Is(err, os.ErrNotExist) {
		_, err := os.Create("data")
		if err != nil {
			return err
		}
	}

	return nil
}
