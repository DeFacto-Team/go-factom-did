package factomdid

import (
	"testing"

	"github.com/FactomProject/factom"
	"github.com/stretchr/testify/assert"
)

func TestGenerateNonce(t *testing.T) {

	n1 := generateNonce()
	n2 := generateNonce()

	assert.NotEmpty(t, n1)
	assert.NotEmpty(t, n2)
	assert.NotEqual(t, n1, n2)

}

func TestCalculateChainID(t *testing.T) {

	d := &DID{}
	d.ExtIDs = append(d.ExtIDs, []byte("a32bf0ca979d58cd6970def0fcc666c1f4c262532ed53851bf041be204d2bd92"))

	chainID, err := calculateChainID(d.ExtIDs)

	assert.NoError(t, err)
	assert.Equal(t, "301a57c2e753d061928cf6b6a692ea052885d75d2af5640e9b5cbc8897bbf7d5", chainID)

}

func TestCalculateEntrySize(t *testing.T) {

	entry := &factom.Entry{}
	entry.Content = []byte("1")
	entry.ExtIDs = append(entry.ExtIDs, []byte("a32bf0ca979d58cd6970def0fcc666c1f4c262532ed53851bf041be204d2bd92"))

	s := calculateEntrySize(entry)

	assert.Equal(t, 65, s)

}
