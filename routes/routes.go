package routes

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
	"github.com/google/uuid"
)

var BASEURL string
var BOLT *repository.Bolt

func RegisterRoutes(bolt *repository.Bolt, server *gin.Engine) {
	BOLT = bolt
	BASEURL = BOLT.ConfigRepo.GetBaseURL()
	server.NoRoute(noRoute)
	server.GET("/", home)
	server.POST("/config", configUpdate)
	server.POST("/sale", saleTest)
}

func home(c *gin.Context) {
	clientToken, secretKey := BOLT.ConfigRepo.GetClientAndSecretKey()
	transactions := BOLT.TransactionRepo.GetAllTransactions()
	orderID, _ := uuid.NewUUID()
	saleReq := models.SaleRequest{OrderID: orderID.String()}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"clientToken":  clientToken,
		"secretKey":    secretKey,
		"transactions": transactions,
		"saleReq":      saleReq,
	})
}

func noRoute(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", nil)
}

func configUpdate(c *gin.Context) {
	var data map[string]string
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if data["clientToken"] == "" || data["secretKey"] == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "clientToken or secretKey is empty"})
		return
	}
	err := BOLT.ConfigRepo.UpdateClientAndSecretKey(data["clientToken"], data["secretKey"])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func saleTest(c *gin.Context) {
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

	req, _ := http.NewRequest("POST", BASEURL+"/api/Payment/Sale", bytes.NewBuffer(saleReqJson))

	req.Header = utils.CalculateSignature(string(saleReqJson), BOLT)
	if len(req.Header.Get("x_signature")) < 1 {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "clientToken or secretKey is empty"})
		return
	}

	err = BOLT.TransactionRepo.LogRequest("sale", "request", saleReq.OrderID, saleReqJson, req.Header)
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
	err = BOLT.TransactionRepo.LogRequest("sale", "response", saleReq.OrderID, response.Result, resp.Header)
	c.JSON(http.StatusOK, response.Result)

}
