package secure

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

func ThreedsHosting(c *gin.Context, bolt *repository.Bolt) {
	//c.HTML(0, "wait.html", nil)
	baseURL := bolt.ConfigRepo.GetBaseURL()
	log.Println(baseURL)
	err := c.Request.ParseForm()
	if err != nil {
		c.HTML(http.StatusBadRequest, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
	}
	threedshostingreq := models.ThreeDSHostingRequest{}
	err = c.Bind(&threedshostingreq)
	if err != nil {
		c.HTML(http.StatusBadRequest, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
	}
	if threedshostingreq.ReturnUrl == "" {
		threedshostingreq.ReturnUrl = "http://localhost:8080/return?orderID=" + threedshostingreq.OrderId
	}

	var req *http.Request

	req, _ = http.NewRequest("POST", baseURL+"/api/Order/CreateOrder3D", bytes.NewBuffer([]byte(threedshostingreq.ToJson(false))))

	req.Header = utils.CalculateSignature(threedshostingreq.ToJson(false), bolt)

	err = bolt.TransactionRepo.LogRequest("threedshosting", "request", threedshostingreq.OrderId, []byte(threedshostingreq.ToJson(false)), req.Header)
	if err != nil {
		log.Printf("could not log the threeds request %s\n", err.Error())
	}

	if len(req.Header.Get("x_signature")) < 1 {
		c.HTML(http.StatusBadRequest, "result.html", gin.H{"state": 0, "result": "clientToken or secretKey is empty"})
		return
	}

	client := &http.Client{
		Timeout: time.Minute,
	}
	resp, err := client.Do(req)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
	}

	var response models.Response
	var threedshostingresp models.ThreeDSHostingResponse
	_ = json.Unmarshal(respBody, &response)
	_ = json.Unmarshal(response.Result, &threedshostingresp)
	responseIndent, _ := json.MarshalIndent(response, "", "	")

	err = bolt.TransactionRepo.LogRequest("threedshosting", "response", threedshostingreq.OrderId, response.Result, resp.Header)
	if err != nil {
		log.Printf("could not log the threeds request %s\n", err.Error())
	}

	c.HTML(http.StatusOK, "result.html", gin.H{"state": response.State, "result": string(responseIndent)})
}
