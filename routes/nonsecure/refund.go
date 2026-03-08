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

	req.Header, err = utils.CalculateSignature(string(refundReqJson), bolt)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
	}

	err = bolt.TransactionRepo.Log("refund", "request", refundReq.OrderID, refundReqJson, req.Header)
	if err != nil {
		log.Panic(err)
	}

	client := &http.Client{
		Timeout: time.Minute,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Panicf("There is an error: %s", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Panic(err)
		}
	}(resp.Body)
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	var response models.Response
	_ = json.Unmarshal(resBody, &response)
	err = bolt.TransactionRepo.Log("refund", "response", refundReq.OrderID, response.Result, resp.Header)
	result, _ := json.MarshalIndent(response, "", "  ")
	c.HTML(http.StatusOK, "result.html", gin.H{"state": response.State, "result": string(result)})
}
