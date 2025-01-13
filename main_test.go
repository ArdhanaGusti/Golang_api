package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func setupTestDB() *gorm.DB {
	dsn := "root@tcp(127.0.0.1:3306)/go-api?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}

	return database
}

func clearDB() {
	testDB.Exec("DELETE FROM users")
}

func TestRegisterUser(t *testing.T) {
	// testDB = setupTestDB()
	// defer clearDB()

	router := setupRouter()
	// router.POST("/users", func(c *gin.Context) {
	// 	var userPayload validation.RegisterUserPayload
	// 	if err := c.ShouldBind(&userPayload); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{
	// 			"error": err.Error(),
	// 		})
	// 		return
	// 	}

	// 	hash, err := bcrypt.GenerateFromPassword([]byte(userPayload.Password), bcrypt.DefaultCost)
	// 	if err != nil {
	// 		c.JSON(400, gin.H{"status": "Hashing failed"})
	// 		c.Abort()
	// 		return
	// 	}

	// 	newUser := models.User{
	// 		Username: userPayload.Username,
	// 		Fullname: userPayload.Fullname,
	// 		Email:    userPayload.Email,
	// 		Password: string(hash),
	// 	}

	// 	if err := testDB.Create(&newUser).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
	// 		return
	// 	}

	// 	c.JSON(http.StatusCreated, newUser)
	// })

	payload := `{
		"Username": "Rena",
		"Fullname": "Rena Aliana",
		"Email": "rena.aliana@yahoo.com",
		"Password": "admin123"
	}`
	// body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Body.String())

	// assert.Equal(t, http.StatusCreated, w.Code)

	// var response models.User
	// err := json.Unmarshal(w.Body.Bytes(), &response)
	// assert.NoError(t, err)
	// assert.Equal(t, payload["Username"], response.Username)
	// assert.Equal(t, payload["Fullname"], response.Fullname)
	// assert.Equal(t, payload["Email"], response.Email)
	// err = bcrypt.CompareHashAndPassword([]byte(response.Password), []byte(payload["Password"]))
	// assert.NoError(t, err)

	// var user models.User
	// result := testDB.First(&user, "email = ?", payload["Email"])
	// assert.Nil(t, result.Error)
	// assert.Equal(t, payload["Username"], user.Username)
	// assert.Equal(t, payload["Fullname"], user.Fullname)
	// assert.Equal(t, payload["Email"], user.Email)
	// err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload["Password"]))
	// assert.NoError(t, err)
}
