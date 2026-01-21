package routes

import (
	"net/http"
	"roof/vpos/models"
	"roof/vpos/repository"
	"roof/vpos/routes/nonsecure"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterRoutes(bolt *repository.Bolt, server *gin.Engine) {
	server.NoRoute(noRoute)
	server.GET("/", func(c *gin.Context) {
		home(c, bolt)
	})
	server.POST("/config", func(c *gin.Context) {
		configUpdate(c, bolt)
	})
	server.POST("/sale", func(c *gin.Context) {
		nonsecure.Sale(c, bolt)
	})
}

func home(c *gin.Context, bolt *repository.Bolt) {
	clientToken, secretKey := bolt.ConfigRepo.GetClientAndSecretKey()
	transactions := bolt.TransactionRepo.GetAllTransactions()
	orderID, _ := uuid.NewUUID()
	saleReq := models.SaleRequest{OrderID: orderID.String()}
	c.HTML(http.StatusOK, "home.html", gin.H{
		"clientToken":  clientToken,
		"secretKey":    secretKey,
		"transactions": transactions,
		"saleReq":      saleReq,
	})
}

func noRoute(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", nil)
}

func configUpdate(c *gin.Context, bolt *repository.Bolt) {
	var data map[string]string
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if data["clientToken"] == "" || data["secretKey"] == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "clientToken or secretKey is empty"})
		return
	}
	err := bolt.ConfigRepo.UpdateClientAndSecretKey(data["clientToken"], data["secretKey"])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
