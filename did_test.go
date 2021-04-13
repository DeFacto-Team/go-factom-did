package factomdid

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDID(t *testing.T) {

	did := NewDID()

	assert.Equal(t, 3, len(did.ExtIDs))
	assert.Equal(t, []byte(EntryTypeCreate), did.ExtIDs[0])
	assert.NotEmpty(t, did.ExtIDs[1])
	assert.NotEmpty(t, did.String())
	assert.Equal(t, 64, len(did.GetChainID()))

}

func TestString(t *testing.T) {

	// DID strings
	didString := "did:factom:301a57c2e753d061928cf6b6a692ea052885d75d2af5640e9b5cbc8897bbf7d5"
	didStringTestnet := "did:factom:testnet:301a57c2e753d061928cf6b6a692ea052885d75d2af5640e9b5cbc8897bbf7d5"
	didStringMainnet := "did:factom:mainnet:301a57c2e753d061928cf6b6a692ea052885d75d2af5640e9b5cbc8897bbf7d5"

	did := &DID{}

	// test no network specified
	did.ID = didString
	assert.Equal(t, didString, did.String())

	// test mainnet specified
	did.SetMainnet()
	assert.Equal(t, didStringMainnet, did.String())

	// test testnet specified
	did.SetTestnet()
	assert.Equal(t, didStringTestnet, did.String())

}

func TestGetChainID(t *testing.T) {

	// DID strings
	chainID := "301a57c2e753d061928cf6b6a692ea052885d75d2af5640e9b5cbc8897bbf7d5"
	didString := "did:factom:301a57c2e753d061928cf6b6a692ea052885d75d2af5640e9b5cbc8897bbf7d5"
	didStringTestnet := "did:factom:testnet:301a57c2e753d061928cf6b6a692ea052885d75d2af5640e9b5cbc8897bbf7d5"
	didStringMainnet := "did:factom:mainnet:301a57c2e753d061928cf6b6a692ea052885d75d2af5640e9b5cbc8897bbf7d5"

	did := &DID{}

	// test no network specified
	did.ID = didString
	assert.Equal(t, chainID, did.GetChainID())

	// test mainnet specified
	did.ID = didStringMainnet
	assert.Equal(t, chainID, did.GetChainID())

	// test testnet specified
	did.ID = didStringTestnet
	assert.Equal(t, chainID, did.GetChainID())

}

func TestSetMainnet(t *testing.T) {

	did := &DID{}
	did.SetMainnet()
	assert.Equal(t, NetworkMainnet, did.Network)

}

func TestSetTestnet(t *testing.T) {

	did := &DID{}
	did.SetTestnet()
	assert.Equal(t, NetworkTestnet, did.Network)

}

func TestAddDIDKey(t *testing.T) {

	// add invalid DID key
	did := NewDID()
	invalidDIDKey := &DIDKey{}
	invalidDIDKey.Alias = "test"
	did, err := did.AddDIDKey(invalidDIDKey)
	assert.Nil(t, did)
	assert.Error(t, err)

	// add DID Key without Purpose
	did = NewDID()
	didKey, _ := NewDIDKey("test", KeyTypeRSA)
	did, err = did.AddDIDKey(didKey)
	assert.Nil(t, did)
	assert.Error(t, err)

	// add valid DID Key
	did = NewDID()
	didKey.AddPurpose(KeyPurposePublic)
	did, err = did.AddDIDKey(didKey)
	assert.Equal(t, 1, len(did.DIDKeys))
	assert.NoError(t, err)

	// try to add key with duplicate alias
	duplicateKey1, _ := NewDIDKey("test", KeyTypeECDSA)
	duplicateKey1.AddPurpose(KeyPurposePublic)
	did, err = did.AddDIDKey(duplicateKey1)
	assert.Nil(t, did)
	assert.Error(t, err)

}

func TestAddManagementKey(t *testing.T) {

	// add invalid Management key
	did := NewDID()
	invalidMKey := &ManagementKey{}
	invalidMKey.Alias = "test"
	did, err := did.AddManagementKey(invalidMKey)
	assert.Nil(t, did)
	assert.Error(t, err)

	// add valid Management Key
	did = NewDID()
	validMKey, _ := NewManagementKey("test", KeyTypeRSA, 0)
	did, err = did.AddManagementKey(validMKey)
	assert.Equal(t, 1, len(did.ManagementKeys))
	assert.NoError(t, err)

	// try to add key with duplicate alias
	duplicateKey1, _ := NewManagementKey("test", KeyTypeECDSA, 0)
	did, err = did.AddManagementKey(duplicateKey1)
	assert.Nil(t, did)
	assert.Error(t, err)

}

func TestAddService(t *testing.T) {

	var err error

	// add invalid Service
	did := NewDID()
	invalidService := &Service{}
	invalidService.Alias = "test"
	did, err = did.AddService(invalidService)
	assert.Nil(t, did)
	assert.Error(t, err)

	// add valid Service
	did = NewDID()
	validService, _ := NewService("test", "Test", "https://endpoint.com")
	did, err = did.AddService(validService)
	assert.Equal(t, 1, len(did.Services))
	assert.NoError(t, err)

	// try to add service with duplicate alias
	duplicateService1, _ := NewService("test", "Test2", "https://another.endpoint.com")
	did, err = did.AddService(duplicateService1)
	assert.Nil(t, did)
	assert.Error(t, err)

}

func TestRevokeDIDKey(t *testing.T) {

	var err error

	did := NewDID()
	didKey, _ := NewDIDKey("test", KeyTypeRSA)
	didKey.AddPurpose(KeyPurposePublic)
	did.AddDIDKey(didKey)

	// revoke not existent did key
	_, err = did.RevokeDIDKey("not-exist")
	assert.Error(t, err)

	// revoke existent did key
	did.RevokeDIDKey("test")
	assert.Equal(t, 0, len(did.DIDKeys))

	// revoke did key from empty slice []*DIDKey
	did, err = did.RevokeDIDKey("test")
	assert.Nil(t, did)
	assert.Error(t, err)

}

func TestCreate(t *testing.T) {

	did := NewDID()

	didKey, _ := NewDIDKey("default-did-key", KeyTypeECDSA)
	didKey.AddPurpose(KeyPurposePublic)
	mgmtKey, _ := NewManagementKey("default-mgmt-key", KeyTypeECDSA, 0)
	service, _ := NewService("demo", "Demo", "https://demo.com")

	did.AddDIDKey(didKey)
	did.AddManagementKey(mgmtKey)
	did.AddService(service)

	fe, err := did.Create()
	assert.NotNil(t, fe)
	assert.NoError(t, err)

}

func TestDeactivate(t *testing.T) {

	did := NewDID()

	didKey, _ := NewDIDKey("default-did-key", KeyTypeECDSA)
	didKey.AddPurpose(KeyPurposePublic)
	mgmtKey, _ := NewManagementKey("default-mgmt-key", KeyTypeECDSA, 0)
	mgmtKey1, _ := NewManagementKey("secondary-mgmt-key", KeyTypeECDSA, 1)

	did.AddDIDKey(didKey)
	did.AddManagementKey(mgmtKey)
	did.AddManagementKey(mgmtKey1)

	// try to deactivate with non existent key
	fe, err := did.Deactivate("not-existent-mgmt-key")
	assert.Nil(t, fe)
	assert.Error(t, err)

	// try to deactivate with priority=1 ManagementKey
	fe, err = did.Deactivate("secondary-mgmt-key")
	assert.Nil(t, fe)
	assert.Error(t, err)

	// try to deactivate with priority=0 ManagementKey
	fe, err = did.Deactivate("default-mgmt-key")
	assert.NotNil(t, fe)
	assert.NoError(t, err)

	// check signature
	signingKeyFullID := strings.Join([]string{did.ID, "default-mgmt-key"}, "#")
	v, err := did.ManagementKeys[0].Verify([]byte(strings.Join([]string{EntryTypeDeactivation, LatestEntrySchema, signingKeyFullID}, "")), fe.ExtIDs[3])
	assert.True(t, v)
	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {

	did := NewDID()

	didKey, _ := NewDIDKey("did1", KeyTypeECDSA)
	didKey.AddPurpose(KeyPurposeAuthentication)
	mKey1, _ := NewManagementKey("m1", KeyTypeECDSA, 0)
	mKey2, _ := NewManagementKey("m2", KeyTypeECDSA, 1)
	mKey1.SetPriorityRequirement(0)
	mKey2.SetPriorityRequirement(1)

	did.AddDIDKey(didKey)
	did.AddManagementKey(mKey1)
	did.AddManagementKey(mKey2)

	updatedDID := did.Copy()

	updatedDID.RevokeManagementKey("m2")
	mKey3, _ := NewManagementKey("m3", KeyTypeEdDSA, 1)
	mKey3.SetPriorityRequirement(0)
	updatedDID.AddManagementKey(mKey3)

	// try to update with 1 priority key (invalid)
	fe, err := did.Update(updatedDID, "m2")
	assert.Nil(t, fe)
	assert.Error(t, err)

	// update with 0 priority key (valid)
	fe, err = did.Update(updatedDID, "m1")
	assert.NotEmpty(t, fe)
	assert.NoError(t, err)

	// check signature
	signingKeyFullID := strings.Join([]string{did.ID, "m1"}, "#")
	v, err := did.ManagementKeys[0].Verify([]byte(strings.Join([]string{EntryTypeUpdate, LatestEntrySchema, signingKeyFullID, string(fe.Content)}, "")), fe.ExtIDs[3])
	assert.True(t, v)
	assert.NoError(t, err)
}

func TestValidate(t *testing.T) {

	var err error

	// DID document without keys
	did1 := NewDID()
	err = did1.Validate()
	assert.Error(t, err)

	// DID document without DIDKey
	mgmtKey, _ := NewManagementKey("1", KeyTypeRSA, 1)
	did1.AddManagementKey(mgmtKey)
	err = did1.Validate()
	assert.Error(t, err)

	// DID document with DIDKey and ManagementKey with Priority > 0
	didKey, _ := NewDIDKey("2", KeyTypeRSA)
	didKey.AddPurpose(KeyPurposeAuthentication)
	did1.AddDIDKey(didKey)
	err = did1.Validate()
	assert.Error(t, err)

	// DID document with DIDKey and ManagementKey
	mgmtKey2, _ := NewManagementKey("3", KeyTypeRSA, 0)
	did1.AddManagementKey(mgmtKey2)
	err = did1.Validate()
	assert.NoError(t, err)

	// DID document with duplicate key aliases
	did2 := NewDID()
	didKey1, _ := NewDIDKey("1", KeyTypeRSA)
	didKey1.AddPurpose(KeyPurposeAuthentication)
	didKey2, _ := NewDIDKey("1", KeyTypeRSA)
	didKey2.AddPurpose(KeyPurposeAuthentication)
	did2.AddDIDKey(didKey1)
	did2.AddDIDKey(didKey2)
	err = did2.Validate()
	assert.Error(t, err)

	// DID document with duplicate key aliases
	did3 := NewDID()
	didKey3, _ := NewManagementKey("2", KeyTypeRSA, 0)
	s1, _ := NewService("s1", "KYC", "https://kyc.com")
	s2, _ := NewService("s1", "KYC", "https://kyc.com")
	did3.AddDIDKey(didKey1)
	did3.AddManagementKey(didKey3)
	did3.AddService(s1)
	did3.AddService(s2)
	err = did3.Validate()
	assert.Error(t, err)

}
