package models

type Transaction struct {
	OrderID string
	LogDate string

	SaleRequest         SaleRequest
	SaleRequestBody     string
	SaleRequestHeaders  string
	SaleResponse        SaleResponse
	SaleResponseBody    string
	SaleResponseHeaders string

	VoidRequest         VoidRequest
	VoidRequestBody     string
	VoidRequestHeaders  string
	VoidResponse        VoidResponse
	VoidResponseBody    string
	VoidResponseHeaders string
}
