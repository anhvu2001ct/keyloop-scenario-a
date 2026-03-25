package model

type ServiceType struct {
	CommonModel

	Name            string
	DurationMinutes int
}

func (*ServiceType) TableName() string {
	return "service_types"
}
