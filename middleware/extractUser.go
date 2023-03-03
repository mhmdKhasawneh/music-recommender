package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mhmdKhasawneh/musicrecommendationapp/initializers"
	"github.com/mhmdKhasawneh/musicrecommendationapp/models"
	"gorm.io/gorm"
)

func ExtractUserFromLocalStorage(c *gin.Context) {
	// get token string from header
	bearerToken := c.Request.Header["Authorization"]
	if bearerToken[0] == "null" {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	tokenString := strings.Split(bearerToken[0], " ")[1]
	
	//validate and decode the token
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
	
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//check expiry 
		if expired := float64(time.Now().Unix()) > claims["exp"].(float64); expired {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		//find user
		var user models.User
		result := initializers.Db.First(&user, "email=?", claims["sub"])
		if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.Error != nil{
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Set("email", user.Email)

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	//next
	c.Next()
}
