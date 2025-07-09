package types

// TxResponse handles the transaction response for client operations.
// Not every response has Events, so client implementations need to handle Events field appropriately.
type TxResponse struct {
	TxHash string
	Events []byte
}
