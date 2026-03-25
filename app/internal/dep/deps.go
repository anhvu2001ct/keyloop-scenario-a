package dep

import (
	"scenario-a/internal/repository"
	"scenario-a/internal/service"

	"gorm.io/gorm"
)

func Init(gormDB *gorm.DB) *Dependencies {
	// ---- Repository ---- //
	appointmentRepo := repository.NewAppointment(gormDB)
	serviceTypeRepo := repository.NewServiceType(gormDB)
	dealershipRepo := repository.NewDealership(gormDB)
	technicianRepo := repository.NewTechnician(gormDB)
	serviceBayRepo := repository.NewServiceBay(gormDB)
	customerRepo := repository.NewCustomer(gormDB)

	// ---- Service ---- //
	appointmentService := service.NewAppointment(
		appointmentRepo,
		serviceTypeRepo,
		dealershipRepo,
		technicianRepo,
		serviceBayRepo,
		customerRepo,
	)

	res := &Dependencies{
		Repository: repositories{
			Appointment: appointmentRepo,
			ServiceType: serviceTypeRepo,
			Dealership:  dealershipRepo,
			Technician:  technicianRepo,
			Customer:    customerRepo,
		},
		Service: services{
			Appointment: appointmentService,
		},
	}

	return res
}
