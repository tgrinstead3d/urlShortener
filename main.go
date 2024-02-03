package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Shortener struct {
	urls map[string]string
}

func main() {
	shortener := &Shortener{
		urls: make(map[string]string),
	}

	http.HandleFunc("/", handleForm)
	http.HandleFunc("/shorten", shortener.HandleShorten)
	http.HandleFunc("/short/", shortener.HandleRedirect)

	fmt.Println("URL Shortener is running on :8080")
	http.ListenAndServe(":8080", nil)
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/shorten", http.StatusSeeOther)
		return
	}

	http.ServeFile(w, r, "index.html")
}

func (us *Shortener) HandleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}

	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	shortKey := generateShortKey()
	us.urls[shortKey] = originalURL

	shortenedURL := fmt.Sprintf("http://localhost:8080/short/%s", shortKey)


	
	w.Write([]byte(shortenedURL))

}

func (us *Shortener) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path[len("/short/"):]
	if shortKey == "" {
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		return
	}

	originalURL, found := us.urls[shortKey]
	if !found {
		http.Error(w, "Shortened key not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
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
