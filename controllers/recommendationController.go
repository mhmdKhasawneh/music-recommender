package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mhmdKhasawneh/musicrecommendationapp/initializers"
	"github.com/mhmdKhasawneh/musicrecommendationapp/models"
	"gorm.io/gorm"
)

type Recommendation struct {
	To_user string `json:"to_user"`
	URL     string `json:"url"`
}

type SpotifyResponse struct {
	Name  string `json:"name"`
	Album struct {
		Images []map[string]interface{} `json:"images"`
	} `json:"album"`
}

func Recommend(c *gin.Context) {
	//get recommendee from context
	fromUser := c.GetString("email")

	var recommendation Recommendation
	if c.BindJSON(&recommendation) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to bind recommendation",
		})
		return
	}

	if !isUrlSpotify(recommendation.URL) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "URL must be a Spotify url.",
		})
		return
	}

	trackId := extractSpotifyId(recommendation.URL)
	SpotifyResponse := extractTrackMetaData(trackId)

	recommendationModel := models.Recommendation{
		To_user:   recommendation.To_user,
		From_user: fromUser,
		Url:       recommendation.URL,
		Name:      SpotifyResponse.Name,
		ImgUrl:    SpotifyResponse.Album.Images[2]["url"].(string),
	}

	result := initializers.Db.Create(&recommendationModel)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not push the recommendation to db.",
		})
		return
	}

	c.JSON(http.StatusOK, recommendationModel)
}

func GetRecommendations(c *gin.Context) {
	userEmail := c.GetString("email")

	var recommendationModels []models.Recommendation
	result := initializers.Db.Find(&recommendationModels, "to_user=?", userEmail)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not fetch recommendations for user.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recommendations": recommendationModels,
	})
}

func isUrlSpotify(url string) bool {
	return strings.Contains(url, "https://open.spotify.com/track/")
}

func extractSpotifyId(url string) string {
	from := url[31:]
	length := len(from)
	id := ""
	for i := 0; i < length; i++ {
		if from[i] == '?' {
			break
		}
		id = id + string(from[i])
	}
	return id
}

func extractTrackMetaData(track string) SpotifyResponse {
	client := &http.Client{}
	apiReq := "https://api.spotify.com/v1/tracks/" + track
	req, err := http.NewRequest("GET", apiReq, nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer BQB5OwAiblqcBiEODrcZLGZ7R8coZu1ujSJco40yiBaVsnL6EGiOxbic0T67_ZtxDXNmHlJzf3jpUhZmAa3QeKtNhcBrmgvCEi5KJ9IdwxVmhiTWH8byCKkmphpaxTqn7bqbT4HemaxlwPxy-t8GHF9Su0VNU-GxIoHKmRM33rd3DH4AOd2qQYnWTA40F5XP")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var responseObject SpotifyResponse
	json.Unmarshal(bodyBytes, &responseObject)
	return responseObject
}
