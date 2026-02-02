package models

import "time"

type PostSaleRequest struct {
	ProcessID      string `json:"processId" form:"postsale-process-id"`
	OrderID        string `json:"orderId" form:"postsale-order-id"`
	Lang           string `json:"lang" form:"postsale-lang"`
	Amount         int    `json:"amount" form:"postsale-amount"`
	AdditionalInfo string `json:"additionalInfo" form:"postsale-additional-info"`
	InvoiceInfo    string `json:"invoiceInfo" form:"postsale-invoice-info"`
}

type PostSaleResponse struct {
	ProcessID         string    `json:"processId"`
	OrderID           string    `json:"orderId"`
	ResultMessage     string    `json:"resultMessage"`
	ResultCode        string    `json:"resultCode"`
	ProcReturnCode    string    `json:"procReturnCode"`
	AuthCode          string    `json:"authCode"`
	SecureType        string    `json:"secureType"`
	TxnType           string    `json:"txnType"`
	CardMask          string    `json:"cardMask"`
	Amount            int       `json:"amount"`
	PointAmount       int       `json:"pointAmount"`
	InstallmentCount  int       `json:"installmentCount"`
	MerchantID        string    `json:"merchantId"`
	TerminalID        string    `json:"terminalId"`
	BatchNo           int       `json:"batchNo"`
	ProcessDate       time.Time `json:"processDate"`
	RRN               string    `json:"rrn"`
	AdditionalInfo    string    `json:"additionalInfo"`
	InvoiceInfo       string    `json:"invoiceInfo"`
	EcommerceTxnType  int       `json:"ecommerceTxnType"`
	TxnInitiationType int       `json:"txnInitiationType"`
}
