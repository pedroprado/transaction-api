package entity

type BalanceProvision struct {
	ProvisionID    string
	Value          float64
	Origin         string
	Destination    string
	IsCompensation bool
}
