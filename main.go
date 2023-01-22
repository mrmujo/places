package main

import (
	"places/configs" //add this
	"places/routes"  //add this

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	configs.ConnectDB()

	routes.UserRoutes(router)

	router.Run("localhost:6000")
}
