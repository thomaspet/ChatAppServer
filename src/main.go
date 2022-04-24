package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gocraft/web"
	_ "github.com/heroku/x/hmetrics/onload"
)

var apiToken string

func main() {
	defer SseServer.Shutdown()

	// err := godotenv.Load()
	// if err != nil {
	// 	fmt.Println("Error loading .env file")
	// }

	apiToken := os.Getenv("API_TOKEN")
	apiToken = "testtoken"
	if apiToken == "" {
		log.Fatalln("API_TOKEN not set")
		return
	}

	server := web.New(HttpContext{})

	server.Middleware(web.LoggerMiddleware)
	server.Middleware(web.ShowErrorsMiddleware)

	server.Middleware(func(c *HttpContext, rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
		if strings.HasPrefix(req.URL.Path, "/events/messages") {
			rw.Header().Set("Access-Control-Allow-Origin", "*")
			rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			rw.Header().Set("Content-Type", "text/event-stream")
			rw.Header().Set("Cache-Control", "no-cache")
		}
		next(rw, req)
	})
	server.Middleware(func(c *HttpContext, rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
		// if route is events/messages skip validation
		// TODO - Implement auth for subscribing?
		if strings.HasPrefix(req.URL.Path, "/events/messages") {
			next(rw, req)
			return
		}

		authHeader := req.Header.Get("Authorization")
		apiTokenHeader := req.Header.Get("X-Api-Token")

		if authHeader != "" {
			// Get token from bearer token
			authHeader = strings.TrimPrefix(authHeader, "Bearer ")

			// Split token on .
			authHeaderParts := strings.Split(authHeader, ".")
			if len(authHeaderParts) != 3 {
				log.Println("Invalid auth header")
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			// base64 decode 2nd part of tokenparts
			decoded, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(authHeaderParts[1])
			if err != nil {
				log.Println("Error decoding auth header", err)
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			// json decode decoded
			var decodedJson map[string]interface{}
			err = json.Unmarshal(decoded, &decodedJson)
			if err != nil {
				log.Println("Error decoding auth header", err)
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Get preferred_username, name or email from decoded json
			username, ok := decodedJson["name"].(string)
			if !ok {
				username, ok = decodedJson["preferred_username"].(string)
				if !ok {
					username, ok = decodedJson["email"].(string)
					if !ok {
						log.Println("Error decoding auth header")
						rw.WriteHeader(http.StatusUnauthorized)
						return
					}
				}
			}

			guid, ok := decodedJson["sub"].(string)
			if !ok {
				log.Println("Error decoding auth header")
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			c.Guid = guid
			c.UserName = username
		} else {
			// Write 401 response
			log.Println("Missing Auth header")
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("401 - Unauthorized"))
			return
		}

		if apiTokenHeader != apiToken {
			// Write 401 response
			log.Println("Invalid API token")
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("401 - Unauthorized"))
			return

		}

		next(rw, req)
	})
	server.Post("/messages", (*HttpContext).SendMessage)
	server.Get("/messages", (*HttpContext).GetMessages)
	server.Get("/events/:channel", func(rw web.ResponseWriter, req *web.Request) {
		SseServer.ServeHTTP(rw, req.Request)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		fmt.Println("HOST not set, defaulting to", port)
	}

	log.Println(http.ListenAndServe(":"+port, server))
}
