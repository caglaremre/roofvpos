package models

type PointRequest struct {
	OrderID        string `json:"orderId" form:"point-order-id"`
	CardNo         string `json:"cardNo" form:"point-card-no"`
	Expiry         int    `json:"expiry" form:"point-expiry"`
	Lang           string `json:"lang,omitempty" form:"point-lang"`
	AdditionalInfo string `json:"additionalInfo,omitempty" form:"point-additional-info"`
	InvoiceInfo    string `json:"invoiceInfo,omitempty" form:"point-invoice-info"`
}

type PointResponse struct {
	ResultCode           string `json:"resultCode"`
	ProcReturnCode       string `json:"procReturnCode"`
	ResultMessage        string `json:"resultMessage"`
	AvailablePoint       int    `json:"availablePoint"`
	PointCode            string `json:"pointCode"`
	PointMultiplier      int    `json:"pointMultiplier"`
	AvailablePointAmount int    `json:"availablePointAmount"`
}
