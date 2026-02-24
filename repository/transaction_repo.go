package repository

import (
	"cmp"
	"encoding/json"
	"net/http"
	"roof/vpos/models"
	"slices"
	"strings"
	"time"

	"go.etcd.io/bbolt"
)

type TransactionRepository struct {
	DB *bbolt.DB
}

func (t *TransactionRepository) GetAllTransactions() []models.Transaction {
	var transactions []models.Transaction
	err := t.DB.View(func(tx *bbolt.Tx) error {
		transactionsBucket := tx.Bucket([]byte("transactions"))
		transactionsCursor := transactionsBucket.Cursor()
		for orderId, _ := transactionsCursor.First(); orderId != nil; orderId, _ = transactionsCursor.Next() {
			var transaction models.Transaction

			orderIdBucket := transactionsBucket.Bucket(orderId)
			transaction.OrderID = string(orderId)

			logDateByte := orderIdBucket.Get([]byte("logDate"))
			transaction.LogDate = string(logDateByte)
			getTransactionDetails(&transaction, orderIdBucket)
			transactions = append(transactions, transaction)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
	slices.SortFunc(transactions, func(a, b models.Transaction) int {
		return cmp.Compare(b.LogDate, a.LogDate)
	})
	return transactions
}

func getTransactionDetails(transaction *models.Transaction, orderBucket *bbolt.Bucket) {
	orderBucket.ForEachBucket(func(transactionType []byte) error {
		actionBucket := orderBucket.Bucket([]byte(transactionType))
		if actionBucket != nil {
			requestBucket := actionBucket.Bucket([]byte("request"))
			requestHeaders := http.Header{}
			_ = json.Unmarshal(requestBucket.Get([]byte("headers")), &requestHeaders)

			responseBucket := actionBucket.Bucket([]byte("response"))
			responseHeaders := http.Header{}
			_ = json.Unmarshal(responseBucket.Get([]byte("headers")), &responseHeaders)

			for key := range responseHeaders {
				if !strings.HasPrefix(key, "X") {
					responseHeaders.Del(key)
				}
			}

			switch string(transactionType) {
			case "sale", "presale":
				transaction.SaleRequestHeaders = requestHeaders
				transaction.SaleResponseHeaders = responseHeaders

				requestBody := models.SaleRequest{}
				_ = json.Unmarshal(requestBucket.Get([]byte("body")), &requestBody)
				transaction.SaleRequest = requestBody

				responseBody := models.SaleResponse{}
				_ = json.Unmarshal(responseBucket.Get([]byte("body")), &responseBody)
				transaction.SaleResponse = responseBody

			case "postsale":
				transaction.PostSaleRequestHeaders = requestHeaders
				transaction.PostSaleResponseHeaders = responseHeaders

				requestBody := models.PostSaleRequest{}
				_ = json.Unmarshal(requestBucket.Get([]byte("body")), &requestBody)
				transaction.PostSaleRequest = requestBody

				responseBody := models.PostSaleResponse{}
				_ = json.Unmarshal(responseBucket.Get([]byte("body")), &responseBody)
				transaction.PostSaleResponse = responseBody

			case "void":
				transaction.VoidRequestHeaders = requestHeaders
				transaction.VoidResponseHeaders = responseHeaders

				requestBody := models.VoidRequest{}
				_ = json.Unmarshal(requestBucket.Get([]byte("body")), &requestBody)
				transaction.VoidRequest = requestBody

				responseBody := models.VoidResponse{}
				_ = json.Unmarshal(responseBucket.Get([]byte("body")), &responseBody)
				transaction.VoidResponse = responseBody

			case "refund":
				transaction.RefundRequestHeaders = requestHeaders
				transaction.RefundResponseHeaders = responseHeaders

				requestBody := models.RefundRequest{}
				_ = json.Unmarshal(requestBucket.Get([]byte("body")), &requestBody)
				transaction.RefundRequest = requestBody

				responseBody := models.RefundResponse{}
				_ = json.Unmarshal(responseBucket.Get([]byte("body")), &responseBody)
				transaction.RefundResponse = responseBody
			case "point":
				transaction.PointRequestHeaders = requestHeaders
				transaction.PointResponseHeaders = responseHeaders

				requestBody := models.PointRequest{}
				_ = json.Unmarshal(requestBucket.Get([]byte("body")), &requestBody)
				transaction.PointRequest = requestBody

				responseBody := models.PointResponse{}
				_ = json.Unmarshal(responseBucket.Get([]byte("body")), &responseBody)
				transaction.PointResponse = responseBody

			case "threeds":
				transaction.ThreeDSRequestHeaders = requestHeaders
				transaction.ThreeDSResponseHeaders = responseHeaders

				requestBody := models.ThreeDSRequest{}
				_ = json.Unmarshal(requestBucket.Get([]byte("body")), &requestBody)
				transaction.ThreeDSRequest = requestBody

				responseBody := models.ThreeDSResponse{}
				_ = json.Unmarshal(responseBucket.Get([]byte("body")), &responseBody)
				transaction.ThreeDSResponse = responseBody

			case "token":
				transaction.TokenRequestHeaders = requestHeaders
				transaction.TokenResponseHeaders = responseHeaders

				requestBody := models.TokenRequest{}
				_ = json.Unmarshal(requestBucket.Get([]byte("body")), &requestBody)
				transaction.TokenRequest = requestBody

				responseBody := models.TokenResponse{}
				_ = json.Unmarshal(responseBucket.Get([]byte("body")), &responseBody)
				transaction.TokenResponse = responseBody

			case "completepayment":
				transaction.CompletePaymentRequestHeaders = requestHeaders
				transaction.CompletePaymentResponseHeaders = responseHeaders

				requestBody := models.CompletePaymentRequest{}
				_ = json.Unmarshal(requestBucket.Get([]byte("body")), &requestBody)
				transaction.CompletePaymentRequest = requestBody

				responseBody := models.CompletePaymentResponse{}
				_ = json.Unmarshal(responseBucket.Get([]byte("body")), &responseBody)
				transaction.CompletePaymentResponse = responseBody

			case "threedshosting":
				transaction.ThreeDSHostingRequestHeaders = requestHeaders
				transaction.ThreeDSHostingResponseHeaders = responseHeaders

				requestBody := models.ThreeDSHostingRequest{}
				_ = json.Unmarshal(requestBucket.Get([]byte("body")), &requestBody)
				transaction.ThreeDSHostingRequest = requestBody

				responseBody := models.ThreeDSHostingResponse{}
				_ = json.Unmarshal(responseBucket.Get([]byte("body")), &responseBody)
				transaction.ThreeDSHostingResponse = responseBody

			}

		}

		return nil
	})
}

func (t *TransactionRepository) LogRequest(transactionType, action, orderID string, body []byte, headers http.Header) error {
	err := t.DB.Update(func(tx *bbolt.Tx) error {
		transactionsBucket := tx.Bucket([]byte("transactions"))

		orderIDBucket, _ := transactionsBucket.CreateBucketIfNotExists([]byte(orderID))

		logDate := time.Now().Format(time.RFC3339)
		err := orderIDBucket.Put([]byte("logDate"), []byte(logDate))
		if err != nil {
			return err
		}

		headerBytes, err := json.Marshal(headers)
		if err != nil {
			return err
		}

		transactionTypeBucket, _ := orderIDBucket.CreateBucketIfNotExists([]byte(transactionType))
		actionBucket, _ := transactionTypeBucket.CreateBucketIfNotExists([]byte(action))
		_ = actionBucket.Put([]byte("headers"), headerBytes)
		_ = actionBucket.Put([]byte("body"), body)
		return nil
	})
	return err
}
