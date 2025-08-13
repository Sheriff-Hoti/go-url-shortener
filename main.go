package main

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/Sheriff-Hoti/go-url-shortener/database"
	templates "github.com/Sheriff-Hoti/go-url-shortener/templates"
	templ "github.com/a-h/templ"
	_ "github.com/glebarez/go-sqlite"
)

const baseUrl = "http://localhost:3000/"

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// stolen from: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func RandStringBytesMaskImprSrcSB(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

var urlRegex = regexp.MustCompile(`^(https?://)([^\s/$.?#].[^\s]*)$`)

//another alternative is to read the schema file

//go:embed schema.sql
var ddl string

func main() {
	ctx := context.Background()
	//another alternative is to put a file in place of the :memory: and the database will be created in that file and the data will persist
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		fmt.Println(err)
		return
	}

	queries := database.New(db)

	router := http.NewServeMux()

	urlForm := templates.UrlForm()

	router.Handle("/", templ.Handler(urlForm))
	router.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	router.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {

		res, err := queries.GetUrlByShortenedUrl(ctx, r.PathValue("id"))

		if errors.Is(err, sql.ErrNoRows) {
			// URL not found return 404
			w.WriteHeader(http.StatusNotFound)
			templates.Error404().Render(r.Context(), w)
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			templates.Error404().Render(r.Context(), w)
			return
		}

		http.Redirect(w, r, res.OriginalUrl, http.StatusFound)
	})
	router.HandleFunc("POST /shorten", func(w http.ResponseWriter, r *http.Request) {

		// Parse form data
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		originalURL := r.FormValue("url")
		if originalURL == "" {
			w.WriteHeader(http.StatusBadRequest)
			templates.Result("Missing url", true).Render(r.Context(), w)
			return
		}
		if _, err := url.ParseRequestURI(originalURL); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			templates.Result("Invalid URL format. Please use http:// or https://", true).Render(r.Context(), w)
			return
		}

		res, err := queries.GetUrl(ctx, originalURL)

		if errors.Is(err, sql.ErrNoRows) {
			// URL not found, proceed to create a new shortened URL
			shortenedURL := RandStringBytesMaskImprSrcSB(5) // Example shortening logic
			_, err = queries.CreateUrl(ctx, database.CreateUrlParams{
				OriginalUrl:  originalURL,
				ShortenedUrl: shortenedURL,
			})
			msg := fmt.Sprintf("URL '%s%s' shortened successfully!", baseUrl, shortenedURL)
			templates.Result(msg, false).Render(r.Context(), w)
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, "Failed to create shortened URL", http.StatusInternalServerError)
			return
		}

		if res.OriginalUrl == originalURL {
			// URL already exists, return the existing shortened URL
			msg := fmt.Sprintf("URL already shortened: '%s%s'", baseUrl, res.ShortenedUrl)
			templates.Result(msg, false).Render(r.Context(), w)
			return
		}

		// Here you would handle the URL shortening logic
		// For now, just respond with a success message

		// Return HTML fragment to replace #result
	})

	server := http.Server{
		Addr:    ":3000",
		Handler: router,
	}
	fmt.Println("Listening on :3000")

	server.ListenAndServe()
	defer server.Close()

}
