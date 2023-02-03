package transactionStatusApi

type SearchTransactionStatusRequest struct {
	TransactionID string `form:"transaction_id" binding:"required"`
}
