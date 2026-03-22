package routes

import (
	"net/http"
	"roof/vpos/models"
	"roof/vpos/repository"
	"roof/vpos/routes/check"
	"roof/vpos/routes/nonsecure"
	"roof/vpos/routes/secure"

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
	server.POST("/void", func(c *gin.Context) {
		nonsecure.Void(c, bolt)
	})
	server.POST("/refund", func(c *gin.Context) {
		nonsecure.Refund(c, bolt)
	})
	server.POST("/pre", func(c *gin.Context) {
		nonsecure.Presale(c, bolt)
	})
	server.POST("/post", func(c *gin.Context) {
		nonsecure.PostSale(c, bolt)
	})
	server.POST("/point", func(c *gin.Context) {
		nonsecure.Point(c, bolt)
	})
	server.POST("/threeds", func(c *gin.Context) {
		secure.ThreeDS(c, bolt)
	})
	server.GET("/return", func(c *gin.Context) {
		secure.ThreeDSResult(c, bolt)
	})
	server.POST("/completepayment", func(c *gin.Context) {
		secure.CompletePayment(c, bolt)
	})
	server.POST("/threedshosting", func(c *gin.Context) {
		secure.ThreedsHosting(c, bolt)
	})
	server.POST("/checktoken", func(c *gin.Context) {
		check.Token(c, bolt)
	})
	server.POST("/checkorder", func(c *gin.Context) {
		check.OrderId(c, bolt)
	})
	server.POST("/checkprocess", func(c *gin.Context) {
		check.ProcessId(c, bolt)
	})
	server.POST("/listorder", func(c *gin.Context) {
		check.ListOrderId(c, bolt)
	})
	server.GET("/transaction/:orderid", func(c *gin.Context) {
		getTransactionDetail(c, bolt)
	})
}

func home(c *gin.Context, bolt *repository.Bolt) {
	// TODO return only transactions orderids then fetch when needed
	clientToken, secretKey := bolt.ConfigRepo.GetClientAndSecretKey()
	baseURL := bolt.ConfigRepo.GetBaseURL()
	transactions := bolt.TransactionRepo.GetAllTransactions()
	orderID, _ := uuid.NewUUID()
	saleReq := models.SaleRequest{OrderID: orderID.String()}
	c.HTML(http.StatusOK, "home.html", gin.H{
		"clientToken":  clientToken,
		"secretKey":    secretKey,
		"transactions": transactions,
		"saleReq":      saleReq,
		"baseUrl":      baseURL,
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
	if data["clientToken"] == "" || data["secretKey"] == "" || data["baseUrl"] == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "clientToken, secretKey or baseUrl is empty"})
		return
	}
	err := bolt.ConfigRepo.UpdateConfig(data["clientToken"], data["secretKey"], data["baseUrl"])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func getTransactionDetail(c *gin.Context, bolt *repository.Bolt) {
	var transaction models.Transaction
	transaction.OrderID = c.Param("orderid")
	err := bolt.TransactionRepo.GetTransactionDetail(&transaction)
	if err != nil {
		noRoute(c)
		return
	}
	c.HTML(http.StatusOK, "transaction.html", gin.H{"transaction": transaction})
}
