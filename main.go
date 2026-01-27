package main

import (
	"log"
	"os"

	"github.com/colin-harness/go-example/handler"
	"github.com/colin-harness/go-example/store"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	s := store.NewMemoryStore()
	h := handler.NewHandler(s)

	r := setupRouter(h)

	log.Printf("Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func setupRouter(h *handler.Handler) *gin.Engine {
	r := gin.Default()

	r.Static("/static", "./static")
	r.GET("/", h.Index)
	r.GET("/paste/:id", h.GetPaste)
	r.POST("/paste", h.CreatePaste)
	r.GET("/api/paste/:id", h.GetPasteJSON)
	r.DELETE("/api/paste/:id", h.DeletePaste)

	return r
}
