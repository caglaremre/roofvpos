package models

import (
	"encoding/json"
)

type ThreeDSHostingRequest struct {
	OrderId          string `json:"orderId" form:"threedshosting-order-id"`
	OrderType        string `json:"orderType" form:"threedshosting-order-type"`
	RequestType      string `json:"requestType" form:"threedshosting-request-type"`
	TxnType          string `json:"txnType" form:"threedshosting-transaction-type"`
	Amount           int    `json:"amount" form:"threedshosting-amount"`
	Currency         int    `json:"currency" form:"threedshosting-currency"`
	InstallmentCount int    `json:"installmentCount,omitempty" form:"threedshosting-installment-count"`
	Lang             string `json:"lang,omitempty" form:"threedshosting-lang"`
	AdditionalInfo   string `json:"additionalInfo,omitempty" form:"threedshosting-additional-info"`
	InvoiceInfo      string `json:"invoiceInfo,omitempty" form:"threedshosting-invoice-info"`
	ReturnUrl        string `json:"returnUrl,omitempty" form:"threedshosting-return-url"`
	CallbackUrl      string `json:"callbackUrl,omitempty" form:"threedshosting-callback-url"`
}

type ThreeDSHostingResponse struct {
	ProcessId     string `json:"processId"`
	ResultMessage string `json:"resultMessage"`
	ResultCode    string `json:"resultCode"`
	Link          string `json:"link"`
}

func (req *ThreeDSHostingRequest) ToJson(indent bool) string {
	threedshostingreqJson := jsonize(req, indent)
	return string(threedshostingreqJson)
}

func (res *ThreeDSHostingResponse) ToJson(indent bool) string {
	threedshostingreqJson := jsonize(res, indent)
	return string(threedshostingreqJson)
}

func jsonize(v any, indent bool) []byte {
	if indent {
		req, _ := json.MarshalIndent(v, "", " ")
		return req
	}
	req, _ := json.Marshal(v)
	return req
}
