package models

import (
	"net/http"
)

type Transaction struct {
	OrderID             string
	LogDate             string
	SaleRequest         SaleRequest
	SaleRequestBody     string
	SaleRequestsHeaders http.Header
	SaleResponse        SaleResponse
	SaleResponseBody    string
	SaleResponseHeaders http.Header
	SaleProcessID       string
}
