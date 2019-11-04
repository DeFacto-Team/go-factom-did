package encryptor

import "github.com/DeFacto-Team/go-factom-did/model"

// Encryptor component encrypts and decrypts keys with given password
type Encryptor interface {
	EncryptKeys(did *model.DID, password string) string
	DecryptKeys(encrypted string, password string) (*model.DID, error)
}
