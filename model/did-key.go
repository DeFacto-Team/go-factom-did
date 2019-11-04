package model

import "github.com/FactomProject/factom"

// DIDKey is a key used to sign updates for an existing DID
type DIDKey struct {
	Alias               string          `json:"alias" form:"alias" query:"alias" validate:"required"`
	KeyType             string          `json:"keyType" form:"keyType" query:"keyType" validate:"required,oneof=ECDSA EdDSA RSA"`
	Controller          string          `json:"controller" form:"controller" query:"controller" validate:"required"`
	PrivateKey          string          `json:"privateKey" form:"privateKey" query:"privateKey"`
	PublicKey           string          `json:"publicKey" form:"publicKey" query:"publicKey"`
	Purpose             []DIDKeyPurpose `json:"purpose" form:"purpose" query:"purpose" validate:"dive,required"`
	PriorityRequirement int             `json:"priorityRequirement" form:"priorityRequirement" query:"priorityRequirement"`
}

// DIDKeyPurpose shows what purpose(s) the key serves
// DIDKey.Purpose may be PublicKey, AuthenticationKey or both
type DIDKeyPurpose struct {
	Purpose string `json:"purpose" form:"purpose" query:"purpose" validate:"oneof=publicKey authenticationKey"`
}

// NewDIDKeyFromFactomEntry fills DIDKey struct from Factom entry
func NewDIDKeyFromFactomEntry(fe *factom.Entry) *DIDKey {

	didkey := DIDKey{}

	// TBD

	return &didkey

}

// ConvertToFactomEntry converts DIDKey struct to Factom Entry
func (didkey *DIDKey) ConvertToFactomEntry() *factom.Entry {

	fe := &factom.Entry{}

	// TBD

	return fe

}

// Sign a message with the existing private key and signature type
// The message is hashed before being signed, with the provided hash function. The default hash function used is SHA-256.
func (didkey *DIDKey) Sign(message string, hashFunction string) []byte {

	// TBD

	return nil

}

// Verify the signature of the given message
func (didkey *DIDKey) Verify(message string, signature string, hashFunction string) []byte {

	// TBD

	return nil

}
