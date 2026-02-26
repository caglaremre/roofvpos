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

	ThreeDSRequest         ThreeDSRequest
	ThreeDSRequestHeaders  http.Header
	ThreeDSResponse        ThreeDSResponse
	ThreeDSResponseHeaders http.Header

	TokenRequest         TokenRequest
	TokenRequestHeaders  http.Header
	TokenResponse        TokenResponse
	TokenResponseHeaders http.Header

	CompletePaymentRequest         CompletePaymentRequest
	CompletePaymentRequestHeaders  http.Header
	CompletePaymentResponse        CompletePaymentResponse
	CompletePaymentResponseHeaders http.Header

	ThreeDSHostingRequest         ThreeDSHostingRequest
	ThreeDSHostingRequestHeaders  http.Header
	ThreeDSHostingResponse        ThreeDSHostingResponse
	ThreeDSHostingResponseHeaders http.Header

	CheckOrderRequest         CheckOrderRequest
	CheckOrderRequestHeaders  http.Header
	CheckOrderResponse        CheckOrderResponse
	CheckOrderResponseHeaders http.Header
}
