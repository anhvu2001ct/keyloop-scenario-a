package model

// this struct exists to keep track table name only
type TechnicianServiceType struct {
}

func (TechnicianServiceType) TableName() string {
	return "technician_service_types"
}
