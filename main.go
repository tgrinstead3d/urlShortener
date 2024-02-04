package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Shortener struct {
	urls map[string]string
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	shortener := &Shortener{
		urls: make(map[string]string),
	}

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.POST("/shorten", shortener.HandleShorten)
	r.GET("/shortened", shortener.HandleShorten)
	r.GET("/short/:key", shortener.HandleRedirect)

	fmt.Println("URL Shortener app is running on :8080")
	r.Run(":8080")
}

func (us *Shortener) HandleShorten(c *gin.Context) {
	originalURL := c.PostForm("url")
	if originalURL == "" {
		c.String(http.StatusBadRequest, "URL parameter is missing")
		return
	}

	shortKey := generateShortKey()
	us.urls[shortKey] = originalURL

	shortenedURL := fmt.Sprintf("localhost:8080/short/%s", shortKey)
	fmt.Printf("New shortened URL is %s", shortenedURL)
	c.String(http.StatusOK, shortenedURL)
}

func (us *Shortener) HandleRedirect(c *gin.Context) {
	shortKey := c.Param("key")
	originalURL, found := us.urls[shortKey]
	if !found {
		c.String(http.StatusNotFound, "Shortened key not found")
		return
	}
	c.Redirect(http.StatusMovedPermanently, originalURL)
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[random.Intn(len(charset))]
	}
	return string(shortKey)
}
