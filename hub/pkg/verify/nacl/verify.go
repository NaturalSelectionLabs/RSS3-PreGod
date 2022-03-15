package nacl

import (
	"fmt"

	"golang.org/x/crypto/nacl/sign"
)

const (
	cryptoSignBytes       = 64
	cryptoSignPubkeyBytes = 32
)

func Verify(msg, sig, pubkey []byte) (bool, error) {
	if (len(sig) != cryptoSignBytes) || (len(pubkey) != cryptoSignPubkeyBytes) {
		return false, fmt.Errorf("sig must be 64 bytes and pubkey must be 32 bytes, currently len(sig)=%d and len(pubkey)=%d", len(sig), len(pubkey))
	}

	var pubkeyBytes [cryptoSignPubkeyBytes]byte

	copy(pubkeyBytes[:], pubkey[:cryptoSignPubkeyBytes])

	_, ok := sign.Open(nil, append(sig, msg...), &pubkeyBytes)
	if !ok {
		return false, fmt.Errorf("signature nacl verification failed")
	}

	return ok, nil
}
