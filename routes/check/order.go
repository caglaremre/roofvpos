package check

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

func OrderId(c *gin.Context, bolt *repository.Bolt) {
	baseURL := bolt.ConfigRepo.GetBaseURL()
	err := c.Request.ParseForm()
	if err != nil {
		c.HTML(http.StatusBadRequest, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
	}

	var orderreq models.CheckOrderRequest
	err = c.Bind(&orderreq)
	if err != nil {
		c.HTML(http.StatusBadRequest, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
	}

	orderreqbody, _ := json.Marshal(orderreq)
	req, _ := http.NewRequest("POST", baseURL+"/api/Check/ByOrderId", bytes.NewBuffer(orderreqbody))
	req.Header, err = utils.CalculateSignature(string(orderreqbody), bolt)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
	}
	err = bolt.TransactionRepo.Log("checkorder", "request", orderreq.OrderId, orderreqbody, req.Header)
	if err != nil {
		log.Printf("Could not log the checkorder request %s\n", err.Error())
	}
	client := http.Client{
		Timeout: time.Minute,
	}
	res, err := client.Do(req)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{"state": 0, "result": err.Error()})
	}
	defer res.Body.Close()

	resBody, _ := io.ReadAll(res.Body)
	var response models.Response
	_ = json.Unmarshal(resBody, &response)
	err = bolt.TransactionRepo.Log("checkorder", "response", orderreq.OrderId, response.Result, res.Header)
	if err != nil {
		log.Printf("Could not log the checkorder request %s\n", err.Error())
	}
	result, _ := json.MarshalIndent(response.Result, "", "	")

	c.HTML(http.StatusOK, "result.html", gin.H{"state": response.State, "result": string(result)})
}

func ListOrderId(c *gin.Context, bolt *repository.Bolt) {
	baseURL := bolt.ConfigRepo.GetBaseURL()
	err := c.Request.ParseForm()
	if err != nil {
		c.HTML(http.StatusBadRequest, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
	}
	orderreq := models.CheckOrderRequest{
		OrderId: c.PostForm("listorder-orderid"),
		Lang:    c.PostForm("listorder-lang"),
	}

	orderreqbody, _ := json.Marshal(orderreq)
	req, _ := http.NewRequest("POST", baseURL+"/api/Check/ListByOrderId", bytes.NewBuffer(orderreqbody))
	req.Header, err = utils.CalculateSignature(string(orderreqbody), bolt)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
	}
	err = bolt.TransactionRepo.Log("listorder", "request", orderreq.OrderId, orderreqbody, req.Header)
	if err != nil {
		log.Printf("Could not log the listorder request %s\n", err.Error())
	}
	client := http.Client{
		Timeout: time.Minute,
	}
	res, err := client.Do(req)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{"state": 0, "result": err.Error()})
	}
	defer res.Body.Close()

	resBody, _ := io.ReadAll(res.Body)
	var response models.Response
	_ = json.Unmarshal(resBody, &response)
	err = bolt.TransactionRepo.Log("listorder", "response", orderreq.OrderId, response.Result, res.Header)
	if err != nil {
		log.Printf("Could not log the listorder request %s\n", err.Error())
	}
	result, _ := json.MarshalIndent(response.Result, "", "	")

	c.HTML(http.StatusOK, "result.html", gin.H{"state": response.State, "result": string(result)})
}
