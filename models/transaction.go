package models

import "net/http"

type Transaction struct {
	OrderID string
	LogDate string

	SaleRequest         SaleRequest
	SaleRequestHeaders  http.Header
	SaleResponse        SaleResponse
	SaleResponseHeaders http.Header

	PostSaleRequest         PostSaleRequest
	PostSaleRequestHeaders  http.Header
	PostSaleResponse        PostSaleResponse
	PostSaleResponseHeaders http.Header

	VoidRequest         VoidRequest
	VoidRequestHeaders  http.Header
	VoidResponse        VoidResponse
	VoidResponseHeaders http.Header

	RefundRequest         RefundRequest
	RefundRequestHeaders  http.Header
	RefundResponse        RefundResponse
	RefundResponseHeaders http.Header

	PointRequest         PointRequest
	PointRequestHeaders  http.Header
	PointResponse        PointResponse
	PointResponseHeaders http.Header
}
