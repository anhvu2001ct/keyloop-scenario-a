package model

type Technician struct {
	CommonModel

	DealershipID int64
	Name         string
}

func (*Technician) TableName() string {
	return "technicians"
}
