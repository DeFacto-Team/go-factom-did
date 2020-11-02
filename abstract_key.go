package factomdid

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"

	"github.com/frankbraun/dcrd/dcrec/secp256k1"
	"gopkg.in/go-playground/validator.v9"
)

// AbstractKey represents the common fields and functionality in a ManagementKey and a DIDKey.
type AbstractKey struct {
	Alias               string `json:"alias" form:"alias" query:"alias" validate:"required"`
	KeyType             string `json:"keyType" form:"keyType" query:"keyType" validate:"required,oneof=ECDSASecp256k1VerificationKey Ed25519VerificationKey RSAVerificationKey"`
	Controller          string `json:"controller" form:"controller" query:"controller" validate:"required"`
	PriorityRequirement *int   `json:"priorityRequirement" form:"priorityRequirement" query:"omitempty,priorityRequirement"`
	PublicKey           []byte `json:"publicKey" form:"publicKey" query:"publicKey" validate:"required"`
	PrivateKey          []byte `json:"privateKey" form:"privateKey" query:"privateKey" validate:"required"`
}

const (
	// KeyTypeECDSA is a constant for "ECDSASecp256k1VerificationKey"
	KeyTypeECDSA = "ECDSASecp256k1VerificationKey"
	// KeyTypeEdDSA is a constant for "Ed25519VerificationKey"
	KeyTypeEdDSA = "Ed25519VerificationKey"
	// KeyTypeRSA is a constant for "RSAVerificationKey"
	KeyTypeRSA = "RSAVerificationKey"
)

// Sign a message with the existing private key and signature type
// The message is hashed (SHA-256) before being signed
func (key *AbstractKey) Sign(message []byte) ([]byte, error) {

	validate := validator.New()
	err := validate.StructPartial(key, "PrivateKey", "KeyType")
	if err != nil {
		return nil, err
	}

	// signing sha256 hash of the message
	hashed := sha256.Sum256(message)

	switch key.KeyType {
	case KeyTypeECDSA:

		privKey, _ := secp256k1.PrivKeyFromBytes(secp256k1.S256(), key.PrivateKey)

		signature, err := privKey.Sign(hashed[:])
		if err != nil {
			return nil, err
		}

		signatureBytes := signature.Serialize()
		if err != nil {
			return nil, err
		}

		return signatureBytes, nil

	case KeyTypeEdDSA:

		signature := ed25519.Sign(key.PrivateKey, hashed[:])

		return signature, nil

	case KeyTypeRSA:

		p, err := x509.ParsePKCS1PrivateKey(key.PrivateKey)
		if err != nil {
			return nil, err
		}

		signature, err := rsa.SignPKCS1v15(rand.Reader, p, crypto.SHA256, hashed[:])
		if err != nil {
			return nil, err
		}

		return signature, nil

	}

	return nil, fmt.Errorf("Invalid key.KeyType")

}

// Verify the signature of the given message.
// The message is hashed (SHA-256) before being verified
func (key *AbstractKey) Verify(message []byte, signature []byte) (bool, error) {

	validate := validator.New()
	err := validate.StructPartial(key, "PublicKey", "KeyType")
	if err != nil {
		return false, err
	}

	// signing sha256 hash of the message
	hashed := sha256.Sum256(message)

	switch key.KeyType {
	case KeyTypeECDSA:

		pubKey, err := secp256k1.ParsePubKey(key.PublicKey, secp256k1.S256())
		if err != nil {
			return false, err
		}

		sig, err := secp256k1.ParseSignature(signature, secp256k1.S256())
		if err != nil {
			return false, err
		}

		v := sig.Verify(hashed[:], pubKey)

		return v, nil

	case KeyTypeEdDSA:

		v := ed25519.Verify(key.PublicKey, hashed[:], signature)

		return v, nil

	case KeyTypeRSA:

		p, err := x509.ParsePKCS1PublicKey(key.PublicKey)
		if err != nil {
			return false, err
		}

		err = rsa.VerifyPKCS1v15(p, crypto.SHA256, hashed[:], signature)
		if err != nil {
			return false, err
		}
		return true, nil

	}

	return false, fmt.Errorf("Invalid key.KeyType")

}

// internal helper function
// GenerateRandomKeys generates random keypair of AbstractKey.KeyType
func (key *AbstractKey) generateRandomKeys() *AbstractKey {

	switch key.KeyType {
	case KeyTypeECDSA:

		// Generate random ecdsa key
		privateKey, _ := secp256k1.GeneratePrivateKey(secp256k1.S256())
		key.PrivateKey = privateKey.Serialize()

		// Get corresponding public key
		x, y := privateKey.Public()
		key.PublicKey = secp256k1.NewPublicKey(secp256k1.S256(), x, y).Serialize()

	case KeyTypeEdDSA:

		// Generate random eddsa key
		key.PublicKey, key.PrivateKey, _ = ed25519.GenerateKey(rand.Reader)

	case KeyTypeRSA:

		// Generate random rsa key
		privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)

		// Get corresponding public key
		publicKey := privateKey.Public()

		// Marshall private key to []byte
		key.PrivateKey = x509.MarshalPKCS1PrivateKey(privateKey)

		// Marshal public key to []byte
		key.PublicKey = x509.MarshalPKCS1PublicKey(publicKey.(*rsa.PublicKey))

	}

	return key

}

// SetPriorityRequirement sets PriorityRequirement for AbstractKey
func (key *AbstractKey) SetPriorityRequirement(i int) *AbstractKey {
	key.PriorityRequirement = &i
	return key
}
