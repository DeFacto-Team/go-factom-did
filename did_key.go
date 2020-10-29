package factomdid

import (
	"encoding/pem"
	"strings"

	"github.com/FactomProject/btcutil/base58"
	"gopkg.in/go-playground/validator.v9"
)

// DIDKey is a key used to sign updates for an existing DID
// DIDKey.Purpose may be publicKey, authentication or both
type DIDKey struct {
	AbstractKey
	Purpose []DIDKeyPurpose `json:"purpose" form:"purpose" query:"purpose" validate:"len=1|len=2,unique,required,dive"`
}

// DIDKeyPurpose shows what purpose(s) the key serves
type DIDKeyPurpose struct {
	Purpose string `json:"purpose" form:"purpose" query:"purpose" validate:"required,oneof=publicKey authentication"`
}

const (
	// KeyPurposeAuthentication is authentication purpose
	KeyPurposeAuthentication = "authentication"
	// KeyPurposePublic is publicKey purpose
	KeyPurposePublic = "publicKey"
	// OnChainPubKeyName is name of public key on-chain
	OnChainPubKeyName = "publicKeyBase58"
)

// NewDIDKey generates new DIDKey with alias and keyType
func NewDIDKey(alias string, keyType string) (*DIDKey, error) {

	key := &DIDKey{}
	key.Alias = alias
	key.KeyType = keyType

	// validate Alias and KeyType
	validate := validator.New()
	err := validate.StructPartial(key.AbstractKey, "Alias", "KeyType")
	if err != nil {
		return nil, err
	}

	key.generateRandomKeys()

	return key, nil

}

// AddPurpose adds purpose into DIDKey
func (didkey *DIDKey) AddPurpose(purpose string) (*DIDKey, error) {

	p := DIDKeyPurpose{Purpose: purpose}

	// validate Purpose
	validate := validator.New()
	err := validate.Struct(p)
	if err != nil {
		return nil, err
	}

	didkey.Purpose = append(didkey.Purpose, p)

	return didkey, nil

}

// helper function to convert DIDKey into DIDKeySchema
func (didkey *DIDKey) toSchema(DID string) (*DIDKeySchema, error) {

	// validate DIDKey
	validate := validator.New()
	err := validate.Struct(didkey)
	if err != nil {
		return nil, err
	}

	s := &DIDKeySchema{}
	s.Controller = didkey.Controller
	s.ID = strings.Join([]string{DID, didkey.Alias}, "#")
	s.PriorityRequirement = didkey.PriorityRequirement
	s.Type = didkey.KeyType

	for i := range didkey.Purpose {
		s.Purpose = append(s.Purpose, didkey.Purpose[i].Purpose)
	}

	if didkey.KeyType == KeyTypeRSA {
		s.PublicKeyPem = string(pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: didkey.PublicKey}))
	} else {
		s.PublicKeyBase58 = base58.Encode(didkey.PublicKey)
	}

	return s, nil

}

// helper function to convert DIDKey into RevokeIDSchema
func (didkey *DIDKey) toRevokeIDSchema(DID string) (*RevokeIDSchema, error) {

	// validate DID Key
	validate := validator.New()
	err := validate.StructPartial(didkey, "Alias")
	if err != nil {
		return nil, err
	}

	s := &RevokeIDSchema{}
	s.ID = strings.Join([]string{DID, didkey.Alias}, "#")

	return s, nil

}
