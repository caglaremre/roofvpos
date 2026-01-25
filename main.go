package main

import (
	"log"
	"roof/vpos/repository"
	"roof/vpos/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	bolt, err := repository.InitBolt("./bolt.db")
	if err != nil {
		panic(err)
	}
	err = bolt.InitialBuckets()
	if err != nil {
		panic(err)
	}
	defer bolt.CloseBolt()

	//init gin
	router := gin.Default()

	router.StaticFile("/favicon.ico", "./assets/favicon.ico")
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*")

	routes.RegisterRoutes(bolt, router)

	router.Use(func(context *gin.Context) {
		context.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		context.Next()
	})

	log.Fatal(router.Run("localhost:8080"))
}
