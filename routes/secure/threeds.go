package secure

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"roof/vpos/models"
	"roof/vpos/repository"
	"roof/vpos/routes/check"
	"roof/vpos/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func ThreeDS(c *gin.Context, bolt *repository.Bolt) {
	baseURL := bolt.ConfigRepo.GetBaseURL()
	err := c.Request.ParseForm()
	if err != nil {
		c.HTML(http.StatusBadRequest, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
	}

	ttm := c.Request.Form.Get("threeds-transaction-mode")

	threedsreq := models.ThreeDSRequest{}
	err = c.Bind(&threedsreq)

	if err != nil {
		c.HTML(http.StatusBadRequest, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
	}
	if threedsreq.ReturnUrl == "" {
		threedsreq.ReturnUrl = "http://localhost:8080/return?orderID=" + threedsreq.OrderId
	}
	threedsreqJson, _ := json.Marshal(&threedsreq)

	var req *http.Request

	switch ttm {
	case "sale":
		req, _ = http.NewRequest("POST", baseURL+"/api/ThreeD/Sale", bytes.NewBuffer(threedsreqJson))

	case "presale":
		req, _ = http.NewRequest("POST", baseURL+"/api/ThreeD/PreSale", bytes.NewBuffer(threedsreqJson))

	default:
		c.HTML(http.StatusBadRequest, "result.html", gin.H{"state": 0, "result": "couldn't recognize threeds transaction mode"})
		return
	}

	req.Header, err = utils.CalculateSignature(string(threedsreqJson), bolt)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{"state": 0, "result": err.Error()})
	}

	err = bolt.TransactionRepo.Log("threeds", "request", threedsreq.OrderId, threedsreqJson, req.Header)
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
	var threedsresp models.ThreeDSResponse
	_ = json.Unmarshal(respBody, &response)
	_ = json.Unmarshal(response.Result, &threedsresp)

	err = bolt.TransactionRepo.Log("threeds", "response", threedsreq.OrderId, response.Result, resp.Header)
	if err != nil {
		log.Printf("could not log the threeds request %s\n", err.Error())
	}

	if response.State == 0 {
		responseIndent, _ := json.MarshalIndent(response, "", "	")
		c.HTML(http.StatusOK, "result.html", gin.H{"state": response.State, "result": string(responseIndent)})
		return
	}

	htmlcontent, _ := base64.StdEncoding.DecodeString(threedsresp.HtmlContent)
	// do not show the loading gif.
	re := regexp.MustCompile("<img.*>")
	redirect := re.ReplaceAll(htmlcontent, []byte(""))
	c.Data(http.StatusOK, "text/html; charset=utf-8", redirect)
}

func ThreeDSResult(c *gin.Context, bolt *repository.Bolt) {
	token := utils.TransformToken(c.Query("x_body"))
	tokenreq := models.TokenRequest{Token: token, Lang: "tr"}

	response, err := check.CheckToken(bolt, tokenreq)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
	}
	responseIndent, _ := json.MarshalIndent(response, "", "	")
	c.HTML(200, "result.html", gin.H{"state": response.State, "result": string(responseIndent)})
}

func CompletePayment(c *gin.Context, bolt *repository.Bolt) {
	baseURL := bolt.ConfigRepo.GetBaseURL()
	err := c.Request.ParseForm()
	if err != nil {
		c.HTML(http.StatusBadRequest, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
	}

	var request models.CompletePaymentRequest
	err = c.Bind(&request)
	if err != nil {
		c.HTML(http.StatusBadRequest, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
	}

	requestBody, _ := json.Marshal(&request)

	var httpRequest *http.Request

	httpRequest, _ = http.NewRequest("POST", baseURL+"/api/ThreeD/Complete3DSPayment", bytes.NewBuffer(requestBody))
	httpRequest.Header, err = utils.CalculateSignature(string(requestBody), bolt)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
	}

	err = bolt.TransactionRepo.Log("completepayment", "request", request.OrderID, requestBody, httpRequest.Header)
	if err != nil {
		log.Printf("could not log the complete payment request %s\n", err.Error())
	}

	client := http.Client{
		Timeout: time.Minute,
	}

	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
	}
	defer httpResponse.Body.Close()

	responseBody, _ := io.ReadAll(httpResponse.Body)
	var response models.Response
	_ = json.Unmarshal(responseBody, &response)
	responseIndent, _ := json.MarshalIndent(response, "", "	")
	err = bolt.TransactionRepo.Log("completepayment", "response", request.OrderID, response.Result, httpResponse.Header)
	if err != nil {
		log.Printf("could not log the complete payment response %s\n", err.Error())
	}

	c.HTML(200, "result.html", gin.H{"state": response.State, "result": string(responseIndent)})

}
