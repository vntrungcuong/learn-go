// Password utils
package util

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash      = errors.New("the encoded hash is not in the correct format")
	ErrPasswordMismatch = errors.New("password does not match")
)

// PasswordConfig for password hashing
type PasswordConfig struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

// Config for password hashing
var config = &PasswordConfig{
	time:    1,
	memory:  64 * 1024,
	threads: 4,
	keyLen:  32,
}

// HashPassword hashes a password using Argon2id.
func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, config.time, config.memory, config.threads, config.keyLen)

	// Format: salt.hash (base64)
	return fmt.Sprintf("%s.%s",
		base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(hash)), nil
}

// CheckPassword checks a password against a hashed password.s
func CheckPassword(password string, hashedPassword string) error {
	// 1. Split string format "salt.hash"
	parts := strings.Split(hashedPassword, ".")
	if len(parts) != 2 {
		return ErrInvalidHash
	}

	// 2. Decode salt and hash from base64
	salt, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return err
	}

	originalHash, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return err
	}

	// 3. Hash new passsword with old salt and config
	comparisonHash := argon2.IDKey([]byte(password), salt, config.time, config.memory, config.threads, config.keyLen)

	// 4. Compare by ConstantTimeCompare to prevent Timing Attacks
	if subtle.ConstantTimeCompare(originalHash, comparisonHash) == 1 {
		return nil
	}

	return ErrPasswordMismatch
}
