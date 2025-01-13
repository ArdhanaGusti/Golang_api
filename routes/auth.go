package routes

import (
	"net/http"
	"os"
	"time"

	"github.com/ArdhanaGusti/Golang_api/config"
	"github.com/ArdhanaGusti/Golang_api/handler/failed"
	"github.com/ArdhanaGusti/Golang_api/handler/validation"
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
		c.JSON(500, failed.FailedResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
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
		c.JSON(500, failed.FailedResponse{
			StatusCode: 500,
			Message:    err.Error(),
		})
		return
	}

	newUser, ers := getOrRegisterUser(provider, (*structs.User)(user))
	if ers != nil {
		c.JSON(500, failed.FailedResponse{
			StatusCode: 500,
			Message:    ers.Error(),
		})
	}
	jwtToken, ert := getToken(&newUser)
	if ert != nil {
		c.JSON(500, failed.FailedResponse{
			StatusCode: 500,
			Message:    ert.Error(),
		})
	}
	c.JSON(200, gin.H{
		"data":    newUser,
		"token":   jwtToken,
		"message": "Berhasil",
	})
}

func getOrRegisterUser(provider string, user *structs.User) (models.User, error) {
	var userData models.User

	if err := config.DB.Where("provider = ? AND social_id = ?", provider, user.ID).First(&userData).Error; err != nil {
		return models.User{}, err
	}

	if userData.ID == 0 {
		newUser := models.User{
			Username: user.Username,
			Fullname: user.FullName,
			Email:    user.Email,
			SocialID: user.ID,
			Provider: provider,
			Avatar:   user.Avatar,
		}
		if err := config.DB.Create(&newUser).Error; err != nil {
			return models.User{}, err
		}
		return newUser, nil
	} else {
		return userData, nil
	}
}

func getToken(user *models.User) (string, error) {
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"user_role": user.Role,
		"exp":       time.Now().AddDate(0, 0, 7).Unix(),
		"iat":       time.Now().Unix(),
	})

	tokenString, err := newToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CheckToken(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Berhasil",
	})
}

func RegisterUser(c *gin.Context) {
	var userPayload validation.RegisterUserPayload

	if err := c.ShouldBind(&userPayload); err != nil {
		c.JSON(400, failed.FailedResponse{
			StatusCode: 400,
			Message:    err.Error(),
		})
		return
	}

	var existedUser models.User
	if err := config.DB.First(&existedUser, "email = ?", userPayload.Email).Error; err == nil {
		c.JSON(409, failed.FailedResponse{
			StatusCode: 409,
			Message:    "User is exist",
		})
		c.Abort()
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(userPayload.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(400, failed.FailedResponse{
			StatusCode: 400,
			Message:    "Hashing is failed",
		})
		c.Abort()
		return
	}

	newUser := models.User{
		Username: userPayload.Username,
		Fullname: userPayload.Fullname,
		Email:    userPayload.Email,
		Password: string(hash),
	}

	if err := config.DB.Create(&newUser).Error; err != nil {
		c.JSON(500, failed.FailedResponse{
			StatusCode: 500,
			Message:    err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"status": "berhasil",
		"data":   newUser,
	})
}

func LoginUser(c *gin.Context) {
	var userPayload validation.LoginUserPayload

	if err := c.ShouldBind(&userPayload); err != nil {
		c.JSON(400, failed.FailedResponse{
			StatusCode: 400,
			Message:    err.Error(),
		})
		return
	}

	var existedUser models.User
	if err := config.DB.First(&existedUser, "email = ?", userPayload.Email).Error; err != nil {
		c.JSON(404, failed.FailedResponse{
			StatusCode: 404,
			Message:    "User don't exist",
		})
		c.Abort()
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(userPayload.Password)); err != nil {
		c.JSON(400, failed.FailedResponse{
			StatusCode: 400,
			Message:    err.Error(),
		})
		c.Abort()
		return
	}

	jwtToken, ert := getToken(&existedUser)
	if ert != nil {
		c.JSON(500, failed.FailedResponse{
			StatusCode: 500,
			Message:    ert.Error(),
		})
	}
	c.JSON(200, gin.H{
		"email":   existedUser.Email,
		"token":   jwtToken,
		"message": "Berhasil",
	})
}

func ChangeEmail(c *gin.Context) {
	var existedUser models.User
	if err := config.DB.First(&existedUser, "id = ?", uint(c.MustGet("jwt_user_id").(float64))).Error; err != nil {
		c.JSON(404, failed.FailedResponse{
			StatusCode: 404,
			Message:    "User don't exist",
		})
		c.Abort()
		return
	}

	var email = c.PostForm("Email")
	if email == "" {
		c.JSON(404, failed.FailedResponse{
			StatusCode: 404,
			Message:    "Email can't be empty",
		})
		return
	}

	if err := config.DB.Model(&existedUser).Where("id = ?", uint(c.MustGet("jwt_user_id").(float64))).Updates(models.User{Email: email}).Error; err != nil {
		c.JSON(500, failed.FailedResponse{
			StatusCode: 500,
			Message:    err.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"email":   existedUser.Email,
		"message": "Berhasil",
	})
}

func GetProfile(c *gin.Context) {
	var user models.User
	user_id := uint(c.MustGet("jwt_user_id").(float64))

	if err := config.DB.Where("id = ?", user_id).Preload("Articles", "user_id = ?", user_id).Find(&user).Error; err != nil {
		c.JSON(404, failed.FailedResponse{
			StatusCode: 404,
			Message:    err.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"data": user,
	})
}
