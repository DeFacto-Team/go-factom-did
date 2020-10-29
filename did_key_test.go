package factomdid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/validator.v9"
)

func TestNewDIDKey(t *testing.T) {

	// valid Alias and KeyType
	k1, err := NewDIDKey("test", KeyTypeECDSA)
	assert.NoError(t, err)
	assert.NotEmpty(t, k1.PrivateKey)
	assert.NotEmpty(t, k1.PublicKey)

	k2, err := NewDIDKey("test", KeyTypeEdDSA)
	assert.NoError(t, err)
	assert.NotEmpty(t, k2.PrivateKey)
	assert.NotEmpty(t, k2.PublicKey)

	k3, err := NewDIDKey("test", KeyTypeRSA)
	assert.NoError(t, err)
	assert.NotEmpty(t, k3.PrivateKey)
	assert.NotEmpty(t, k3.PublicKey)

	// invalid KeyType
	k4, err := NewDIDKey("test", "WrongKeyType")
	assert.Error(t, err)
	assert.Nil(t, k4)
	k5, err := NewDIDKey("test", "")
	assert.Error(t, err)
	assert.Nil(t, k5)

	// invalid Alias
	k6, err := NewDIDKey("", KeyTypeECDSA)
	assert.Error(t, err)
	assert.Nil(t, k6)

	// purposes tests
	validate := validator.New()
	k7, _ := NewDIDKey("test", KeyTypeECDSA)

	// 0. no purposes
	err = validate.StructPartial(k7, "Purpose")
	assert.Error(t, err)

	// 1. single valid purpose
	k7.AddPurpose(KeyPurposeAuthentication)
	err = validate.StructPartial(k7, "Purpose")
	assert.NoError(t, err)

	// 2. two valid purposes
	k7.AddPurpose(KeyPurposePublic)
	err = validate.StructPartial(k7, "Purpose")
	assert.NoError(t, err)

	// 3. add third valid purpose
	k7.AddPurpose(KeyPurposePublic)
	err = validate.StructPartial(k7, "Purpose")
	assert.Error(t, err)

	// 4. two duplicate valid purposes
	k9, _ := NewDIDKey("test", KeyTypeECDSA)
	k9.AddPurpose(KeyPurposePublic)
	k9.AddPurpose(KeyPurposePublic)
	err = validate.StructPartial(k9, "Purpose")
	assert.Error(t, err)

}

func TestAddPurpose(t *testing.T) {

	k, err := NewDIDKey("test", KeyTypeECDSA)

	// valid Purposes
	k, err = k.AddPurpose(KeyPurposeAuthentication)
	assert.NoError(t, err)
	k, err = k.AddPurpose(KeyPurposePublic)
	assert.NoError(t, err)

	// invalid Purposes
	k, err = k.AddPurpose("WrongKeyPurpose")
	assert.Error(t, err)
	assert.Nil(t, k)
	k, err = k.AddPurpose("")
	assert.Error(t, err)
	assert.Nil(t, k)

}

func TestDIDKeyToSchema(t *testing.T) {

	did := "did:factom:301a57c2e753d061928cf6b6a692ea052885d75d2af5640e9b5cbc8897bbf7d5"

	// 1. Test RSA key
	k, _ := NewDIDKey("test", KeyTypeRSA)

	// failed validation due to no purpose
	f, err := k.toSchema(did)
	assert.Nil(t, f)
	assert.Error(t, err)

	// add purpose
	k.AddPurpose(KeyPurposePublic)

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
	k2, _ := NewDIDKey("test", KeyTypeECDSA)

	// add purpose
	k2.AddPurpose(KeyPurposePublic)
	k2.AddPurpose(KeyPurposeAuthentication)

	// add placeholder controller
	k2.Controller = did

	s2, err := k2.toSchema(did)
	assert.NotEmpty(t, s2)
	assert.NoError(t, err)

}
