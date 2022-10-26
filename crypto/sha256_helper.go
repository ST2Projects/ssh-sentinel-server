package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256sum(payload []byte) string {
	sha := sha256.New()
	sha.Write(payload)
	hash := sha.Sum(nil)

	return hex.EncodeToString(hash)
}
