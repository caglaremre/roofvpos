package models

import "time"

type SaleRequest struct {
	OrderID              string `json:"orderId" form:"sale-order-id"`
	CardNo               string `json:"cardNo" form:"sale-card-no"`
	Expiry               int    `json:"expiry" form:"sale-expiry"`
	CVV2                 string `json:"cvv2" form:"sale-cvv"`
	Amount               int    `json:"amount" form:"sale-amount"`
	PointAmount          int    `json:"pointAmount,omitempty" form:"sale-point"`
	Currency             int    `json:"currency" form:"sale-currency"`
	EcommerceTxnType     int    `json:"ecommerceTxnType,omitempty" form:"sale-ecommerce-txn-type"`
	TxnInitiationType    int    `json:"txnInitiationType,omitempty" form:"sale-txn-initiation-type"`
	InstallmentCount     int    `json:"installmentCount,omitempty" form:"sale-installment-count"`
	CardHolderIp         string `json:"cardHolderIp,omitempty" form:"sale-cardholder-ip"`
	MerchantIp           string `json:"merchantIp,omitempty" form:"sale-merchant-ip"`
	SubmerchantIp        string `json:"submerchantIp,omitempty" form:"sale-submerchant-ip"`
	Lang                 string `json:"lang,omitempty" form:"sale-lang"`
	AdditionalInfo       string `json:"additionalInfo,omitempty" form:"sale-additional-info"`
	InvoiceInfo          string `json:"invoiceInfo,omitempty" form:"sale-invoice-info"`
	IdentificationNumber string `json:"identificationNumber,omitempty" form:"sale-identificationNumber"`
}

type SaleResponse struct {
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
