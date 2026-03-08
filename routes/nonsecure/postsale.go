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

func PostSale(c *gin.Context, bolt *repository.Bolt) {
	baseURL := bolt.ConfigRepo.GetBaseURL()
	err := c.Request.ParseForm()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{"state": 0, "result": err.Error()})
	}

	postSaleReq := models.PostSaleRequest{}
	err = c.Bind(&postSaleReq)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
	}

	postSaleReqJson, _ := json.Marshal(&postSaleReq)

	req, _ := http.NewRequest("POST", baseURL+"/api/Payment/PostSale", bytes.NewBuffer(postSaleReqJson))

	req.Header, err = utils.CalculateSignature(string(postSaleReqJson), bolt)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{"state": 0, "result": err.Error()})
	}

	err = bolt.TransactionRepo.Log("postsale", "request", postSaleReq.OrderID, postSaleReqJson, req.Header)
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
	err = bolt.TransactionRepo.Log("postsale", "response", postSaleReq.OrderID, response.Result, resp.Header)
	result, _ := json.MarshalIndent(response, "", "  ")
	c.HTML(http.StatusOK, "result.html", gin.H{"state": response.State, "result": string(result)})
}
