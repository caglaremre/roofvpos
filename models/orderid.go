package models

import "time"

type CheckOrderRequest struct {
	OrderId string `json:"orderId" form:"checkorder-orderid"`
	Lang    string `json:"lang" form:"checkorder-lang"`
}

type CheckOrderResponse struct {
	Eci                  string    `json:"eci"`
	Cavv                 string    `json:"cavv"`
	TxnStatus            string    `json:"txnStatus"`
	TerminalAuthType     int       `json:"terminalAuthType"`
	RequestType          string    `json:"requestType"`
	ReturnUrl            string    `json:"returnUrl"`
	PostSaleAmount       int       `json:"postSaleAmount"`
	AmountDci            string    `json:"amountDci"`
	BatchStatus          string    `json:"batchStatus"`
	QrStatus             string    `json:"qrStatus"`
	AdditionalContent    string    `json:"additionalContent"`
	ProcessExpireDate    string    `json:"processExpireDate"`
	LinkStatus           string    `json:"linkStatus"`
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
