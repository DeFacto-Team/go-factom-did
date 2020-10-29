package factomdid

import (
	"encoding/pem"
	"strings"

	"github.com/FactomProject/btcutil/base58"
	"gopkg.in/go-playground/validator.v9"
)

// ManagementKey is a key used to sign updates for an existing DID
type ManagementKey struct {
	AbstractKey
	Priority int `json:"priority" form:"priority" query:"priority" validate:"min=0"`
}

// NewManagementKey generates new ManagementKey with alias and keyType
func NewManagementKey(alias string, keyType string, priority int) (*ManagementKey, error) {

	key := &ManagementKey{}
	key.Alias = alias
	key.KeyType = keyType
	key.Priority = priority

	// validate Alias and KeyType
	validate := validator.New()
	err := validate.StructPartial(key.AbstractKey, "Alias", "KeyType")
	if err != nil {
		return nil, err
	}
	// validate Priority
	err = validate.StructPartial(key, "Priority")
	if err != nil {
		return nil, err
	}

	key.generateRandomKeys()

	return key, nil

}

// helper function to convert ManagementKey into ManagementKeySchema
func (mgmtkey *ManagementKey) toSchema(DID string) (*ManagementKeySchema, error) {

	// validate ManagementKey
	validate := validator.New()
	err := validate.Struct(mgmtkey)
	if err != nil {
		return nil, err
	}

	s := &ManagementKeySchema{}
	s.Controller = mgmtkey.Controller
	s.ID = strings.Join([]string{DID, mgmtkey.Alias}, "#")
	s.Priority = mgmtkey.Priority
	s.PriorityRequirement = mgmtkey.PriorityRequirement
	s.Type = mgmtkey.KeyType

	if mgmtkey.KeyType == KeyTypeRSA {
		s.PublicKeyPem = string(pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: mgmtkey.PublicKey}))
	} else {
		s.PublicKeyBase58 = base58.Encode(mgmtkey.PublicKey)
	}

	return s, nil

}

// helper function to convert ManagementKey into RevokeIDSchema
func (mgmtkey *ManagementKey) toRevokeIDSchema(DID string) (*RevokeIDSchema, error) {

	// validate Management Key
	validate := validator.New()
	err := validate.StructPartial(mgmtkey, "Alias")
	if err != nil {
		return nil, err
	}

	s := &RevokeIDSchema{}
	s.ID = strings.Join([]string{DID, mgmtkey.Alias}, "#")

	return s, nil

}
