package models

import "time"

type CompletePaymentRequest struct {
	ProcessID      string `json:"processId" form:"completepayment-process-id"`
	OrderID        string `json:"orderId" form:"completepayment-order-id"`
	Lang           string `json:"lang" form:"completepayment-lang"`
	AdditionalInfo string `json:"additionalInfo" form:"completepayment-additional-info"`
	InvoiceInfo    string `json:"invoiceInfo" form:"completepayment-invoice-info"`
}

type CompletePaymentResponse struct {
	RequestType          string    `json:"requestType"`
	ProcessId            string    `json:"processId"`
	OrderId              string    `json:"orderId"`
	ResultMessage        string    `json:"resultMessage"`
	ResultCode           string    `json:"resultCode"`
	ProcReturnCode       string    `json:"procReturnCode"`
	AuthCode             string    `json:"authCode"`
	SecureType           string    `json:"secureType"`
	TxnType              string    `json:"txnType"`
	CardMask             string    `json:"cardMask"`
	Amount               int       `json:"amount"`
	PointAmount          int       `json:"pointAmount"`
	PointCode            string    `json:"pointCode"`
	PointMultiplier      int       `json:"pointMultiplier"`
	AvailablePointAmount string    `json:"availablePointAmount"`
	InstallmentCount     string    `json:"installmentCount"`
	InstallmentIndex     string    `json:"installmentIndex"`
	MerchantId           string    `json:"merchantId"`
	TerminalId           string    `json:"terminalId"`
	BatchNo              int       `json:"batchNo"`
	ProcessDate          time.Time `json:"processDate"`
	Rrn                  string    `json:"rrn"`
	AdditionalInfo       string    `json:"additionalInfo"`
	InvoiceInfo          string    `json:"invoiceInfo"`
	EcommerceTxnType     int       `json:"ecommerceTxnType"`
	TxnInitiationType    int       `json:"txnInitiationType"`
}
