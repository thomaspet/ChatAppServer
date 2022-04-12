package main

import (
	"net/http"

	"github.com/gocraft/web"
)

func main() {
	defer SseServer.Shutdown()
	server := web.New(HttpContext{})
	server.Middleware(func(c *HttpContext, rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		rw.Header().Set("Content-Type", "text/event-stream")
		rw.Header().Set("Cache-Control", "no-cache")
		next(rw, req)
	})
	server.Post("/messages", (*HttpContext).SendMessage)
	server.Get("/messages", (*HttpContext).GetMessages)
	server.Get("/events/:channel", func(rw web.ResponseWriter, req *web.Request) {
		SseServer.ServeHTTP(rw, req.Request)
	})

	http.ListenAndServe("localhost:8080", server)
}
