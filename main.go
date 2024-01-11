package main

import (
	"gin-fleamarket/controllers"
	"gin-fleamarket/infra"
	// "gin-fleamarket/models"
	"gin-fleamarket/repositories"
	"gin-fleamarket/services"

	"github.com/gin-gonic/gin"
)

func main() {
	infra.Initialize()
	db := infra.SetupDB()
	// items := []models.Item {
	// 	{ID: 1, Name: "Item 1", Price: 100, Description: "This is item 1", SoldOut: false},
	// 	{ID: 2, Name: "Item 2", Price: 200, Description: "This is item 2", SoldOut: false},
	// 	{ID: 3, Name: "Item 3", Price: 300, Description: "This is item 3", SoldOut: true},
	// }

	// itemRepository := repositories.NewItemMemoryRepository(items)

	itemRepository := repositories.NewItemRepository(db)
	itemService := services.NewItemService(itemRepository)
	itemController := controllers.NewItemController(itemService)
	
	r := gin.Default()
	itemRouter := r.Group("/items")
	itemRouter.GET("", itemController.FindAll)
	itemRouter.GET("/:id", itemController.FindById)
	itemRouter.POST("", itemController.Create)
	itemRouter.PUT("/:id", itemController.Update)
	itemRouter.DELETE("/:id", itemController.Delete)
	r.Run("localhost:8080")
}