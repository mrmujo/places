package main

import (
	"places/configs"
	"places/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	configs.ConnectDB()

	routes.PlaceRoutes(router)

	router.Run("localhost:6000")
}
