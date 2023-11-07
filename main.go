package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors(r *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"*"}
	config.AllowMethods = []string{"*"}
	r.Use(cors.New(config))
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	Cors(app)
	proxiesHandle(app)
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	}
	err := app.Run(port)
	if err != nil {
		log.Println(err)
	}
}
