package repository

import (
	"cmp"
	"encoding/json"
	"net/http"
	"roof/vpos/models"
	"slices"
	"time"

	"go.etcd.io/bbolt"
)

type TransactionRepository struct {
	DB *bbolt.DB
}

// todo seperata the transactions and its details to different methods
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

			saleBucket := orderIdBucket.Bucket([]byte("sale"))

			requestBucket := saleBucket.Bucket([]byte("request"))

			saleRequestBody := models.SaleRequest{}
			_ = json.Unmarshal(requestBucket.Get([]byte("body")), &saleRequestBody)
			transaction.SaleRequest = saleRequestBody
			saleRequestBodyJson, _ := json.MarshalIndent(saleRequestBody, "", "	")
			transaction.SaleRequestBody = string(saleRequestBodyJson)

			saleRequestHeaders := http.Header{}
			_ = json.Unmarshal(requestBucket.Get([]byte("headers")), &saleRequestHeaders)
			transaction.SaleRequestsHeaders = saleRequestHeaders

			responseBucket := saleBucket.Bucket([]byte("response"))

			saleResponseBody := models.SaleResponse{}
			_ = json.Unmarshal(responseBucket.Get([]byte("body")), &saleResponseBody)
			transaction.SaleResponse = saleResponseBody
			saleResponseBodyJson, _ := json.MarshalIndent(saleResponseBody, "", "	")
			transaction.SaleResponseBody = string(saleResponseBodyJson)

			saleResponseHeaders := http.Header{}
			_ = json.Unmarshal(responseBucket.Get([]byte("headers")), &saleResponseHeaders)
			transaction.SaleResponseHeaders = saleResponseHeaders

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
			refundBucket, _ := orderIDBucket.CreateBucket([]byte("refund"))
			refundRequestBucket, _ := refundBucket.CreateBucket([]byte(action))
			_ = refundRequestBucket.Put([]byte("headers"), headerBytes)
			_ = refundRequestBucket.Put([]byte("body"), request)
		case "void":
			voidBucket, _ := orderIDBucket.CreateBucket([]byte("void"))
			voidRequestBucket, _ := voidBucket.CreateBucket([]byte(action))
			_ = voidRequestBucket.Put([]byte("headers"), headerBytes)
			_ = voidRequestBucket.Put([]byte("body"), request)
		}

		return nil
	})
	return err
}
