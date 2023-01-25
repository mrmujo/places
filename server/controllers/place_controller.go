package controllers

import (
	"context"
	"net/http"
	"places/configs"
	"places/models"
	"places/responses"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var placeCollection *mongo.Collection = configs.GetCollection(configs.DB, "places")
var validate = validator.New()

func CreatePlace() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var place models.Place
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&place); err != nil {
			c.JSON(http.StatusBadRequest, responses.PlaceResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&place); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.PlaceResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newPlace := models.Place{
			Id:       primitive.NewObjectID(),
			Name:     place.Name,
			Description: place.Description,
			Latitude:    place.Latitude,
			Longitude:    place.Longitude,
		}

		result, err := placeCollection.InsertOne(ctx, newPlace)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.PlaceResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.PlaceResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetAPlace() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		placeId := c.Param("placeId")
		var place models.Place
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(placeId)

		err := placeCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&place)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.PlaceResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.PlaceResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": place}})
	}
}

func EditAPlace() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		placeId := c.Param("placeId")
		var place models.Place
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(placeId)

		//validate the request body
		if err := c.BindJSON(&place); err != nil {
			c.JSON(http.StatusBadRequest, responses.PlaceResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&place); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.PlaceResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{"name": place.Name, "description": place.Description, "latitude": place.Latitude, "longitude": place.Longitude}
		result, err := placeCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.PlaceResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//get updated user details
		var updatedUser models.Place
		if result.MatchedCount == 1 {
			err := placeCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.PlaceResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.PlaceResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedUser}})
	}
}

func DeleteAPlace() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		placeId := c.Param("placeId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(placeId)

		result, err := placeCollection.DeleteOne(ctx, bson.M{"id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.PlaceResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.PlaceResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Place with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			responses.PlaceResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Place successfully deleted!"}},
		)
	}
}

func GetAllPlaces() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var places []models.Place
		defer cancel()

		results, err := placeCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.PlaceResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleUser models.Place
			if err = results.Decode(&singleUser); err != nil {
				c.JSON(http.StatusInternalServerError, responses.PlaceResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			places = append(places, singleUser)
		}

		c.JSON(http.StatusOK,
			responses.PlaceResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": places}},
		)
	}
}
