package factomdid

type DIDManagementEntrySchema struct {
	DIDKey           []*DIDKeySchema        `json:"didKey" form:"didKey" query:"didKey"`
	DIDMethodVersion string                 `json:"didMethodVersion" form:"didMethodVersion" query:"didMethodVersion"`
	ManagementKey    []*ManagementKeySchema `json:"managementKey" form:"managementKey" query:"managementKey"`
	Service          []*ServiceSchema       `json:"service,omitempty" form:"service" query:"service"`
}

type DIDUpdateEntrySchema struct {
	Add struct {
		DIDKey        []*DIDKeySchema        `json:"didKey,omitempty" form:"didKey" query:"didKey"`
		ManagementKey []*ManagementKeySchema `json:"managementKey,omitempty" form:"managementKey" query:"managementKey"`
		Service       []*ServiceSchema       `json:"service,omitempty" form:"service" query:"service"`
	}
	Revoke struct {
		DIDKey        []*RevokeIDSchema `json:"didKey,omitempty" form:"didKey" query:"didKey"`
		ManagementKey []*RevokeIDSchema `json:"managementKey,omitempty" form:"managementKey" query:"managementKey"`
		Service       []*RevokeIDSchema `json:"service,omitempty" form:"service" query:"service"`
	}
}

type DIDDeactivationEntrySchema struct{}

type DIDKeySchema struct {
	Controller          string   `json:"controller" form:"controller" query:"controller"`
	ID                  string   `json:"id" form:"id" query:"id"`
	PublicKeyBase58     string   `json:"publicKeyBase58,omitempty" form:"publicKeyBase58" query:"publicKeyBase58"`
	PublicKeyPem        string   `json:"publicKeyPem,omitempty" form:"publicKeyPem" query:"publicKeyPem"`
	Purpose             []string `json:"purpose" form:"purpose" query:"purpose"`
	Type                string   `json:"type" form:"type" query:"type"`
	PriorityRequirement *int     `json:"priorityRequirement,omitempty" form:"priorityRequirement" query:"priorityRequirement"`
	BIP44               string   `json:"bip44,omitempty" form:"bip44" query:"bip44"`
}

type ManagementKeySchema struct {
	Controller          string `json:"controller" form:"controller" query:"controller"`
	ID                  string `json:"id" form:"id" query:"id"`
	PublicKeyBase58     string `json:"publicKeyBase58,omitempty" form:"publicKeyBase58" query:"publicKeyBase58"`
	PublicKeyPem        string `json:"publicKeyPem,omitempty" form:"publicKeyPem" query:"publicKeyPem"`
	Type                string `json:"type" form:"type" query:"type"`
	Priority            int    `json:"priority" form:"priority" query:"priority"`
	PriorityRequirement *int   `json:"priorityRequirement,omitempty" form:"priorityRequirement" query:"priorityRequirement"`
	BIP44               string `json:"bip44,omitempty" form:"bip44" query:"bip44"`
}

type ServiceSchema struct {
	ID                  string `json:"id" form:"id" query:"id"`
	Type                string `json:"type" form:"type" query:"type"`
	ServiceEndpoint     string `json:"serviceEndpoint" form:"serviceEndpoint" query:"serviceEndpoint"`
	PriorityRequirement *int   `json:"priorityRequirement,omitempty" form:"priorityRequirement" query:"priorityRequirement"`
}

type RevokeIDSchema struct {
	ID string `json:"id" form:"id" query:"id"`
}
