package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

const hashFormat = "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"

type PasswordConfig struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

func (c PasswordConfig) DefaultConfig() *PasswordConfig {
	return &PasswordConfig{
		time:    1,
		memory:  64 * 1024,
		threads: 4,
		keyLen:  32,
	}
}

func GenerateHash(config *PasswordConfig, s string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(s), salt, config.time, config.memory, config.threads, config.keyLen)

	saltB64 := base64.RawStdEncoding.EncodeToString(salt)
	hashB64 := base64.RawStdEncoding.EncodeToString(hash)

	finalHash := fmt.Sprintf(hashFormat, argon2.Version, config.memory, config.time, config.threads, saltB64, hashB64)

	return finalHash, nil
}

func Validate(s, hash string) (bool, error) {

	hashParts := strings.Split(hash, "$")

	config := &PasswordConfig{}

	_, err := fmt.Sscanf(hashParts[3], "m=%d,t=%d,p=%d", &config.memory, &config.time, &config.threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(hashParts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(hashParts[5])

	config.keyLen = uint32(len(decodedHash))

	comparisonHash := argon2.IDKey([]byte(s), salt, config.time, config.memory, config.threads, config.keyLen)

	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1, nil
}
