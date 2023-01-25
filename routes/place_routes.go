package routes

import (
	"places/controllers"

	"github.com/gin-gonic/gin"
)

func PlaceRoutes(router *gin.Engine) {
	router.POST("/place", controllers.CreatePlace())
	router.GET("/place/:placeId", controllers.GetAPlace())
	router.PUT("/place/:placeId", controllers.EditAPlace())
	router.DELETE("/place/:placeId", controllers.DeleteAPlace())
	router.GET("/places", controllers.GetAllPlaces())
}
