package models

type ThreeDSRequest struct {
	OrderId           string `json:"orderId" form:"threeds-order-id"`
	CardNo            string `json:"cardNo" form:"threeds-card-no"`
	Expiry            int    `json:"expiry" form:"threeds-expiry"`
	Cvv2              string `json:"cvv2,omitempty" form:"threeds-cvv"`
	Amount            int    `json:"amount" form:"threeds-amount"`
	PointAmount       int    `json:"pointAmount,omitempty" form:"threeds-point"`
	Currency          int    `json:"currency" form:"threeds-currency"`
	InstallmentCount  int    `json:"installmentCount,omitempty" form:"threeds-installment-count"`
	RequestType       string `json:"requestType" form:"threeds-request-type"`
	CardHolderIp      string `json:"cardHolderIp,omitempty" form:"threeds-cardholder-ip"`
	MerchantIp        string `json:"merchantIp,omitempty" form:"threeds-merchant-ip"`
	SubMerchantIp     string `json:"subMerchantIp,omitempty" form:"threeds-submerchant-ip"`
	TxnInitiationType int    `json:"txnInitiationType,omitempty" form:"threeds-txn-initiation-type"`
	Lang              string `json:"lang,omitempty" form:"threeds-lang"`
	ReturnUrl         string `json:"returnUrl,omitempty" form:"threeds-return-url"`
	CallbackUrl       string `json:"callbackUrl,omitempty" form:"threeds-callback-url"`
	CardholderName    string `json:"cardholderName,omitempty" form:"threeds-cardholder-name"`
	CardholderEmail   string `json:"cardholderEmail,omitempty" form:"threeds-cardholder-email"`
	AdditionalInfo    string `json:"additionalInfo,omitempty" form:"threeds-additional-info"`
	InvoiceInfo       string `json:"invoiceInfo,omitempty" form:"threeds-invoice-info"`
}

type ThreeDSResponse struct {
	ResultCode     string `json:"resultCode" form:""`
	ProcReturnCode string `json:"procReturnCode" form:""`
	ResultMessage  string `json:"resultMessage" form:""`
	ReturnUrl      string `json:"returnUrl" form:""`
	ProcessId      string `json:"processId" form:""`
	Token          string `json:"token" form:""`
	HtmlContent    string `json:"htmlContent" form:""`
}
