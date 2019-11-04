package model

import (
	"github.com/FactomProject/factom"
)

// DID describes DID document
type DID struct {
	Network        string          `json:"network" form:"network" query:"network"`
	ManagementKeys []ManagementKey `json:"managementKeys" form:"managementKeys" query:"managementKeys"`
	DIDKeys        []DIDKey        `json:"didKeys" form:"didKeys" query:"didKeys"`
	Services       []Service       `json:"services" form:"services" query:"services"`
}

const (

	// DID Key Type enum
	KeyTypeECDSA = "ECDSASecp256k1VerificationKey"
	KeyTypeEdDSA = "Ed25519VerificationKey"
	KeyTypeRSA   = "RSAVerificationKey"

	// DID Entry Type enum
	EntryTypeCreate         = "DIDManagement"
	EntryTypeDeactivation   = "DIDDeactivation"
	EntryTypeUpdate         = "DIDUpdate"
	EntryTypeVersionUpgrade = "DIDMethodVersionUpgrade"

	// DID Key Purpose enum
	KeyPurposeAuthentication = "authentication"
	KeyPurposePublic         = "publicKey"

	// DID Network enum
	NetworkMainnet     = "mainnet"
	NetworkTestnet     = "testnet"
	NetworkUnspecified = ""
)

// ChainID calculates ChainID of DID chain on Factom
func (did *DID) ChainID() string {

	// TBD

	return ""

}

// AddDIDKey adds a DID key to the DID object
func (did *DID) AddDIDKey(alias string, purpose string, keyType string, controller string, priorityRequirement int) error {

	// TBD

	return nil

}

// AddManagementKey adds a new Management key to the DID object
func (did *DID) AddManagementKey(alias string, priority int, keyType string, controller string, priorityRequirement int) error {

	// TBD

	return nil

}

// AddService adds a new Service to the DID object
func (did *DID) AddService(alias string, serviceType string, endpoint string, priorityRequirement int) error {

	// TBD

	return nil

}

// SetMainnet sets the DID network to mainnet
func (did *DID) SetMainnet() error {

	// TBD

	return nil

}

// SetTestnet sets the DID network to testnet
func (did *DID) SetTestnet() error {

	// TBD

	return nil

}

// GenerateUpdateEntry prepares the update entry for writing on-chain
func (did *DID) GenerateUpdateEntry() (*factom.Entry, error) {

	fe := &factom.Entry{}

	// TBD

	return fe, nil

}

// ExportEncrypted exports JSON string with encrypted keys
func (did *DID) ExportEncrypted(password string) (string, error) {

	// TBD

	return "", nil

}

// RevokeDIDKey revokes DID Key from DID object
func (did *DID) RevokeDIDKey(alias string) error {

	// TBD

	return nil

}

// RevokeManagementKey revokes Management Key from DID object
func (did *DID) RevokeManagementKey(alias string) error {

	// TBD

	return nil

}

// RevokeService revokes Service from DID object
func (did *DID) RevokeService(alias string) error {

	// TBD

	return nil

}

// RotateDIDKey rotates DID Key from DID object
func (did *DID) RotateDIDKey(alias string) error {

	// TBD

	return nil

}

// RotateManagementKey rotates Management Key from DID object
func (did *DID) RotateManagementKey(alias string) error {

	// TBD

	return nil

}
