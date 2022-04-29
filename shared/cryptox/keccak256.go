package cryptox

import "golang.org/x/crypto/sha3"

func Keccak256(p []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(p)

	return hash.Sum(nil)
}
