package model

type Customer struct {
	CommonModel

	Name  string
	Email string
	Phone string
}

func (*Customer) TableName() string {
	return "customers"
}
