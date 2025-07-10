package types

// TxResponse handles the transaction response in the interface ConsumerController
// Not every consumer has Events thing in their response,
// so consumer client implementations need to care about Events field.
type TxResponse struct {
	TxHash string
	Events []byte
}
