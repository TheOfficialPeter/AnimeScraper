package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type AnimeResponse struct {
	Data []struct {
		MalID    int    `json:"mal_id"`
		Title    string `json:"title"`
		TitleEnglish string `json:"title_english"`
		Images   struct {
			JPG struct {
				ImageURL string `json:"image_url"`
			} `json:"jpg"`
		} `json:"images"`
		Score      float64 `json:"score"`
		Status     string  `json:"status"`
		Episodes   int     `json:"episodes"`
		Synopsis   string  `json:"synopsis"`
		Year       int     `json:"year"`
		Genres     []Genre `json:"genres"`
		Season     string  `json:"season"`
		Popularity int     `json:"popularity"`
		Trailer    struct {
			YoutubeID string `json:"youtube_id"`
			URL       string `json:"url"`
		} `json:"trailer"`
		Studios []struct {
			Name string `json:"name"`
		} `json:"studios"`
		Duration string `json:"duration"`
		Rating   string `json:"rating"`
	} `json:"data"`
}

// Updated to match AniList API response format
type RecommendationsResponse struct {
	Data struct {
		Media struct {
			Recommendations struct {
				Nodes []struct {
					MediaRecommendation struct {
						ID    int `json:"id"`
						Title struct {
							Romaji  string `json:"romaji"`
							English string `json:"english"`
							Native  string `json:"native"`
						} `json:"title"`
						CoverImage struct {
							Medium     string `json:"medium"`
							Large      string `json:"large"`
							ExtraLarge string `json:"extraLarge"`
						} `json:"coverImage"`
						AverageScore float64 `json:"averageScore"`
					} `json:"mediaRecommendation"`
					User struct {
						Name string `json:"name"`
					} `json:"user"`
				} `json:"nodes"`
			} `json:"recommendations"`
		} `json:"Media"`
	} `json:"data"`
}

// AniList API response for anime details
type AniListAnimeResponse struct {
	Data struct {
		Media struct {
			ID    int `json:"id"`
			Title struct {
				Romaji  string `json:"romaji"`
				English string `json:"english"`
				Native  string `json:"native"`
			} `json:"title"`
			CoverImage struct {
				Medium     string `json:"medium"`
				Large      string `json:"large"`
				ExtraLarge string `json:"extraLarge"`
			} `json:"coverImage"`
			AverageScore float64 `json:"averageScore"`
		} `json:"Media"`
	} `json:"data"`
}

type Genre struct {
	MalID int    `json:"mal_id"`
	Name  string `json:"name"`
}

var client = resty.New().
	SetBaseURL("https://api.jikan.moe/v4").
	SetTimeout(10 * time.Second)

// AniList GraphQL API endpoint
const aniListURL = "https://graphql.anilist.co"

func main() {
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/recommendations", func(c *gin.Context) {
		c.HTML(http.StatusOK, "recommendations.html", nil)
	})

	// API endpoints
	api := r.Group("/api")
	{
		api.GET("/trending", getTrendingAnime)
		api.GET("/upcoming", getUpcomingAnime)
		api.GET("/search", searchAnime)
		api.GET("/isekai", getIsekaiAnime)
		api.GET("/ecchi", getEcchiAnime)
		api.GET("/romance", getRomanceAnime)
		api.GET("/suggestions", getAnimeSuggestions)
		api.GET("/recommendations/:id", getAnimeRecommendations)
		api.GET("/anime/:id", getAnimeDetails)
	}

	log.Fatal(r.Run(":8080"))
}

func getTrendingAnime(c *gin.Context) {
	var response AnimeResponse
	resp, err := client.R().
		SetResult(&response).
		Get("/top/anime")

	if err != nil || resp.IsError() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trending anime"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func getUpcomingAnime(c *gin.Context) {
	var response AnimeResponse
	resp, err := client.R().
		SetResult(&response).
		Get("/seasons/upcoming")

	if err != nil || resp.IsError() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch upcoming anime"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func searchAnime(c *gin.Context) {
	query := c.Query("q")
	genre := c.Query("genre")
	season := c.Query("season")
	year := c.Query("year")
	status := c.Query("status")
	order := c.Query("order_by")

	var response AnimeResponse
	req := client.R().SetResult(&response)

	// Build query parameters
	params := make(map[string]string)
	if query != "" {
		params["q"] = query
	}
	if genre != "" {
		params["genres"] = genre
	}
	if season != "" {
		params["season"] = season
	}
	if year != "" {
		params["year"] = year
	}
	if status != "" {
		params["status"] = status
	}
	if order != "" {
		params["order_by"] = order
		params["sort"] = "desc"
	}

	resp, err := req.SetQueryParams(params).Get("/anime")

	if err != nil || resp.IsError() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search anime"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func getIsekaiAnime(c *gin.Context) {
	var response AnimeResponse
	// Using genre ID 62 for Isekai (which is a sub-category of Fantasy)
	resp, err := client.R().
		SetResult(&response).
		SetQueryParam("genres", "62").
		SetQueryParam("order_by", "popularity").
		SetQueryParam("sort", "desc").
		SetQueryParam("limit", "24").
		Get("/anime")

	if err != nil || resp.IsError() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch isekai anime"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func getEcchiAnime(c *gin.Context) {
	var response AnimeResponse
	// Using genre ID 9 for Ecchi
	resp, err := client.R().
		SetResult(&response).
		SetQueryParam("genres", "9").
		SetQueryParam("order_by", "popularity").
		SetQueryParam("sort", "desc").
		SetQueryParam("limit", "24").
		Get("/anime")

	if err != nil || resp.IsError() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ecchi anime"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func getRomanceAnime(c *gin.Context) {
	var response AnimeResponse
	// Using genre ID 22 for Romance
	resp, err := client.R().
		SetResult(&response).
		SetQueryParam("genres", "22").
		SetQueryParam("order_by", "popularity").
		SetQueryParam("sort", "desc").
		SetQueryParam("limit", "24").
		Get("/anime")

	if err != nil || resp.IsError() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch romance anime"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func getAnimeSuggestions(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter is required"})
		return
	}

	// GraphQL query for anime search
	graphqlQuery := `
	query ($search: String) {
		Page(page: 1, perPage: 5) {
			media(search: $search, type: ANIME, sort: POPULARITY_DESC) {
				id
				title {
					romaji
					english
					native
				}
				coverImage {
					medium
					large
					extraLarge
				}
				averageScore
			}
		}
	}
	`

	// Create request variables
	variables := map[string]interface{}{
		"search": query,
	}

	// Create request body
	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     graphqlQuery,
		"variables": variables,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Make request to AniList API
	resp, err := http.Post(aniListURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch anime suggestions"})
		return
	}
	defer resp.Body.Close()

	// Parse response
	var response struct {
		Data struct {
			Page struct {
				Media []struct {
					ID    int `json:"id"`
					Title struct {
						Romaji  string `json:"romaji"`
						English string `json:"english"`
						Native  string `json:"native"`
					} `json:"title"`
					CoverImage struct {
						Medium     string `json:"medium"`
						Large      string `json:"large"`
						ExtraLarge string `json:"extraLarge"`
					} `json:"coverImage"`
					AverageScore float64 `json:"averageScore"`
				} `json:"media"`
			} `json:"Page"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse anime suggestions"})
		return
	}

	// Convert to format expected by frontend
	var convertedResponse AnimeResponse
	convertedResponse.Data = make([]struct {
		MalID    int    `json:"mal_id"`
		Title    string `json:"title"`
		TitleEnglish string `json:"title_english"`
		Images   struct {
			JPG struct {
				ImageURL string `json:"image_url"`
			} `json:"jpg"`
		} `json:"images"`
		Score      float64 `json:"score"`
		Status     string  `json:"status"`
		Episodes   int     `json:"episodes"`
		Synopsis   string  `json:"synopsis"`
		Year       int     `json:"year"`
		Genres     []Genre `json:"genres"`
		Season     string  `json:"season"`
		Popularity int     `json:"popularity"`
		Trailer    struct {
			YoutubeID string `json:"youtube_id"`
			URL       string `json:"url"`
		} `json:"trailer"`
		Studios []struct {
			Name string `json:"name"`
		} `json:"studios"`
		Duration string `json:"duration"`
		Rating   string `json:"rating"`
	}, len(response.Data.Page.Media))

	for i, anime := range response.Data.Page.Media {
		convertedResponse.Data[i].MalID = anime.ID
		convertedResponse.Data[i].Title = anime.Title.Romaji
		convertedResponse.Data[i].TitleEnglish = anime.Title.English
		if convertedResponse.Data[i].TitleEnglish == "" {
			convertedResponse.Data[i].TitleEnglish = anime.Title.Romaji
		}
		// Use extraLarge image if available, fallback to large, then medium
		if anime.CoverImage.ExtraLarge != "" {
			convertedResponse.Data[i].Images.JPG.ImageURL = anime.CoverImage.ExtraLarge
		} else if anime.CoverImage.Large != "" {
			convertedResponse.Data[i].Images.JPG.ImageURL = anime.CoverImage.Large
		} else {
			convertedResponse.Data[i].Images.JPG.ImageURL = anime.CoverImage.Medium
		}
		convertedResponse.Data[i].Score = anime.AverageScore / 10
	}

	c.JSON(http.StatusOK, convertedResponse)
}

// Updated to use AniList API
func getAnimeRecommendations(c *gin.Context) {
	animeID := c.Param("id")
	if animeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anime ID is required"})
		return
	}

	// GraphQL query for recommendations
	query := `
	query ($id: Int) {
		Media(id: $id, type: ANIME) {
			recommendations(sort: RATING_DESC) {
				nodes {
					mediaRecommendation {
						id
						title {
							romaji
							english
							native
						}
						coverImage {
							medium
							large
							extraLarge
						}
						averageScore
					}
					user {
						name
					}
				}
			}
		}
	}
	`

	// Create request variables
	variables := map[string]interface{}{
		"id": animeID,
	}

	// Create request body
	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Make request to AniList API
	resp, err := http.Post(aniListURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recommendations"})
		return
	}
	defer resp.Body.Close()

	// Parse response
	var response RecommendationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse recommendations"})
		return
	}

	// Convert to format expected by frontend
	var convertedResponse struct {
		Data []struct {
			Entry struct {
				MalID       int     `json:"mal_id"`
				Title       string  `json:"title"`
				TitleEnglish string  `json:"title_english"`
				Score       float64 `json:"score"`
				Images      struct {
					JPG struct {
						ImageURL string `json:"image_url"`
					} `json:"jpg"`
				} `json:"images"`
			} `json:"entry"`
			Content string `json:"content"`
			User    struct {
				Username string `json:"username"`
			} `json:"user"`
		} `json:"data"`
	}

	// Map AniList response to the format expected by the frontend
	convertedResponse.Data = make([]struct {
		Entry struct {
			MalID       int     `json:"mal_id"`
			Title       string  `json:"title"`
			TitleEnglish string  `json:"title_english"`
			Score       float64 `json:"score"`
			Images      struct {
				JPG struct {
					ImageURL string `json:"image_url"`
				} `json:"jpg"`
			} `json:"images"`
		} `json:"entry"`
		Content string `json:"content"`
		User    struct {
			Username string `json:"username"`
		} `json:"user"`
	}, len(response.Data.Media.Recommendations.Nodes))

	for i, rec := range response.Data.Media.Recommendations.Nodes {
		convertedResponse.Data[i].Entry.MalID = rec.MediaRecommendation.ID
		convertedResponse.Data[i].Entry.Title = rec.MediaRecommendation.Title.Romaji
		convertedResponse.Data[i].Entry.TitleEnglish = rec.MediaRecommendation.Title.English
		if convertedResponse.Data[i].Entry.TitleEnglish == "" {
			convertedResponse.Data[i].Entry.TitleEnglish = rec.MediaRecommendation.Title.Romaji
		}
		// Use extraLarge image if available, fallback to large, then medium
		if rec.MediaRecommendation.CoverImage.ExtraLarge != "" {
			convertedResponse.Data[i].Entry.Images.JPG.ImageURL = rec.MediaRecommendation.CoverImage.ExtraLarge
		} else if rec.MediaRecommendation.CoverImage.Large != "" {
			convertedResponse.Data[i].Entry.Images.JPG.ImageURL = rec.MediaRecommendation.CoverImage.Large
		} else {
			convertedResponse.Data[i].Entry.Images.JPG.ImageURL = rec.MediaRecommendation.CoverImage.Medium
		}
		// Set the score in the Entry structure
		convertedResponse.Data[i].Entry.Score = rec.MediaRecommendation.AverageScore / 10
		convertedResponse.Data[i].Content = "Score: " + fmt.Sprintf("%.2f", rec.MediaRecommendation.AverageScore/10)
		convertedResponse.Data[i].User.Username = rec.User.Name
		if convertedResponse.Data[i].User.Username == "" {
			convertedResponse.Data[i].User.Username = "AniList User"
		}
	}

	c.JSON(http.StatusOK, convertedResponse)
}

// Updated to use AniList API
func getAnimeDetails(c *gin.Context) {
	animeID := c.Param("id")
	if animeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anime ID is required"})
		return
	}

	// GraphQL query for anime details
	query := `
	query ($id: Int) {
		Media(id: $id, type: ANIME) {
			id
			title {
				romaji
				english
				native
			}
			coverImage {
				medium
				large
				extraLarge
			}
			averageScore
		}
	}
	`

	// Create request variables
	variables := map[string]interface{}{
		"id": animeID,
	}

	// Create request body
	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Make request to AniList API
	resp, err := http.Post(aniListURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch anime details"})
		return
	}
	defer resp.Body.Close()

	// Parse response
	var response AniListAnimeResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse anime details"})
		return
	}

	// Convert to format expected by frontend
	var convertedResponse struct {
		Data struct {
			MalID       int     `json:"mal_id"`
			Title       string  `json:"title"`
			TitleEnglish string  `json:"title_english"`
			Score       float64 `json:"score"`
			Images      struct {
				JPG struct {
					ImageURL string `json:"image_url"`
				} `json:"jpg"`
			} `json:"images"`
		} `json:"data"`
	}

	convertedResponse.Data.MalID = response.Data.Media.ID
	convertedResponse.Data.Title = response.Data.Media.Title.Romaji
	convertedResponse.Data.TitleEnglish = response.Data.Media.Title.English
	if convertedResponse.Data.TitleEnglish == "" {
		convertedResponse.Data.TitleEnglish = response.Data.Media.Title.Romaji
	}
	convertedResponse.Data.Score = response.Data.Media.AverageScore / 10
	// Use extraLarge image if available, fallback to large, then medium
	if response.Data.Media.CoverImage.ExtraLarge != "" {
		convertedResponse.Data.Images.JPG.ImageURL = response.Data.Media.CoverImage.ExtraLarge
	} else if response.Data.Media.CoverImage.Large != "" {
		convertedResponse.Data.Images.JPG.ImageURL = response.Data.Media.CoverImage.Large
	} else {
		convertedResponse.Data.Images.JPG.ImageURL = response.Data.Media.CoverImage.Medium
	}

	c.JSON(http.StatusOK, convertedResponse)
} 