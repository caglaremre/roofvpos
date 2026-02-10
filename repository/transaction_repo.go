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

	transactionTypeList := []string{"sale", "presale", "postsale", "void", "refund", "point", "threeds", "token"}
	for _, action := range transactionTypeList {
		actionBucket := orderIdBucket.Bucket([]byte(action))
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

			switch action {
			case "sale", "presale":
				transaction.SaleRequestHeaders = requestHeaders
				transaction.SaleResponseHeaders = responseHeaders

				saleRequestBody := models.SaleRequest{}
				_ = json.Unmarshal(requestBucket.Get([]byte("body")), &saleRequestBody)
				transaction.SaleRequest = saleRequestBody

				saleResponseBody := models.SaleResponse{}
				_ = json.Unmarshal(responseBucket.Get([]byte("body")), &saleResponseBody)
				transaction.SaleResponse = saleResponseBody

			case "postsale":
				transaction.PostSaleRequestHeaders = requestHeaders
				transaction.PostSaleResponseHeaders = responseHeaders

				postSaleRequestBody := models.PostSaleRequest{}
				_ = json.Unmarshal(requestBucket.Get([]byte("body")), &postSaleRequestBody)
				transaction.PostSaleRequest = postSaleRequestBody

				postSaleResponseBody := models.PostSaleResponse{}
				_ = json.Unmarshal(responseBucket.Get([]byte("body")), &postSaleResponseBody)
				transaction.PostSaleResponse = postSaleResponseBody

			case "void":
				transaction.VoidRequestHeaders = requestHeaders
				transaction.VoidResponseHeaders = responseHeaders

				voidRequestBody := models.VoidRequest{}
				_ = json.Unmarshal(requestBucket.Get([]byte("body")), &voidRequestBody)
				transaction.VoidRequest = voidRequestBody

				voidResponseBody := models.VoidResponse{}
				_ = json.Unmarshal(responseBucket.Get([]byte("body")), &voidResponseBody)
				transaction.VoidResponse = voidResponseBody

			case "refund":
				transaction.RefundRequestHeaders = requestHeaders
				transaction.RefundResponseHeaders = responseHeaders

				RefundRequestBody := models.RefundRequest{}
				_ = json.Unmarshal(requestBucket.Get([]byte("body")), &RefundRequestBody)
				transaction.RefundRequest = RefundRequestBody

				RefundResponseBody := models.RefundResponse{}
				_ = json.Unmarshal(responseBucket.Get([]byte("body")), &RefundResponseBody)
				transaction.RefundResponse = RefundResponseBody
			case "point":
				transaction.PointRequestHeaders = requestHeaders
				transaction.PointResponseHeaders = responseHeaders

				PointRequestBody := models.PointRequest{}
				_ = json.Unmarshal(requestBucket.Get([]byte("body")), &PointRequestBody)
				transaction.PointRequest = PointRequestBody

				PointResponseBody := models.PointResponse{}
				_ = json.Unmarshal(responseBucket.Get([]byte("body")), &PointResponseBody)
				transaction.PointResponse = PointResponseBody

			case "threeds":
				transaction.ThreeDSRequestHeaders = requestHeaders
				transaction.ThreeDSResponseHeaders = responseHeaders

				threedsRequestBody := models.ThreeDSRequest{}
				_ = json.Unmarshal(requestBucket.Get([]byte("body")), &threedsRequestBody)
				transaction.ThreeDSRequest = threedsRequestBody

				threedsResponseBody := models.ThreeDSResponse{}
				_ = json.Unmarshal(responseBucket.Get([]byte("body")), &threedsResponseBody)
				transaction.ThreeDSResponse = threedsResponseBody

			case "token":
				transaction.TokenRequestHeaders = requestHeaders
				transaction.TokenResponseHeaders = responseHeaders

				tokenRequestBody := models.TokenRequest{}
				_ = json.Unmarshal(requestBucket.Get([]byte("body")), &tokenRequestBody)
				transaction.TokenRequest = tokenRequestBody

				tokenResponseBody := models.TokenResponse{}
				_ = json.Unmarshal(responseBucket.Get([]byte("body")), &tokenResponseBody)
				transaction.TokenResponse = tokenResponseBody
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
		//TODO need to cleans this, we already have the action no need for switch
		switch transactionType {
		case "sale":
			saleBucket, _ := orderIDBucket.CreateBucketIfNotExists([]byte("sale"))
			saleActionBucket, _ := saleBucket.CreateBucketIfNotExists([]byte(action))
			_ = saleActionBucket.Put([]byte("headers"), headerBytes)
			_ = saleActionBucket.Put([]byte("body"), request)
		case "presale":
			presaleBucket, _ := orderIDBucket.CreateBucketIfNotExists([]byte("presale"))
			presaleActionBucket, _ := presaleBucket.CreateBucketIfNotExists([]byte(action))
			_ = presaleActionBucket.Put([]byte("headers"), headerBytes)
			_ = presaleActionBucket.Put([]byte("body"), request)
		case "postsale":
			postsaleBucket, _ := orderIDBucket.CreateBucketIfNotExists([]byte("postsale"))
			postsaleActionBucket, _ := postsaleBucket.CreateBucketIfNotExists([]byte(action))
			_ = postsaleActionBucket.Put([]byte("headers"), headerBytes)
			_ = postsaleActionBucket.Put([]byte("body"), request)
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
		case "point":
			pointBucket, _ := orderIDBucket.CreateBucketIfNotExists([]byte("point"))
			pointActionBucket, _ := pointBucket.CreateBucketIfNotExists([]byte(action))
			_ = pointActionBucket.Put([]byte("headers"), headerBytes)
			_ = pointActionBucket.Put([]byte("body"), request)
		case "threeds":
			threedsBucket, _ := orderIDBucket.CreateBucketIfNotExists([]byte(transactionType))
			threedsActionBucket, _ := threedsBucket.CreateBucketIfNotExists([]byte(action))
			_ = threedsActionBucket.Put([]byte("headers"), headerBytes)
			_ = threedsActionBucket.Put([]byte("body"), request)
		case "token":
			tokenBucket, _ := orderIDBucket.CreateBucketIfNotExists([]byte(transactionType))
			tokenActionBucket, _ := tokenBucket.CreateBucketIfNotExists([]byte(action))
			_ = tokenActionBucket.Put([]byte("headers"), headerBytes)
			_ = tokenActionBucket.Put([]byte("body"), request)
		}

		return nil
	})
	return err
}
