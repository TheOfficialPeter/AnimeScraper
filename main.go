package main

import (
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

type RecommendationsResponse struct {
	Data []struct {
		Entry struct {
			MalID       int    `json:"mal_id"`
			Title       string `json:"title"`
			TitleEnglish string `json:"title_english"`
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

type Genre struct {
	MalID int    `json:"mal_id"`
	Name  string `json:"name"`
}

var client = resty.New().
	SetBaseURL("https://api.jikan.moe/v4").
	SetTimeout(10 * time.Second)

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

	var response AnimeResponse
	resp, err := client.R().
		SetResult(&response).
		SetQueryParam("q", query).
		SetQueryParam("limit", "5").
		Get("/anime")

	if err != nil || resp.IsError() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch anime suggestions"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func getAnimeRecommendations(c *gin.Context) {
	animeID := c.Param("id")
	if animeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anime ID is required"})
		return
	}

	var response RecommendationsResponse
	resp, err := client.R().
		SetResult(&response).
		SetQueryParam("fields", "title_english").
		Get("/anime/" + animeID + "/recommendations")

	if err != nil || resp.IsError() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recommendations"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func getAnimeDetails(c *gin.Context) {
	animeID := c.Param("id")
	if animeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anime ID is required"})
		return
	}

	var response struct {
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

	resp, err := client.R().
		SetResult(&response).
		Get("/anime/" + animeID)

	if err != nil || resp.IsError() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch anime details"})
		return
	}

	c.JSON(http.StatusOK, response)
} 