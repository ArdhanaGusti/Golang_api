package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ArdhanaGusti/Golang_api/config"
	"github.com/ArdhanaGusti/Golang_api/handler/validation"
	"github.com/ArdhanaGusti/Golang_api/models"
	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"
)

type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

func Initialize() {
	gotenv.Load()
	config.InitDB()
}

func TestRegisterUser(t *testing.T) {
	Initialize()
	config.MigrateFreshDB()
	router := setupRouter()

	newUser := validation.RegisterUserPayload{
		Username: "Rena",
		Fullname: "Rena Aliana",
		Email:    "rena.aliana@yahoo.com",
		Password: "admin123",
	}
	body, _ := json.Marshal(newUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedResponse := fmt.Sprintf(`{"status":"User %s Registered Successfully"}`, newUser.Fullname)
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestLoginUser(t *testing.T) {
	Initialize()
	router := setupRouter()

	existUser := validation.LoginUserPayload{
		Email:    "rena.aliana@yahoo.com",
		Password: "admin123",
	}
	body, _ := json.Marshal(existUser)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var actualResponse LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)
	expectedMessage := "Login " + existUser.Email + " Successfully"
	assert.Equal(t, expectedMessage, actualResponse.Message)
}

func TestGetUser(t *testing.T) {
	Initialize()
	router := setupRouter()

	existUser := validation.LoginUserPayload{
		Email:    "rena.aliana@yahoo.com",
		Password: "admin123",
	}
	body, _ := json.Marshal(existUser)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req1.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w1, req1)
	var actualResponse LoginResponse
	err1 := json.Unmarshal(w1.Body.Bytes(), &actualResponse)
	assert.NoError(t, err1)

	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/auth/profile", nil)
	req2.Header.Set("Authorization", actualResponse.Token)

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	var user models.User
	err2 := json.Unmarshal(w2.Body.Bytes(), &user)
	assert.NoError(t, err2)
	assert.Equal(t, "rena.aliana@yahoo.com", user.Email)
	assert.Equal(t, "Rena Aliana", user.Fullname)
	assert.Equal(t, "Rena", user.Username)
}

func TestChangeEmailUser(t *testing.T) {
	Initialize()
	router := setupRouter()

	existUser := validation.LoginUserPayload{
		Email:    "rena.aliana@yahoo.com",
		Password: "admin123",
	}
	body, _ := json.Marshal(existUser)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req1.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w1, req1)
	var actualResponse LoginResponse
	err1 := json.Unmarshal(w1.Body.Bytes(), &actualResponse)
	assert.NoError(t, err1)

	form := url.Values{}
	form.Add("Role", "admin")

	req2, _ := http.NewRequest(http.MethodPatch, "/api/v1/auth/change-role", strings.NewReader(form.Encode()))
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req2.Header.Set("Authorization", actualResponse.Token)

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	assert.JSONEq(t, `{"message": "Change role to admin Successfully"}`, w2.Body.String())
}

func TestCreateArticle(t *testing.T) {
	Initialize()
	router := setupRouter()

	existUser := validation.LoginUserPayload{
		Email:    "rena.aliana@yahoo.com",
		Password: "admin123",
	}
	body1, _ := json.Marshal(existUser)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w1, req1)
	var actualResponse LoginResponse
	err1 := json.Unmarshal(w1.Body.Bytes(), &actualResponse)
	assert.Equal(t, http.StatusOK, w1.Code)
	assert.NoError(t, err1)

	newArticle := validation.CreateArticlePayload{
		Title: "Tupai terbang",
		Desc:  "Tupai itu terbang ke langit ke 100 dan membawa hadiah.",
		Tag:   "fiction",
	}
	body2, _ := json.Marshal(newArticle)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodPost, "/api/v1/article", bytes.NewBuffer(body2))
	req2.Header.Set("Authorization", actualResponse.Token)
	req2.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	expectedResponse := fmt.Sprintf(`{"message":"Article %s Made Successfully"}`, newArticle.Title)
	assert.Equal(t, expectedResponse, w2.Body.String())
}

func TestGetArticles(t *testing.T) {
	Initialize()
	router := setupRouter()

	existUser := validation.LoginUserPayload{
		Email:    "rena.aliana@yahoo.com",
		Password: "admin123",
	}
	body1, _ := json.Marshal(existUser)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w1, req1)
	var actualResponse LoginResponse
	err1 := json.Unmarshal(w1.Body.Bytes(), &actualResponse)
	assert.Equal(t, http.StatusOK, w1.Code)
	assert.NoError(t, err1)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/", nil)
	req2.Header.Set("Authorization", actualResponse.Token)

	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	var articles []models.Article
	err2 := json.Unmarshal(w2.Body.Bytes(), &articles)
	assert.NoError(t, err2)
	assert.GreaterOrEqual(t, 1, len(articles))
	assert.Equal(t, "Tupai terbang", articles[0].Title)
	assert.GreaterOrEqual(t, "Tupai itu terbang ke langit ke 100 dan membawa hadiah.", articles[0].Desc)
	assert.GreaterOrEqual(t, "fiction", articles[0].Tag)
}

func TestGetArticle(t *testing.T) {
	Initialize()
	router := setupRouter()

	existUser := validation.LoginUserPayload{
		Email:    "rena.aliana@yahoo.com",
		Password: "admin123",
	}
	body1, _ := json.Marshal(existUser)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w1, req1)
	var actualResponse LoginResponse
	err1 := json.Unmarshal(w1.Body.Bytes(), &actualResponse)
	assert.Equal(t, http.StatusOK, w1.Code)
	assert.NoError(t, err1)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/", nil)
	req2.Header.Set("Authorization", actualResponse.Token)

	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	var articles []models.Article
	err2 := json.Unmarshal(w2.Body.Bytes(), &articles)
	assert.NoError(t, err2)

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest(http.MethodGet, "/api/v1/article/"+articles[0].Slug, nil)

	router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)
	var article models.Article
	err3 := json.Unmarshal(w3.Body.Bytes(), &article)
	assert.NoError(t, err3)
	assert.Equal(t, "Tupai terbang", article.Title)
	assert.GreaterOrEqual(t, "Tupai itu terbang ke langit ke 100 dan membawa hadiah.", article.Desc)
	assert.GreaterOrEqual(t, "fiction", article.Tag)
}

func TestUpdateArticle(t *testing.T) {
	Initialize()
	router := setupRouter()

	existUser := validation.LoginUserPayload{
		Email:    "rena.aliana@yahoo.com",
		Password: "admin123",
	}
	body1, _ := json.Marshal(existUser)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w1, req1)
	var actualResponse LoginResponse
	err1 := json.Unmarshal(w1.Body.Bytes(), &actualResponse)
	assert.Equal(t, http.StatusOK, w1.Code)
	assert.NoError(t, err1)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/", nil)
	req2.Header.Set("Authorization", actualResponse.Token)

	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	var articles []models.Article
	err2 := json.Unmarshal(w2.Body.Bytes(), &articles)
	assert.NoError(t, err2)

	newArticle := validation.CreateArticlePayload{
		Title: "Tupai berdiri",
		Desc:  "Tupai itu berdiri ke arah timur dan membawa hadiah.",
		Tag:   "fiction",
	}
	body3, _ := json.Marshal(newArticle)
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest(http.MethodPut, "/api/v1/article/"+articles[0].Slug, bytes.NewBuffer(body3))
	req3.Header.Set("Authorization", actualResponse.Token)
	req3.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)
	expectedResponse := fmt.Sprintf(`{"message":"Article %s Updated Successfully"}`, newArticle.Title)
	assert.Equal(t, expectedResponse, w3.Body.String())
}

func TestDeleteArticle(t *testing.T) {
	Initialize()
	router := setupRouter()

	existUser := validation.LoginUserPayload{
		Email:    "rena.aliana@yahoo.com",
		Password: "admin123",
	}
	body1, _ := json.Marshal(existUser)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w1, req1)
	var actualResponse LoginResponse
	err1 := json.Unmarshal(w1.Body.Bytes(), &actualResponse)
	assert.Equal(t, http.StatusOK, w1.Code)
	assert.NoError(t, err1)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/", nil)
	req2.Header.Set("Authorization", actualResponse.Token)

	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	var articles []models.Article
	err2 := json.Unmarshal(w2.Body.Bytes(), &articles)
	assert.NoError(t, err2)

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest(http.MethodDelete, "/api/v1/article/"+articles[0].Slug, nil)
	req3.Header.Set("Authorization", actualResponse.Token)

	router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)
	expectedResponse := fmt.Sprintf(`{"message":"Article %s Deleted Successfully"}`, articles[0].Title)
	assert.Equal(t, expectedResponse, w3.Body.String())
}
