package main

import (
	"bytes"
	"encoding/json"
	"gin-fleamarket/dto"
	"gin-fleamarket/infra"
	"gin-fleamarket/models"
	"gin-fleamarket/services"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load(".env.test"); err != nil {
		panic("Failed to load env file")
	}

	code := m.Run()

	os.Exit(code)
}

func setupTestData(db *gorm.DB) {
	items := []models.Item {
		{Name: "Item 1", Price: 100, Description: "This is item 1", SoldOut: false, UserID: 1},
		{Name: "Item 2", Price: 200, Description: "This is item 2", SoldOut: false, UserID: 2},
		{Name: "Item 3", Price: 300, Description: "This is item 3", SoldOut: true, UserID: 3},
	}

	users := []models.User {
		{Email: "test1@test.com", Password: "password1"},
		{Email: "test2@test.com", Password: "password2"},
	}

	for _, user := range users {
		db.Create(&user)
	}

	for _, item := range items {
		db.Create(&item)
	}
}

func setup() *gin.Engine {
	db := infra.SetupDB()
	db.AutoMigrate(&models.User{}, &models.Item{})
	setupRouter(db)
	router := setupRouter(db)

	return router
}

func TestFindAll(t *testing.T) {
	router := setup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/items", nil)
	router.ServeHTTP(w, req)

	//APIの実行結果を取得
	var res map[string][]models.Item
	json.Unmarshal([]byte(w.Body.String()), &res)

	//レスポンスのステータスコードが200かどうか
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 3, len(res["data"]))
}

func TestCreate(t *testing.T) {
	router := setup()

	token, err := services.CreateToken(1, "test1@test.com")
	assert.Equal(t, nil, err)

	createItemInput := dto.CreateItemInput{
		Name: "Item 4",
		Price: 400,
		Description: "This is item 4",
	}

	reqBody, _ := json.Marshal(createItemInput)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer " + *token)
	router.ServeHTTP(w, req)

	//APIの実行結果を取得
	var res map[string][]models.Item
	json.Unmarshal([]byte(w.Body.String()), &res)

	//レスポンスのステータスコードが200かどうか
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, uint(4), res["data"].ID)
}


func TestCreateUnauthorized(t *testing.T) {
	router := setup()

	createItemInput := dto.CreateItemInput{
		Name: "Item 4",
		Price: 400,
		Description: "This is item 4",
	}

	reqBody, _ := json.Marshal(createItemInput)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(reqBody))
	router.ServeHTTP(w, req)

	//APIの実行結果を取得
	var res map[string][]models.Item
	json.Unmarshal([]byte(w.Body.String()), &res)

	//レスポンスのステータスコードが200かどうか
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}