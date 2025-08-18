package handler

import (
	"net/http"

	"github.com/Neel-shetty/clarity/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
)

type NoteHandler struct {
	noteService service.NoteService
}

func NewNoteHandler(noteService service.NoteService) *NoteHandler {
	return &NoteHandler{
		noteService,
	}
}

func (h *NoteHandler) CreateNote(c *gin.Context) {
	uid, exists := c.Get("userID")
	userID, err := uuid.FromString(uid.(string))
	if !exists || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user not known"})
		return
	}

	var reqBody struct {
		FileName    string `json:"fileName" binding:"required"`
		ContentType string `json:"contentType" binding:"required"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	note, presignedUrl, err := h.noteService.CreateNote(c.Request.Context(), userID, reqBody.FileName, reqBody.ContentType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "unable to create note", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "note created successfully",
		"note":         note,
		"presignedUrl": presignedUrl,
	})
}

func (h *NoteHandler) GetNote(c *gin.Context) {
	uid, exists := c.Get("userID")
	userID, err := uuid.FromString(uid.(string))
	if !exists || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user not known"})
		return
	}

	noteIDStr := c.Param("id")
	noteID, err := uuid.FromString(noteIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID format"})
		return
	}

	note, err := h.noteService.GetNoteByID(c.Request.Context(), noteID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "note not found", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"note": note,
	})
}

func (h *NoteHandler) DeleteNote(c *gin.Context) {
	uid, exists := c.Get("userID")
	userID, err := uuid.FromString(uid.(string))
	if !exists || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user not known"})
		return
	}

	noteIDStr := c.Param("id")
	noteID, err := uuid.FromString(noteIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID format"})
		return
	}

	err = h.noteService.DeleteNote(c.Request.Context(), noteID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "unable to delete note", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "note deleted successfully",
	})
}
