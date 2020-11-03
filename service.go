package factomdid

import (
	"strings"
)

// Service represents a service associated with a DID.
// A service is an end-point, which can be used to communicate with the DID or to carry out different tasks on behalf of the DID (such as signatures, e.g.)
type Service struct {
	Alias               string `json:"alias" form:"alias" query:"alias" validate:"required"`
	ServiceType         string `json:"serviceType" form:"serviceType" query:"serviceType" validate:"required"`
	Endpoint            string `json:"endpoint" form:"endpoint" query:"endpoint" validate:"required,url"`
	PriorityRequirement *int   `json:"priorityRequirement" form:"priorityRequirement" query:"priorityRequirement" validate:"omitempty,min=0"`
	CustomField         []byte `json:"customFields" form:"customFields" query:"customFields"`
}

// NewService creates new Service
func NewService(alias string, serviceType string, endpoint string) (*Service, error) {

	service := &Service{}
	service.Alias = alias
	service.ServiceType = serviceType
	service.Endpoint = endpoint

	// validate
	err := validate.Struct(service)
	if err != nil {
		return nil, err
	}

	return service, nil

}

// helper function to convert Service into ServiceKeySchema
func (service *Service) toSchema(DID string) (*ServiceSchema, error) {

	// validate DID Key
	err := validate.Struct(service)
	if err != nil {
		return nil, err
	}

	s := &ServiceSchema{}
	s.ID = strings.Join([]string{DID, service.Alias}, "#")
	s.PriorityRequirement = service.PriorityRequirement
	s.Type = service.ServiceType
	s.ServiceEndpoint = service.Endpoint

	return s, nil

}

// helper function to convert Service into RevokeIDSchema
func (service *Service) toRevokeIDSchema(DID string) (*RevokeIDSchema, error) {

	// validate Service
	err := validate.StructPartial(service, "Alias")
	if err != nil {
		return nil, err
	}

	s := &RevokeIDSchema{}
	s.ID = strings.Join([]string{DID, service.Alias}, "#")

	return s, nil

}

// SetPriorityRequirement sets PriorityRequirement for Service
func (service *Service) SetPriorityRequirement(i int) *Service {
	service.PriorityRequirement = &i
	return service
}
