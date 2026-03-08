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

func ProcessId(c *gin.Context, bolt *repository.Bolt) {
	baseURL := bolt.ConfigRepo.GetBaseURL()
	err := c.Request.ParseForm()
	if err != nil {
		c.HTML(http.StatusBadRequest, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
	}

	var process models.CheckProcessRequest
	err = c.Bind(&process)
	if err != nil {
		c.HTML(http.StatusBadRequest, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
	}

	processbody, _ := json.Marshal(process)
	req, _ := http.NewRequest("POST", baseURL+"/api/Check/ByProcessId", bytes.NewBuffer(processbody))
	req.Header, err = utils.CalculateSignature(string(processbody), bolt)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "result.html", gin.H{"state": 0, "result": err.Error()})
		return
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
	var processresp models.CheckProcessResponse
	_ = json.Unmarshal(resBody, &response)
	_ = json.Unmarshal(response.Result, &processresp)

	err = bolt.TransactionRepo.Log("checkprocess", "request", processresp.OrderId, processbody, req.Header)
	if err != nil {
		log.Printf("Could not log the checkprocess request %s\n", err.Error())
	}

	err = bolt.TransactionRepo.Log("checkprocess", "response", processresp.OrderId, response.Result, res.Header)
	if err != nil {
		log.Printf("Could not log the checkprocess request %s\n", err.Error())
	}
	result, _ := json.MarshalIndent(response.Result, "", "	")

	c.HTML(http.StatusOK, "result.html", gin.H{"state": response.State, "result": string(result)})
}
