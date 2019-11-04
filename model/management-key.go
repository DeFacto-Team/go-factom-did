package model

import "github.com/FactomProject/factom"

// ManagementKey is a key used to sign updates for an existing DID
type ManagementKey struct {
	Alias               string `json:"alias" form:"alias" query:"alias"`
	Priority            int    `json:"priority" form:"priority" query:"priority"`
	KeyType             string `json:"keyType" form:"keyType" query:"keyType" validate:"oneof=ECDSA EdDSA RSA"`
	Controller          string `json:"controller" form:"controller" query:"controller"`
	PriorityRequirement int    `json:"priorityRequirement" form:"priorityRequirement" query:"priorityRequirement"`
	PrivateKey          string `json:"privateKey" form:"privateKey" query:"privateKey"`
	PublicKey           string `json:"publicKey" form:"publicKey" query:"publicKey"`
}

// NewManagementKeyFromFactomEntry fills ManagementKey struct from Factom Entry
func NewManagementKeyFromFactomEntry(fe *factom.Entry) *ManagementKey {

	mkey := &ManagementKey{}

	// TBD

	return mkey

}

// ConvertToFactomEntry converts ManagementKey struct to Factom Entry
func (mkey *ManagementKey) ConvertToFactomEntry() *factom.Entry {

	fe := &factom.Entry{}

	// TBD

	return fe

}
