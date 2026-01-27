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

func getTransactionDetails(transaction *models.Transaction, orderIdBucket *bbolt.Bucket) {

	transactionTypeList := []string{"sale", "void", "refund"}
	for _, action := range transactionTypeList {
		actionBucket := orderIdBucket.Bucket([]byte(action))
		if actionBucket != nil {

			requestBucket := actionBucket.Bucket([]byte("request"))
			requestHeaders := http.Header{}
			_ = json.Unmarshal(requestBucket.Get([]byte("headers")), &requestHeaders)

			responseBucket := actionBucket.Bucket([]byte("response"))
			responseHeaders := http.Header{}
			_ = json.Unmarshal(responseBucket.Get([]byte("headers")), &responseHeaders)

			switch action {
			case "sale":
				transaction.SaleRequestHeaders = requestHeaders
				transaction.SaleResponseHeaders = responseHeaders
				for key, _ := range responseHeaders {
					if !strings.HasPrefix(key, "X") {
						responseHeaders.Del(key)
					}
				}

				saleRequestBody := models.SaleRequest{}
				_ = json.Unmarshal(requestBucket.Get([]byte("body")), &saleRequestBody)
				transaction.SaleRequest = saleRequestBody

				saleResponseBody := models.SaleResponse{}
				_ = json.Unmarshal(responseBucket.Get([]byte("body")), &saleResponseBody)
				transaction.SaleResponse = saleResponseBody

			case "void":
				transaction.VoidRequestHeaders = requestHeaders
				transaction.VoidResponseHeaders = responseHeaders
				for key, _ := range responseHeaders {
					if !strings.HasPrefix(key, "X") {
						responseHeaders.Del(key)
					}
				}

				voidRequestBody := models.VoidRequest{}
				_ = json.Unmarshal(requestBucket.Get([]byte("body")), &voidRequestBody)
				transaction.VoidRequest = voidRequestBody

				voidResponseBody := models.VoidResponse{}
				_ = json.Unmarshal(responseBucket.Get([]byte("body")), &voidResponseBody)
				transaction.VoidResponse = voidResponseBody

			case "refund":
				transaction.RefundRequestHeaders = requestHeaders
				transaction.RefundResponseHeaders = responseHeaders
				for key, _ := range responseHeaders {
					if !strings.HasPrefix(key, "X") {
						responseHeaders.Del(key)
					}
				}

				RefundRequestBody := models.RefundRequest{}
				_ = json.Unmarshal(requestBucket.Get([]byte("body")), &RefundRequestBody)
				transaction.RefundRequest = RefundRequestBody

				RefundResponseBody := models.RefundResponse{}
				_ = json.Unmarshal(responseBucket.Get([]byte("body")), &RefundResponseBody)
				transaction.RefundResponse = RefundResponseBody
			}

		}

	}
}

func (t *TransactionRepository) LogRequest(transactionType, action, orderID string, request []byte, headers http.Header) error {
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

		switch transactionType {
		case "sale":
			saleBucket, _ := orderIDBucket.CreateBucketIfNotExists([]byte("sale"))
			saleActionBucket, _ := saleBucket.CreateBucketIfNotExists([]byte(action))
			_ = saleActionBucket.Put([]byte("headers"), headerBytes)
			_ = saleActionBucket.Put([]byte("body"), request)
		case "refund":
			refundBucket, _ := orderIDBucket.CreateBucketIfNotExists([]byte("refund"))
			refundActionBucket, _ := refundBucket.CreateBucketIfNotExists([]byte(action))
			_ = refundActionBucket.Put([]byte("headers"), headerBytes)
			_ = refundActionBucket.Put([]byte("body"), request)
		case "void":
			voidBucket, _ := orderIDBucket.CreateBucketIfNotExists([]byte("void"))
			voidActionBucket, _ := voidBucket.CreateBucketIfNotExists([]byte(action))
			_ = voidActionBucket.Put([]byte("headers"), headerBytes)
			_ = voidActionBucket.Put([]byte("body"), request)
		}

		return nil
	})
	return err
}
