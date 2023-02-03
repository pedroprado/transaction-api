package balanceProvisionApi

type SearchBalanceProvisionsRequest struct {
	TransactionID string `form:"transaction_id" binding:"required"`
}
