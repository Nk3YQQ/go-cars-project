package types

type Fuel string

type Transmission string

const (
	Petrol      Fuel = "бензин"
	Diesel      Fuel = "дизель"
	Electricity Fuel = "электричество"
	Hybrid      Fuel = "гибрид"
)

const (
	Mechanical Transmission = "механическая"
	Automatic  Transmission = "автоматическая"
	Variator   Transmission = "вариатор"
	Robot      Transmission = "робот"
)
