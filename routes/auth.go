package routes

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ArdhanaGusti/Golang_api/config"
	"github.com/ArdhanaGusti/Golang_api/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/danilopolani/gocialite.v1/structs"
)

func RedirectHandler(c *gin.Context) {
	provider := c.Param("provider")

	providerSecrets := map[string]map[string]string{
		"github": {
			"clientID":     os.Getenv("CLIENT_ID_GH"),
			"clientSecret": os.Getenv("CLIENT_SECRET_GH"),
			"redirectURL":  os.Getenv("AUTH_REDIRECT_URL") + "/github/callback",
		},
		"google": {
			"clientID":     os.Getenv("CLIENT_ID_GO"),
			"clientSecret": os.Getenv("CLIENT_SECRET_GO"),
			"redirectURL":  os.Getenv("AUTH_REDIRECT_URL") + "/google/callback",
		},
	}

	providerScopes := map[string][]string{
		"github": []string{"public_repo"},
		"google": []string{},
	}

	providerData := providerSecrets[provider]
	actualScopes := providerScopes[provider]
	authURL, err := config.Gocial.New().
		Driver(provider).
		Scopes(actualScopes).
		Redirect(
			providerData["clientID"],
			providerData["clientSecret"],
			providerData["redirectURL"],
		)

	// Check for errors (usually driver not valid)
	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	// Redirect with authURL
	c.Redirect(http.StatusFound, authURL)
}

func CallbackHandler(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")
	provider := c.Param("provider")

	user, _, err := config.Gocial.Handle(state, code)
	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	var newUser = getOrRegisterUser(provider, (*structs.User)(user))
	var jwtToken = getToken(&newUser)
	c.JSON(200, gin.H{
		"data":    newUser,
		"token":   jwtToken,
		"message": "Berhasil",
	})
}

func getOrRegisterUser(provider string, user *structs.User) models.User {
	var userData models.User

	config.DB.Where("provider = ? AND social_id = ?", provider, user.ID).First(&userData)

	if userData.ID == 0 {
		newUser := models.User{
			Username: user.Username,
			Fullname: user.FullName,
			Email:    user.Email,
			SocialID: user.ID,
			Provider: provider,
			Avatar:   user.Avatar,
		}
		config.DB.Create(&newUser)
		return newUser
	} else {
		return userData
	}
}

func getToken(user *models.User) string {
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"user_role": user.Role,
		"exp":       time.Now().AddDate(0, 0, 7).Unix(),
		"iat":       time.Now().Unix(),
	})

	tokenString, err := newToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		fmt.Println("error", err)
	}

	return tokenString
}

func CheckToken(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Berhasil",
	})
}

func RegisterUser(c *gin.Context) {
	var existedUser models.User
	if err := config.DB.First(&existedUser, "email = ?", c.PostForm("Email")).Error; err == nil {
		c.JSON(409, gin.H{"status": "User is exist"})
		c.Abort()
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(c.PostForm("Password")), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(400, gin.H{"status": "Hashing failed"})
		c.Abort()
		return
	}

	newUser := models.User{
		Username: c.PostForm("Username"),
		Fullname: c.PostForm("Fullname"),
		Email:    c.PostForm("Email"),
		Password: string(hash),
	}

	config.DB.Create(&newUser)

	c.JSON(200, gin.H{
		"status": "berhasil",
		"data":   newUser,
	})
}

func LoginUser(c *gin.Context) {
	var existedUser models.User
	if err := config.DB.First(&existedUser, "email = ?", c.PostForm("Email")).Error; err != nil {
		c.JSON(404, gin.H{"status": "User don't exist"})
		c.Abort()
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(c.PostForm("Password"))); err != nil {
		c.JSON(400, gin.H{"status": err.Error()})
		c.Abort()
		return
	}

	var jwtToken = getToken(&existedUser)
	c.JSON(200, gin.H{
		"email":   existedUser.Email,
		"token":   jwtToken,
		"message": "Berhasil",
	})
}

func GetProfile(c *gin.Context) {
	var user models.User
	user_id := uint(c.MustGet("jwt_user_id").(float64))

	if err := config.DB.Where("id = ?", user_id).Preload("Articles", "user_id = ?", user_id).Find(&user).Error; err != nil {
		c.JSON(404, gin.H{"status": "error", "error": err})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"data": user,
	})
}
