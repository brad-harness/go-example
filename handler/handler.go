package handler

import (
	"net/http"
	"time"

	"github.com/colin-harness/go-example/store"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	store store.Store
}

func NewHandler(s store.Store) *Handler {
	return &Handler{store: s}
}

func (h *Handler) Index(c *gin.Context) {
	c.File("./static/index.html")
}

type CreatePasteRequest struct {
	Content string `json:"content" binding:"required"`
	TTL     int    `json:"ttl"`
}

func (h *Handler) CreatePaste(c *gin.Context) {
	var req CreatePasteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content is required"})
		return
	}

	ttl := time.Duration(req.TTL) * time.Second
	paste, err := h.store.Create(req.Content, ttl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, paste)
}

func (h *Handler) GetPaste(c *gin.Context) {
	id := c.Param("id")
	paste, err := h.store.Get(id)
	if err == store.ErrNotFound {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Paste not found"})
		return
	}
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "paste.html", gin.H{"paste": paste})
}

func (h *Handler) GetPasteJSON(c *gin.Context) {
	id := c.Param("id")
	paste, err := h.store.Get(id)
	if err == store.ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, paste)
}

func (h *Handler) DeletePaste(c *gin.Context) {
	id := c.Param("id")
	err := h.store.Delete(id)
	if err == store.ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "paste deleted"})
}
