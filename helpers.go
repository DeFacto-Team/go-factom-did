package factomdid

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/FactomProject/factom"
)

var rng *rand.Rand
var rngMtx sync.Mutex

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
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

	entry := &factom.Entry{}
	entry.ExtIDs = append(entry.ExtIDs, []byte(EntryTypeCreate))
	entry.ExtIDs = append(entry.ExtIDs, []byte(EntrySchemaV100))

	for _, extid := range extIDs {
		entry.ExtIDs = append(entry.ExtIDs, extid)
	}

	chain := factom.NewChain(entry)

	return chain.ChainID, nil

}

// Calculates entry size
func calculateEntrySize(entry *factom.Entry) int {

	size := len(entry.Content)

	if len(entry.ExtIDs) > 0 {
		for _, extid := range entry.ExtIDs {
			size += len(extid)
		}
	}

	return size

}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
