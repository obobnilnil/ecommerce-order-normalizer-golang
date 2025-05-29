package main

import (
	"log"

	configLoader "e-commerce/internal/e-commerce/config"
	"e-commerce/internal/e-commerce/server"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "X-Auth-Token", "Authorization"}
	router.Use(cors.New(config))

	cleanerMap, err := configLoader.LoadCleanerConfig("internal/e-commerce/config/cleaner.yml")
	if err != nil {
		log.Fatalf("failed to load cleaner config: %v", err)
	}

	server.SetupRoutesEcommerce(router, cleanerMap)

	err = router.Run(":9000")
	if err != nil {
		panic(err.Error())
	}
}
