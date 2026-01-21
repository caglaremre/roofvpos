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

func Sale(c *gin.Context, bolt *repository.Bolt) {
	baseURL := bolt.ConfigRepo.GetBaseURL()
	err := c.Request.ParseForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	saleReq := models.SaleRequest{}
	err = c.Bind(&saleReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	saleReqJson, _ := json.Marshal(&saleReq)

	req, _ := http.NewRequest("POST", baseURL+"/api/Payment/Sale", bytes.NewBuffer(saleReqJson))

	req.Header = utils.CalculateSignature(string(saleReqJson), bolt)
	if len(req.Header.Get("x_signature")) < 1 {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "clientToken or secretKey is empty"})
		return
	}

	err = bolt.TransactionRepo.LogRequest("sale", "request", saleReq.OrderID, saleReqJson, req.Header)
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
	err = bolt.TransactionRepo.LogRequest("sale", "response", saleReq.OrderID, response.Result, resp.Header)
	c.JSON(http.StatusOK, response.Result)

}
