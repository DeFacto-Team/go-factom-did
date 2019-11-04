package model

import "github.com/FactomProject/factom"

// Service represents a service associated with a DID
// A service is an end-point, which can be used to communicate with the DID or to carry out different tasks on behalf of the DID (such as signatures, e.g.)
type Service struct {
	Alias               string `json:"alias" form:"alias" query:"alias" validate:"required"`
	ServiceType         string `json:"serviceType" form:"serviceType" query:"serviceType" validate:"required"`
	Endpoint            string `json:"endpoint" form:"endpoint" query:"endpoint" validate:"required"`
	PriorityRequirement int    `json:"priorityRequirement" form:"priorityRequirement" query:"priorityRequirement"`
}

// FullID returns full id for the service, constituting of the DID_METHOD_NAME, the controller and the service alias.
func (service *Service) FullID() (string, error) {

	// TBD

	return "", nil

}

// ConvertToFactomEntry converts Service struct to Factom Entry
func (service *Service) ConvertToFactomEntry() *factom.Entry {

	fe := &factom.Entry{}

	// TBD

	return fe

}
