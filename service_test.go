package factomdid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {

	// valid Services
	s1, err := NewService("kyc", "KYC", "https://kyc.example.com")
	assert.NotNil(t, s1)
	assert.NoError(t, err)

	s2, err := NewService("kyc", "KYC", "http://kyc.example.com")
	assert.NotNil(t, s2)
	assert.NoError(t, err)

	// invalid Services
	s3, err := NewService("", "KYC", "https://kyc.example.com")
	assert.Nil(t, s3)
	assert.Error(t, err)

	s4, err := NewService("kyc", "", "https://kyc.example.com")
	assert.Nil(t, s4)
	assert.Error(t, err)

	s5, err := NewService("kyc", "KYC", "")
	assert.Nil(t, s5)
	assert.Error(t, err)

	s6, err := NewService("kyc", "KYC", "example.com")
	assert.Nil(t, s6)
	assert.Error(t, err)

}

func TestServiceToSchema(t *testing.T) {

	did := "did:factom:301a57c2e753d061928cf6b6a692ea052885d75d2af5640e9b5cbc8897bbf7d5"

	s, _ := NewService("kyc", "KYC", "https://kyc.example.com")

	// test priorityRequirement 0 not omited
	var i *int
	i = new(int)
	*i = 0
	s.PriorityRequirement = i

	s2, err := s.toSchema(did)
	assert.NotEmpty(t, s2)
	assert.NoError(t, err)

}
