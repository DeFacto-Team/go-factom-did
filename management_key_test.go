package factomdid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewManagementKey(t *testing.T) {

	// valid Alias, KeyType and Priority
	k1, err := NewManagementKey("test", KeyTypeECDSA, 0)
	assert.NoError(t, err)
	assert.NotEmpty(t, k1.PrivateKey)
	assert.NotEmpty(t, k1.PublicKey)

	k2, err := NewManagementKey("test", KeyTypeEdDSA, 1)
	assert.NoError(t, err)
	assert.NotEmpty(t, k2.PrivateKey)
	assert.NotEmpty(t, k2.PublicKey)

	k3, err := NewManagementKey("test", KeyTypeRSA, 2)
	assert.NoError(t, err)
	assert.NotEmpty(t, k3.PrivateKey)
	assert.NotEmpty(t, k3.PublicKey)

	// invalid KeyType
	k4, err := NewManagementKey("test", "WrongKeyType", 0)
	assert.Error(t, err)
	assert.Nil(t, k4)
	k5, err := NewManagementKey("test", "", 0)
	assert.Error(t, err)
	assert.Nil(t, k5)

	// invalid Alias
	k6, err := NewManagementKey("", KeyTypeECDSA, 0)
	assert.Error(t, err)
	assert.Nil(t, k6)

	// invalid Priority
	k7, err := NewManagementKey("test", KeyTypeECDSA, -1)
	assert.Error(t, err)
	assert.Nil(t, k7)

}

func TestManagementKeyToSchema(t *testing.T) {

	did := "did:factom:301a57c2e753d061928cf6b6a692ea052885d75d2af5640e9b5cbc8897bbf7d5"

	// 1. Test RSA key
	k, _ := NewManagementKey("test", KeyTypeRSA, 0)

	// test priorityRequirement 0 not omited
	var i *int
	i = new(int)
	*i = 0
	k.PriorityRequirement = i

	// add placeholder controller
	k.Controller = did

	s, err := k.toSchema(did)
	assert.NotEmpty(t, s)
	assert.NoError(t, err)

	// 2. Test ECDSA key
	k2, _ := NewManagementKey("test", KeyTypeECDSA, 0)

	// add placeholder controller
	k2.Controller = did

	s2, err := k2.toSchema(did)
	assert.NotEmpty(t, s2)
	assert.NoError(t, err)

}
