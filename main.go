package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"roof/vpos/repository"
	"roof/vpos/routes"

	"github.com/gin-gonic/gin"
)

//go:embed templates/*
var templateFS embed.FS

//go:embed assets/*
var staticFS embed.FS

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
	tmpl := template.Must(template.ParseFS(templateFS, "templates/*"))
	router.SetHTMLTemplate(tmpl)

	if len(os.Args) > 1 && os.Args[1] == "e" {
		staticFS, _ := fs.Sub(staticFS, "assets")
		router.StaticFS("/assets", http.FS(staticFS))
		router.GET("/favicon.ico", func(c *gin.Context) {
			fileData, err := fs.ReadFile(staticFS, "favicon.ico")
			if err != nil {
				// If this triggers, your file isn't where you think it is
				c.String(404, "Favicon not found")
				return
			}
			c.Data(200, "image/x-icon", fileData)
		})
	} else {
		router.StaticFile("/favicon.ico", "./assets/favicon.ico")
		router.Static("/assets", "./assets")
		router.LoadHTMLGlob("templates/*")
	}
	routes.RegisterRoutes(bolt, router)

	router.Use(func(context *gin.Context) {
		context.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		context.Next()
	})
	log.Printf("################################################\n")
	log.Printf("# 📡 Server starting on http://localhost:8080  #\n")
	log.Printf("###############################################\n")
	log.Fatal(router.Run("localhost:8080"))
}
