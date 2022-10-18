package middleware

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func IsAuth() gin.HandlerFunc {
	return CheckJwt(false)
}

func IsAdmin() gin.HandlerFunc {
	return CheckJwt(true)
}

func CheckJwt(admin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		// bearerToken := strings.Split(authHeader, " ")
		token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println(claims["user_id"], claims["user_role"])
			userRole := bool(claims["user_role"].(bool))
			c.Set("jwt_user_id", claims["user_id"])

			if admin == true && userRole == false {
				c.JSON(404, gin.H{"msg": "You are not admin"})
				c.Abort()
				return
			}
			// c.Set("jwt_user_role", claims["user_role"])
		} else {
			c.JSON(422, gin.H{"msg": "Invalid token", "error": err})
			c.Abort()
			return
		}
	}
}
