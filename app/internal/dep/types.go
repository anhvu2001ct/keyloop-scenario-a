package dep

import (
	"scenario-a/internal/repository"
	"scenario-a/internal/service"
)

type repositories struct {
	Appointment repository.Appointment
	Dealership  repository.Dealership
	ServiceType repository.ServiceType
	Technician  repository.Technician
	Customer    repository.Customer
}

type services struct {
	Appointment *service.Appointment
}

// Dependencies holds all the dependencies of the application
type Dependencies struct {
	Repository repositories
	Service    services
}
