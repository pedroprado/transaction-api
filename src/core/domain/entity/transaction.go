package entity

type Transaction struct {
	TransactionID      string
	Type               string
	OriginAccount      string
	DestinationAccount string
	Value              float64
	Status             string
	ProvisionID        string
}
