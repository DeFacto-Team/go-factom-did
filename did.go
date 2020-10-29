package factomdid

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/FactomProject/factom"
	"gopkg.in/go-playground/validator.v9"
)

// DID describes DID document
type DID struct {
	ID             string           `json:"id" form:"id" query:"id" validate:"required"`
	Network        string           `json:"network" form:"network" query:"network"`
	ManagementKeys []*ManagementKey `json:"managementKeys" form:"managementKeys" query:"managementKeys" validate:"required"`
	DIDKeys        []*DIDKey        `json:"didKeys" form:"didKeys" query:"didKeys" validate:"required"`
	Services       []*Service       `json:"services" form:"services" query:"services"`
	ExtIDs         [][]byte         `json:"extIDs" form:"extIDs" query:"extIDs"`
}

const (

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
)

// NewDID generates new blank DID document
// DID.ExtIDs is a helper field that stores ExtIDs to be written on-chain to get expected ChainID for new DID chain
// DID.ExtIDs is not a part of DID Document (JSON)
func NewDID() *DID {

	d := &DID{}
	d.ExtIDs = append(d.ExtIDs, []byte(EntryTypeCreate))
	d.ExtIDs = append(d.ExtIDs, []byte(LatestEntrySchema))
	d.ExtIDs = append(d.ExtIDs, generateNonce())

	// d.ExtIDs is not nil, so no need to check for error
	chainID, _ := calculateChainID(d.ExtIDs)
	d.ID = strings.Join([]string{DIDMethodName, chainID}, ":")

	return d

}

// The decentralized identifier, a 32 byte hexadecimal string
func (did *DID) String() string {

	if did.Network == NetworkUnspecified {
		return did.ID
	}

	return strings.Join([]string{DIDMethodName, did.Network, did.GetChainID()}, ":")

}

// GetChainID gets ChainID from DID string
func (did *DID) GetChainID() string {

	s := strings.Split(did.ID, ":")

	return s[len(s)-1]

}

// SetMainnet sets the DID network to mainnet
func (did *DID) SetMainnet() *DID {

	did.Network = NetworkMainnet

	return did

}

// SetTestnet sets the DID network to testnet
func (did *DID) SetTestnet() *DID {

	did.Network = NetworkTestnet

	return did

}

// AddDIDKey adds a DID key to the DID object
func (did *DID) AddDIDKey(key *DIDKey) (*DID, error) {

	key.Controller = did.ID

	validate := validator.New()
	// exclude PrivateKey from validation in case you have PublicKey only to verify signatures
	err := validate.StructExcept(key, "PrivateKey")
	if err != nil {
		return nil, err
	}

	did.DIDKeys = append(did.DIDKeys, key)

	return did, nil

}

// AddManagementKey adds a new Management key to the DID object
func (did *DID) AddManagementKey(key *ManagementKey) (*DID, error) {

	key.Controller = did.ID

	validate := validator.New()
	// exclude PrivateKey from validation in case you have PublicKey only to verify signatures
	err := validate.StructExcept(key, "PrivateKey")
	if err != nil {
		return nil, err
	}

	did.ManagementKeys = append(did.ManagementKeys, key)

	return did, nil

}

// AddService adds a new Service to the DID object
func (did *DID) AddService(service *Service) (*DID, error) {

	validate := validator.New()
	err := validate.Struct(service)
	if err != nil {
		return nil, err
	}

	did.Services = append(did.Services, service)

	return did, nil

}

// RevokeDIDKey revokes DID Key from DID object
func (did *DID) RevokeDIDKey(alias string) (*DID, error) {

	initLen := len(did.DIDKeys)

	if initLen == 0 {
		return nil, fmt.Errorf("No DIDKeys found in this DID document")
	}

	for i, v := range did.DIDKeys {
		if v.Alias == alias {
			did.DIDKeys = append(did.DIDKeys[:i], did.DIDKeys[i+1:]...)
		}
	}

	if initLen == len(did.DIDKeys) {
		return nil, fmt.Errorf("DIDKey with alias %s not found", alias)
	}

	return did, nil

}

// RevokeManagementKey revokes Management Key from DID object
func (did *DID) RevokeManagementKey(alias string) (*DID, error) {

	initLen := len(did.ManagementKeys)

	if initLen == 0 {
		return nil, fmt.Errorf("No ManagementKeys found in this DID document")
	}

	for i, v := range did.ManagementKeys {
		if v.Alias == alias {
			did.ManagementKeys = append(did.ManagementKeys[:i], did.ManagementKeys[i+1:]...)
		}
	}

	if initLen == len(did.ManagementKeys) {
		return nil, fmt.Errorf("ManagementKey with alias %s not found", alias)
	}

	return did, nil

}

// RevokeService revokes Service from DID object
func (did *DID) RevokeService(alias string) (*DID, error) {

	initLen := len(did.Services)

	if initLen == 0 {
		return nil, fmt.Errorf("No Services found in this DID document")
	}

	for i, v := range did.Services {
		if v.Alias == alias {
			did.Services = append(did.Services[:i], did.Services[i+1:]...)
		}
	}

	if initLen == len(did.Services) {
		return nil, fmt.Errorf("Service with alias %s not found", alias)
	}

	return did, nil

}

// Create generates DIDManagement Factom Entry from DID document
func (did *DID) Create() (*factom.Entry, error) {

	// validate DID document
	err := did.Validate()
	if err != nil {
		return nil, err
	}

	// fill entry schema
	s := &DIDManagementEntrySchema{}
	sD := &DIDKeySchema{}
	sM := &ManagementKeySchema{}
	sS := &ServiceSchema{}

	s.DIDMethodVersion = LatestDIDMethodSpec

	for i := range did.DIDKeys {
		sD, err = did.DIDKeys[i].toSchema(did.ID)
		if err != nil {
			return nil, err
		}
		s.DIDKey = append(s.DIDKey, sD)
	}

	for i := range did.ManagementKeys {
		sM, err = did.ManagementKeys[i].toSchema(did.ID)
		if err != nil {
			return nil, err
		}
		s.ManagementKey = append(s.ManagementKey, sM)
	}

	for i := range did.Services {
		sS, err = did.Services[i].toSchema(did.ID)
		if err != nil {
			return nil, err
		}
		s.Service = append(s.Service, sS)
	}

	fe := &factom.Entry{}
	fe.ExtIDs = did.ExtIDs
	fe.Content, err = json.Marshal(s)

	if err != nil {
		return nil, err
	}

	if size := calculateEntrySize(fe); size > MaxEntrySize {
		return nil, fmt.Errorf("You have exceeded the entry size limit")
	}

	return fe, err

}

// Update compares existing and updated DID documents and generates DIDUpdate Factom Entry signed with ManagementKey
func (did *DID) Update(updatedDID *DID, signingKeyAlias string) (*factom.Entry, error) {

	// validate existing DID document
	err := did.Validate()
	if err != nil {
		return nil, err
	}

	// validate updated DID document
	err = updatedDID.Validate()
	if err != nil {
		return nil, err
	}

	// check if it's the same DID document
	if did.ID != updatedDID.ID {
		return nil, fmt.Errorf("ID of updated DID document is not equal to ID of origin DID document")
	}

	// calculate difference between initial and updated DID documents, calculated requiredPriority
	var found bool // equal by alias flag
	var reqPriority = math.MaxInt32

	// empty DID Update Document
	update := &DIDUpdateEntrySchema{}

	// go through all existing ManagementKeys
	for i := range did.ManagementKeys {
		found = false
		// go through all new ManagementKeys
		for j := range updatedDID.ManagementKeys {
			// if aliases are equal
			if did.ManagementKeys[i].Alias == updatedDID.ManagementKeys[j].Alias {
				found = true
				break
			}
		}
		// if not found, then revoke old
		if found == false {
			// revoke old
			r, err := did.ManagementKeys[i].toRevokeIDSchema(did.ID)
			if err != nil {
				return nil, err
			}
			update.Revoke.ManagementKey = append(update.Revoke.ManagementKey, r)
			if did.ManagementKeys[i].PriorityRequirement != nil {
				reqPriority = min(*did.ManagementKeys[i].PriorityRequirement, reqPriority)
			}
		}
	}

	// go through all new ManagementKeys
	for i := range updatedDID.ManagementKeys {
		found = false
		// go through all existing ManagementKeys
		for j := range did.ManagementKeys {
			// if aliases are equal
			if updatedDID.ManagementKeys[i].Alias == did.ManagementKeys[j].Alias {
				found = true
				break
			}
		}
		// if not found, then add new
		if found == false {
			// add new
			a, err := updatedDID.ManagementKeys[i].toSchema(did.ID)
			if err != nil {
				return nil, err
			}
			update.Add.ManagementKey = append(update.Add.ManagementKey, a)
			if updatedDID.ManagementKeys[i].PriorityRequirement != nil {
				reqPriority = min(*updatedDID.ManagementKeys[i].PriorityRequirement, reqPriority)
			}
		}
	}

	// go through all existing DIDKeys
	for i := range did.DIDKeys {
		found = false
		// go through all new DIDKeys
		for j := range updatedDID.DIDKeys {
			// if aliases are equal
			if did.DIDKeys[i].Alias == updatedDID.DIDKeys[j].Alias {
				found = true
				break
			}
		}
		// if not found, then revoke old
		if found == false {
			// revoke old
			r, err := did.DIDKeys[i].toRevokeIDSchema(did.ID)
			if err != nil {
				return nil, err
			}
			update.Revoke.DIDKey = append(update.Revoke.DIDKey, r)
			if did.DIDKeys[i].PriorityRequirement != nil {
				reqPriority = min(*did.DIDKeys[i].PriorityRequirement, reqPriority)
			}
		}
	}

	// go through all new DIDKeys
	for i := range updatedDID.DIDKeys {
		found = false
		// go through all existing DIDKeys
		for j := range did.DIDKeys {
			// if aliases are equal
			if updatedDID.DIDKeys[i].Alias == did.DIDKeys[j].Alias {
				found = true
				break
			}
		}
		// if not found, then add new
		if found == false {
			// add new
			a, err := updatedDID.DIDKeys[i].toSchema(did.ID)
			if err != nil {
				return nil, err
			}
			update.Add.DIDKey = append(update.Add.DIDKey, a)
			if updatedDID.DIDKeys[i].PriorityRequirement != nil {
				reqPriority = min(*updatedDID.DIDKeys[i].PriorityRequirement, reqPriority)
			}
		}
	}

	// go through all existing Services
	for i := range did.Services {
		found = false
		// go through all new Services
		for j := range updatedDID.Services {
			// if aliases are equal
			if did.Services[i].Alias == updatedDID.Services[j].Alias {
				found = true
				break
			}
		}
		// if not found, then revoke old
		if found == false {
			// revoke old
			r, err := did.Services[i].toRevokeIDSchema(did.ID)
			if err != nil {
				return nil, err
			}
			update.Revoke.Service = append(update.Revoke.Service, r)
			if did.Services[i].PriorityRequirement != nil {
				reqPriority = min(*did.Services[i].PriorityRequirement, reqPriority)
			}
		}
	}

	// go through all new Services
	for i := range updatedDID.Services {
		found = false
		// go through all existing Services
		for j := range did.Services {
			// if aliases are equal
			if updatedDID.Services[i].Alias == did.Services[j].Alias {
				found = true
				break
			}
		}
		// if not found, then add new
		if found == false {
			// add new
			a, err := updatedDID.Services[i].toSchema(did.ID)
			if err != nil {
				return nil, err
			}
			update.Add.Service = append(update.Add.Service, a)
			if updatedDID.Services[i].PriorityRequirement != nil {
				reqPriority = min(*updatedDID.Services[i].PriorityRequirement, reqPriority)
			}
		}
	}

	// find ManagementKey
	signingKey := &ManagementKey{}

	if len(did.ManagementKeys) == 0 {
		return nil, fmt.Errorf("No ManagementKeys found in this DID document")
	}

	for _, v := range did.ManagementKeys {
		if v.Alias == signingKeyAlias {
			signingKey = v
		}
	}

	validate := validator.New()
	err = validate.StructPartial(signingKey, "Alias", "PrivateKey")

	if err != nil {
		return nil, err
	}

	// check required priority for signing update
	if signingKey.Priority > reqPriority {
		return nil, fmt.Errorf("The update requires a key with priority <= %d, but the provided signing key priority = %d", reqPriority, signingKey.Priority)
	}

	signingKeyFullID := strings.Join([]string{did.ID, signingKey.Alias}, "#")
	signature, err := signingKey.Sign([]byte(strings.Join([]string{EntryTypeUpdate, LatestEntrySchema, signingKeyFullID}, "")))

	if err != nil {
		return nil, err
	}

	fe := &factom.Entry{}
	fe.ChainID = did.GetChainID()
	fe.ExtIDs = append(fe.ExtIDs, []byte(EntryTypeUpdate))
	fe.ExtIDs = append(fe.ExtIDs, []byte(LatestEntrySchema))
	fe.ExtIDs = append(fe.ExtIDs, []byte(signingKeyFullID))
	fe.ExtIDs = append(fe.ExtIDs, signature)

	fe.Content, err = json.Marshal(update)

	if err != nil {
		return nil, err
	}

	if size := calculateEntrySize(fe); size > MaxEntrySize {
		return nil, fmt.Errorf("You have exceeded the entry size limit")
	}

	return fe, nil

}

// Deactivate generates DIDDeactivation Factom Entry signed with ManagementKey (priority=0 key required)
func (did *DID) Deactivate(signingKeyAlias string) (*factom.Entry, error) {

	// validate existing DID document
	err := did.Validate()
	if err != nil {
		return nil, err
	}

	// find ManagementKey
	signingKey := &ManagementKey{}

	if len(did.ManagementKeys) == 0 {
		return nil, fmt.Errorf("No ManagementKeys found in this DID document")
	}

	for _, v := range did.ManagementKeys {
		if v.Alias == signingKeyAlias {
			signingKey = v
		}
	}

	validate := validator.New()
	err = validate.StructPartial(signingKey, "Alias", "PrivateKey")

	if err != nil {
		return nil, err
	}

	if signingKey.Priority != 0 {
		return nil, fmt.Errorf("You need ManagementKey with 0 priority to deactivate DID")
	}

	signingKeyFullID := strings.Join([]string{did.ID, signingKey.Alias}, "#")
	signature, err := signingKey.Sign([]byte(strings.Join([]string{EntryTypeDeactivation, LatestEntrySchema, signingKeyFullID}, "")))

	if err != nil {
		return nil, err
	}

	fe := &factom.Entry{}
	fe.ExtIDs = append(fe.ExtIDs, []byte(EntryTypeDeactivation))
	fe.ExtIDs = append(fe.ExtIDs, []byte(LatestEntrySchema))
	fe.ExtIDs = append(fe.ExtIDs, []byte(signingKeyFullID))
	fe.ExtIDs = append(fe.ExtIDs, signature)

	if size := calculateEntrySize(fe); size > MaxEntrySize {
		return nil, fmt.Errorf("You have exceeded the entry size limit")
	}

	return fe, nil

}

// Validate validates DID document
func (did *DID) Validate() error {

	var err error
	validate := validator.New()

	// validate DID document (DID)
	err = validate.Struct(did)
	if err != nil {
		return err
	}

	// validate DID ManagementKeys
	if len(did.ManagementKeys) == 0 {
		return fmt.Errorf("DID document must have at least one ManagementKey")
	}
	var hasAtLeastOneZeroPriorityKey bool
	for i := range did.ManagementKeys {
		err = validate.Struct(did.ManagementKeys[i])
		if err != nil {
			return err
		}
		if did.ManagementKeys[i].Priority == 0 {
			hasAtLeastOneZeroPriorityKey = true
		}
	}

	if hasAtLeastOneZeroPriorityKey == false {
		return fmt.Errorf("DID document must have at least one ManagementKey with Priority 0")
	}

	// validate DID DIDKeys
	if len(did.DIDKeys) == 0 {
		return fmt.Errorf("DID document must have at least one DIDKey")
	}
	for i := range did.DIDKeys {
		err = validate.Struct(did.DIDKeys[i])
		if err != nil {
			return err
		}
	}

	// check if DID and Management Keys aliases are unique
	var aliases []string

	for i := range did.ManagementKeys {
		aliases = append(aliases, did.ManagementKeys[i].Alias)
	}
	for j := range did.DIDKeys {
		aliases = append(aliases, did.DIDKeys[j].Alias)
	}

	keys := make(map[string]bool)
	list := []string{}
	for _, item := range aliases {
		if _, value := keys[item]; !value {
			keys[item] = true
			list = append(list, item)
		}
	}

	if len(list) != len(aliases) {
		return fmt.Errorf("All keys aliases must be unique, 2 or more aliases of []*DIDKey and []*ManagementKey are the same")
	}

	// check if services aliases are unique
	var sAliases []string

	for i := range did.Services {
		sAliases = append(sAliases, did.Services[i].Alias)
	}

	services := make(map[string]bool)
	slist := []string{}
	for _, item := range sAliases {
		if _, value := services[item]; !value {
			services[item] = true
			slist = append(slist, item)
		}
	}

	if len(slist) != len(sAliases) {
		return fmt.Errorf("All services aliases must be unique, 2 or more aliases of []*Services are the same")
	}

	return nil

}

// Copy makes a copy of DID Document for update
func (did *DID) Copy() *DID {

	copy := &DID{}
	copy.ID = did.ID
	copy.Network = did.Network
	copy.ExtIDs = did.ExtIDs

	for i := range did.ManagementKeys {
		tmp := *did.ManagementKeys[i]
		copy.ManagementKeys = append(copy.ManagementKeys, &tmp)
	}

	for i := range did.DIDKeys {
		tmp := *did.DIDKeys[i]
		copy.DIDKeys = append(copy.DIDKeys, &tmp)
	}

	for i := range did.Services {
		tmp := *did.Services[i]
		copy.Services = append(copy.Services, &tmp)
	}

	return copy

}
