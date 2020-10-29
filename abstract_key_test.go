package factomdid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/validator.v9"
)

// TestAbstractKeys covers key generation, sign and verify of all type of keys
func TestAbstractKeys(t *testing.T) {

	var err error
	var signature []byte
	var verify bool
	validate := validator.New()

	// test ecdsa key
	ecdsa := &AbstractKey{}
	ecdsa.KeyType = KeyTypeECDSA
	ecdsa.generateRandomKeys()

	err = validate.StructPartial(ecdsa, "PrivateKey", "PublicKey")
	assert.NoError(t, err)

	signature, err = ecdsa.Sign([]byte("Test"))
	assert.NoError(t, err)

	verify, err = ecdsa.Verify([]byte("Test"), signature)
	assert.NoError(t, err)
	assert.True(t, verify)

	// test eddsa key
	eddsa := &AbstractKey{}
	eddsa.KeyType = KeyTypeEdDSA
	eddsa.generateRandomKeys()

	err = validate.StructPartial(eddsa, "PrivateKey", "PublicKey")
	assert.NoError(t, err)

	signature, err = eddsa.Sign([]byte("Test"))
	assert.NoError(t, err)

	verify, err = eddsa.Verify([]byte("Test"), signature)
	assert.NoError(t, err)
	assert.True(t, verify)

	// test rsa key
	rsa := &AbstractKey{}
	rsa.KeyType = KeyTypeRSA
	rsa.generateRandomKeys()

	err = validate.StructPartial(rsa, "PrivateKey", "PublicKey")
	assert.NoError(t, err)

	signature, err = rsa.Sign([]byte("Test"))
	assert.NoError(t, err)

	verify, err = rsa.Verify([]byte("Test"), signature)
	assert.NoError(t, err)
	assert.True(t, verify)

}
