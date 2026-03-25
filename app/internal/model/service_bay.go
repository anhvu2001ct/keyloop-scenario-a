package model

type ServiceBay struct {
	CommonModel

	DealershipID int64
	Name         string
}

func (*ServiceBay) TableName() string {
	return "service_bays"
}
