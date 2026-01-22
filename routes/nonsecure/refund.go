package nonsecure

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"roof/vpos/models"
	"roof/vpos/repository"
	"roof/vpos/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func Refund(c *gin.Context, bolt *repository.Bolt) {
	baseURL := bolt.ConfigRepo.GetBaseURL()
	err := c.Request.ParseForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	refundReq := models.RefundRequest{}
	err = c.Bind(&refundReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	refundReqJson, _ := json.Marshal(&refundReq)

	req, _ := http.NewRequest("POST", baseURL+"/api/Payment/Refund", bytes.NewBuffer(refundReqJson))

	req.Header = utils.CalculateSignature(string(refundReqJson), bolt)
	if len(req.Header.Get("x_signature")) < 1 {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "clientToken or secretKey is empty"})
		return
	}

	err = bolt.TransactionRepo.LogRequest("refund", "request", refundReq.OrderID, refundReqJson, req.Header)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Timeout: time.Minute,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("There is an error: %s", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	var response models.Response
	_ = json.Unmarshal(resBody, &response)
	err = bolt.TransactionRepo.LogRequest("refund", "response", refundReq.OrderID, response.Result, resp.Header)
	result, _ := json.MarshalIndent(response, "", "  ")
	c.HTML(http.StatusOK, "result.html", gin.H{"state": response.State, "result": string(result)})
}
