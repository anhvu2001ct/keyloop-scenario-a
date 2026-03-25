package model_test

import (
	"scenario-a/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomer_TableName(t *testing.T) {
	c := &model.Customer{}
	assert.Equal(t, "customers", c.TableName())
}

func TestServiceBay_TableName(t *testing.T) {
	s := &model.ServiceBay{}
	assert.Equal(t, "service_bays", s.TableName())
}

func TestServiceType_TableName(t *testing.T) {
	s := &model.ServiceType{}
	assert.Equal(t, "service_types", s.TableName())
}

func TestTechnician_TableName(t *testing.T) {
	tech := &model.Technician{}
	assert.Equal(t, "technicians", tech.TableName())
}

func TestTechnicianServiceType_TableName(t *testing.T) {
	techType := model.TechnicianServiceType{}
	assert.Equal(t, "technician_service_types", techType.TableName())
}

func TestVehicle_TableName(t *testing.T) {
	v := &model.Vehicle{}
	assert.Equal(t, "vehicles", v.TableName())
}
