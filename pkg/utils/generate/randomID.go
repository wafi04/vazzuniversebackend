package generate

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type IDOpts struct {
	Prefix *string
	Amount int
}

func GenerateRandomID(opts *IDOpts) string {
	amount := opts.Amount
	if amount <= 0 {
		amount = 16
	}

	// Generate random bytes
	randomBytes := make([]byte, amount)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(fmt.Sprintf("failed to generate random ID: %v", err))
	}

	// Convert to hex string
	randomHex := hex.EncodeToString(randomBytes)

	if opts.Prefix != nil && *opts.Prefix != "" {
		return *opts.Prefix + "-" + randomHex
	}

	return randomHex
}
