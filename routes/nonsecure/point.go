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

func Point(c *gin.Context, bolt *repository.Bolt) {
	baseURL := bolt.ConfigRepo.GetBaseURL()
	err := c.Request.ParseForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	pointReq := models.PointRequest{}
	err = c.Bind(&pointReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	pointReqJson, _ := json.Marshal(&pointReq)

	req, _ := http.NewRequest("POST", baseURL+"/api/Payment/PointInquiry", bytes.NewBuffer(pointReqJson))

	req.Header, err = utils.CalculateSignature(string(pointReqJson), bolt)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
	}

	err = bolt.TransactionRepo.LogRequest("point", "request", pointReq.OrderID, pointReqJson, req.Header)
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
	err = bolt.TransactionRepo.LogRequest("point", "response", pointReq.OrderID, response.Result, resp.Header)
	result, _ := json.MarshalIndent(response, "", "  ")
	c.HTML(http.StatusOK, "result.html", gin.H{"state": response.State, "result": string(result)})
}
