package factomdid

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/FactomProject/factom"
	"gopkg.in/go-playground/validator.v9"
)

var rng *rand.Rand
var rngMtx sync.Mutex
var validate *validator.Validate

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	validate = validator.New()
}

// generateNonce() generates random 32 bytes nonce
func generateNonce() []byte {
	nonce := make([]byte, 64)
	rngMtx.Lock()
	rng.Read(nonce)
	rngMtx.Unlock()

	return nonce
}

// Calculates ChainID based on DID extID
func calculateChainID(extIDs [][]byte) (string, error) {
	if len(extIDs) == 0 {
		return "", fmt.Errorf("extIDs should not be empty")
	}

	prefix := [][]byte{[]byte(EntryTypeCreate), []byte(EntrySchemaV100)}
	return factom.ChainIDFromFields(append(prefix, extIDs...)), nil
}

// Calculates entry size
func calculateEntrySize(entry *factom.Entry) int {
	size := len(entry.Content)

	for _, extid := range entry.ExtIDs {
		size += len(extid)
	}

	return size
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
