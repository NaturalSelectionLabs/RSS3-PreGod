package ethers

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// Verifies the signature of a message from a given address.
func VerifyMessage(msg, sig, addr string) (bool, error) {
	// check format is correct
	if !common.IsHexAddress(addr) {
		return false, fmt.Errorf("Invalid hex address: %s", addr)
	}

	if len(sig) != 65 {
		return false, fmt.Errorf("Invalid signature length (must be of 65): %d", len(sig))
	}

	address := common.HexToAddress(addr)

	sigBytes, err := hexutil.Decode(sig)
	if err != nil {
		return false, fmt.Errorf("Decode sig string failed. Invalid signature: %s", err.Error())
	}

	// legacy issue:
	// https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L442
	if sigBytes[64] != 27 && sigBytes[64] != 28 {
		return false, fmt.Errorf("Invalid signature (V is not 27 or 28): V=%d", sigBytes[64])
	}

	sigBytes[64] -= 27

	recoveredPubkey, err := crypto.SigToPub(signHash([]byte(msg)), sigBytes)
	if err != nil || recoveredPubkey == nil {
		return false, fmt.Errorf("Convert sig to pub failed. Invalid signature: %s", err.Error())
	}
	// recoveredPubkeyBytes := crypto.FromECDSAPub(recoveredPubkey)
	recoveredAddress := crypto.PubkeyToAddress(*recoveredPubkey)
	success := address == recoveredAddress

	return success, nil
}

// signHash is a helper function that calculates a hash for the given message
// that can be safely used to calculate a signature from.
//
// The hash is calculated as
//   keccak256("\x19Ethereum Signed Message:\n"${message length}${message}).
//
// This gives context to the signed message and prevents signing of transactions.
func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)

	return crypto.Keccak256([]byte(msg))
}
