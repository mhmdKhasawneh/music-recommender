package controllers

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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Signup(c *gin.Context) {
	var user User
	if c.BindJSON(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to bind user",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	userModel := models.User{
		Email:    user.Email,
		Password: string(hash),
	}

	result := initializers.Db.Create(&userModel)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to push user to database",
		})
		return
	}

	userModel.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"user": userModel,
	})

}

func Login(c *gin.Context) {
	var err error
	var user User
	if c.BindJSON(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to bind user",
		})
		return
	}

	var retrievedUser models.User
	result := initializers.Db.First(&retrievedUser, "email=?", user.Email)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(retrievedUser.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	var tokenString string
	tokenString, err = token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not create token",
		})
		return
	}

	//c.SetSameSite(http.SameSiteLaxMode)
	//c.SetCookie("Authorization", tokenString, 3600, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func TokenLogin(c *gin.Context) {
	// get token string from cookie
	bearerToken := c.Request.Header["Authorization"]
	if bearerToken[0] == "null" {
		c.Status(401)
		return
	}
	tokenString := strings.Split(bearerToken[0], " ")[1]

	//validate and decode the token
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			c.Status(401)
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//check expiry
		if expired := float64(time.Now().Unix()) > claims["exp"].(float64); expired {
			c.Status(401)
			return
		}
		//find user
		var user models.User
		result := initializers.Db.First(&user, "email=?", claims["sub"])
		if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.Error != nil {
			c.Status(401)
			return
		}
		c.Status(200)
	} else {
		c.Status(401)
		return
	}
}
