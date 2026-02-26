package check

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"roof/vpos/models"
	"roof/vpos/repository"
	"roof/vpos/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func Token(c *gin.Context, bolt *repository.Bolt) {

	err := c.Request.ParseForm()
	if err != nil {
		c.HTML(http.StatusBadRequest, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
	}
	var tokenreq models.TokenRequest
	err = c.Bind(&tokenreq)

	if err != nil {
		c.HTML(http.StatusBadRequest, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
	}

	response, err := CheckToken(bolt, tokenreq)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
	}
	responseIndent, _ := json.MarshalIndent(response.Result, "", "	")
	c.HTML(200, "result.html", gin.H{"state": response.State, "result": string(responseIndent)})

}

func CheckToken(bolt *repository.Bolt, tokenreq models.TokenRequest) (models.Response, error) {
	var response models.Response
	baseURL := bolt.ConfigRepo.GetBaseURL()
	tokenreqjson, _ := json.Marshal(tokenreq)
	var req *http.Request

	req, _ = http.NewRequest("POST", baseURL+"/api/Check/ByToken", bytes.NewBuffer(tokenreqjson))
	req.Header = utils.CalculateSignature(string(tokenreqjson), bolt)
	if len(req.Header.Get("x_signature")) < 1 {
		return response, errors.New("clientToken or secretKey is empty")
	}

	client := &http.Client{
		Timeout: time.Minute,
	}

	resp, err := client.Do(req)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	var tokenresp models.TokenResponse
	_ = json.Unmarshal(respBody, &response)
	_ = json.Unmarshal(response.Result, &tokenresp)

	err = bolt.TransactionRepo.LogRequest("token", "request", tokenresp.OrderId, tokenreqjson, req.Header)
	if err != nil {
		log.Printf("could not log the token request %s\n", err.Error())
	}

	err = bolt.TransactionRepo.LogRequest("token", "response", tokenresp.OrderId, response.Result, resp.Header)
	if err != nil {
		log.Printf("could not log the token response %s\n", err.Error())
	}
	return response, nil
}
