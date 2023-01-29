package values

type TransactionStatus string

const (
	TransactionStatusOpen   TransactionStatus = "OPEN"
	TransactionStatusBooked TransactionStatus = "BOOKED"
	TransactionStatusFailed TransactionStatus = "FAILED"
)
