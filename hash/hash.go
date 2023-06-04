package hash

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	saltAndHash := append(salt, hash...)
	return base64.StdEncoding.EncodeToString(saltAndHash), nil
}

func CheckPassword(password, encodedHash string) bool {
	saltAndHash, err := base64.StdEncoding.DecodeString(encodedHash)
	if err != nil {
		return false
	}
	salt := saltAndHash[:16]
	hash := saltAndHash[16:]

	testHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	return bytes.Equal(hash, testHash)
}
