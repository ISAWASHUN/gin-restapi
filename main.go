package main

import (
	"gin-fleamarket/controllers"
	"gin-fleamarket/models"
	"gin-fleamarket/repositories"
	"gin-fleamarket/services"

	"github.com/gin-gonic/gin"
)

func main() {
	items := []models.Item {
		{ID: 1, Name: "Item 1", Price: 100, Description: "This is item 1", SoldOut: false},
		{ID: 2, Name: "Item 2", Price: 200, Description: "This is item 2", SoldOut: false},
		{ID: 3, Name: "Item 3", Price: 300, Description: "This is item 3", SoldOut: true},
	}

	itemRepository := repositories.NewItemMemoryRepository(items)
	itemService := services.NewItemService(itemRepository)
	itemController := controllers.NewItemController(itemService)
	
	r := gin.Default()
	r.GET("/items", itemController.FindAll)
	r.GET("/items/:id", itemController.FindById)
	r.POST("/items", itemController.Create)
	r.PUT("/items/:id", itemController.Update)
	r.DELETE("/items/:id", itemController.Delete)
	r.Run("localhost:8080")
}