# Factom DID (Decentralized Identities) Golang Lib
The first W3C compatible Golang lib for Factom DID (Decentralized Identities).<br />
Overall overview of this library design:
https://docs.google.com/document/d/19JvrOyVKk6hHTtUuWoH1HS3vGQAu1MjbwsPgj33tIFo/edit

## What this lib can do

* **Generate DID** (DID keys, Management keys, Services)
* **Update DID** and calculate difference between initial and updated DID documents for on-chain update
  * Add/revoke DID keys
  * Add/revoke Management keys
  * Add/revoke Services
* **Deactivate DID**
* **Generate `*factom.Entry{}`** for `DIDManagement`, `DIDUpdate`, `DIDDeactivation` (fully compatible with <a href="https://github.com/FactomProject/factom">Factom Golang Lib)</a>
* **Advanced DID validation**
  * Full validation of DID, DIDKey, ManagementKey, Service structs before generating on-chain entry
  * At least one DIDkey and one ManagementKey required for DID creation
  * At least one ManagementKey with `priority = 0` required for DID creation
  * At least one ManagementKey with `priority = 0` required to remain after DID update
  * DID Deactivation only using ManagementKey with `priority = 0` required
  * Check for no duplicates of aliases among DID and Management keys
  * Check for no duplicates of services aliases
  * Dynamic calculation of max required priority for DID Update and comparing if signing Management Key is equal or lower than the required priority
  * Max Factom Entry size (10KB) validation
* **Sign** and **Verify**
  * **Signing and verifying** any messages with **DID keys** and **Management Keys**
  * **Built-in automatic signing** of generated `DIDUpdate`, `DIDDeactivation` entries
  * **Supported signatures:** `ECDSASecp256k1`, `Ed25519`, `RSA`
* **Automatic public keys conversion** into on-chain format (`Base58` for `ECDSASecp256k1` and `Ed25519`, `PEM` for `RSA`)

## Functions

* **DID**
  * NewDID()
  * String()
  * GetChainID()
  * SetMainnet()
  * SetTestnet()
  * AddDIDKey(key *DIDKEY)
  * AddManagementKey(key *ManagementKey)
  * AddService(service *Service)
  * RevokeDIDKey(alias string)
  * RevokeManagementKey(alias string)
  * RevokeService(alias string)
  * Create()
  * Update(update *DID, signingKeyAlias string)
  * Deactivate(signingKeyAlias string)
  * Validate()
  * Copy()
* **DIDKey**
  * NewDIDKey(alias string, keyType string)
  * AddPurpose(purpose string)
  * SetPriorityRequirement(i int)
  * Sign(message []byte)
  * Verify(message []byte, signature []byte)
* **ManagementKey**
  * NewManagementKey(alias string, keyType string, priority int)
  * SetPriorityRequirement(i int)
  * Sign(message []byte)
  * Verify(message []byte, signature []byte)
* **Service**
  * NewService(alias string, serviceType string, endpoint string)
  * SetPriorityRequirement(i int)

## Enums
```golang
// DID Factom Entry Type enum

// EntryTypeCreate is ExtID "DIDManagement" used in the first entry of DID chain
EntryTypeCreate = "DIDManagement"
// EntryTypeDeactivation is ExtID "DIDDeactivation" used in deactivation entry
EntryTypeDeactivation = "DIDDeactivation"
// EntryTypeUpdate is ExtID "DIDUpdate" used in update entry
EntryTypeUpdate = "DIDUpdate"
// EntryTypeVersionUpgrade is ExtID "DIDMethodVersionUpgrade" used in version upgrade entry
EntryTypeVersionUpgrade = "DIDMethodVersionUpgrade"

// DID Network enum

// NetworkMainnet is "mainnet"
NetworkMainnet = "mainnet"
// NetworkTestnet is "testnet"
NetworkTestnet = "testnet"
// NetworkUnspecified is empty string
NetworkUnspecified = ""

// DIDMethodName is method name for Factom DID, used in DID Document
DIDMethodName = "did:factom"
// MaxEntrySize is maximum size of Factom Entry
MaxEntrySize = 10240

// Versions

// EntrySchemaV100 is version 1.0.0 of entry schema
EntrySchemaV100 = "1.0.0"
// DIDMethodSpecV020 is version 0.2.0 of DID specification
DIDMethodSpecV020 = "0.2.0"

// Latest versions

// LatestEntrySchema is latest available entry schema
LatestEntrySchema = EntrySchemaV100
// LatestDIDMethodSpec is latest available DID specification version
LatestDIDMethodSpec = DIDMethodSpecV020

// KeyTypeECDSA is a constant for "ECDSASecp256k1VerificationKey"
KeyTypeECDSA = "ECDSASecp256k1VerificationKey"
// KeyTypeEdDSA is a constant for "Ed25519VerificationKey"
KeyTypeEdDSA = "Ed25519VerificationKey"
// KeyTypeRSA is a constant for "RSAVerificationKey"
KeyTypeRSA = "RSAVerificationKey"
```

## Example

### Create new DID document
```golang
// Generate blank DID document
// factom.ExtIDs with nonce are automatically generated in did.ExtIDs
did := factomdid.NewDID()

// Generate DID key
// NewDIDKey(alias, keyType)
didKey, err := factomdid.NewDIDKey("did-key-alias", factomdid.KeyTypeECDSA)
if err != nil {
  // handle error
}
didKey.AddPurpose(KeyPurposePublicKey)
didKey.SetPriorityRequirement(0)
did.AddDIDKey(didKey)

// Generate Management Key
// NewManagementKey(alias, keyType, priority)
mgmtKey, err := factomdid.NewManagementKey("mgmt-key-alias", factomdid.KeyTypeRSA, 0)
if err != nil {
  // handle error
}
did.AddManagementKey(mgmtKey)

// Generate Service
// NewService(alias, serviceType, endpoint)
service, err := factomdid.NewService("service-alias", "KYC", "https://kyc.example.com")
if err != nil {
  // handle error
}
service.SetPriorityRequirement(0)
did.AddService(service)

// Generate DIDManagement factom.Entry
entry, err := did.Create()
if err != nil {
  // handle error
}

// publish entry (*factom.Entry) on-chain using Factom Golang lib or Factom Open API
wallet.CommitRevealEntry(entry)
```

### Update existing DID document
```golang
...
// continuation of the code above
// make a copy of DID document for updating it
// we need a copy to calculate difference between 2 DID documents later
update := did.Copy()

// revoke DID key
update.RevokeDIDKey("did-key-alias")

// add second Management key
mgmtKey2, err := factomdid.NewManagementKey("second-mgmt-key-alias", factomdid.KeyTypeEdDSA, 1)
if err != nil {
  // handle error
}
update.AddManagementKey(mgmtKey)

// Generate DIDUpdate Factom entry, signed with Management Key "mgmt-key-alias"
// DID.Update(updatedDIDDocument, signingKeyAlias)
entry, err := did.Update(update, "mgmt-key-alias")

// publish entry (*factom.Entry) on-chain using Factom Golang lib or Factom Open API
wallet.CommitRevealEntry(entry)
```

### Deactivate DID document
```golang
...
// continuation of the code above
// Generate DIDDeactivation Factom entry, signed with Management Key "mgmt-key-alias"
// DID.Deactivate(signingKeyAlias)
entry, err := did.Deactivate("did-key-alias")

// publish entry (*factom.Entry) on-chain using Factom Golang lib or Factom Open API
wallet.CommitRevealEntry(entry)
```
