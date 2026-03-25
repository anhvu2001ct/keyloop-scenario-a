package model

type Vehicle struct {
	CommonModel

	CustomerID *int64
	Name       string
}

func (*Vehicle) TableName() string {
	return "vehicles"
}
