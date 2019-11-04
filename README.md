# Factom DID (Decentralized Identities) Golang Lib
Golang lib for Factom DID (Decentralized Identities)

## Structs

* `DID` (DID Document)
  - `DIDKey` (DID Key)
    - `DIDKeyPurpose` (DID Key Purpose: `PublicKey`, `AuthenticationKey`, both)
  - `ManagementKey` (DID Management Key)
  - `Service` (Service associated with a DID)

## Functions

### DID Document

* `DID.ChainID()`
Calculates ChainID of DID chain on Factom
* `DID.AddDIDKey(…)`
Adds a DID key to the DID object
* `DID.AddManagementKey(…)`
Adds a new Management key to the DID object
* `DID.AddService(…)`
Adds a new Service to the DID object
* `DID.RevokeDIDKey(…)`
Revokes DID key to the DID object
* `DID.RevokeManagementKey(…)`
Revokes Management key to the DID object
* `DID.RevokeService(…)`
Revokes Service to the DID object
* `DID.RotateDIDKey(…)`
Rotates DID key to the DID object
* `DID.RotateManagementKey(…)`
Rotates Management key to the DID object
* `DID.SetMainnet()`
Sets the DID network to mainnet
* `DID.SetTestnet()`
Sets the DID network to testnet
* `DID.GenerateUpdateEntry()`
Prepares the update entry for writing on-chain
* `DID.ExportEncrypted(…)`
Exports JSON string with encrypted keys

### DID Key

* `NewDIDKeyFromFactomEntry()`
Fills DIDKey struct from Factom entry
* `DIDKey.ConvertToFactomEntry()`
Converts DIDKey struct to Factom Entry
* `DIDKey.Sign(…)`
Signs a message with the existing private key and signature type
* `DIDKey.Verify(…)`
Verifies the signature of the given message

### Management Key

* `NewManagementKeyFromFactomEntry()`
Fills ManagementKey struct from Factom entry
* `ManagementKey.ConvertToFactomEntry()`
Converts ManagementKey struct to Factom Entry

### Service

* `NewServiceFromFactomEntry()`
Fills Service struct from Factom entry
* `Service.ConvertToFactomEntry()`
Converts Service struct to Factom Entry

### Encryptor
Encryptor component encrypts and decrypts keys with given password

* `EncryptKeys(…)`
* `DecryptKeys(…)`

## Enums

```golang
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
```
