package handler

import (
	"net/http"

	"github.com/Neel-shetty/clarity/internal/model"
	"github.com/Neel-shetty/clarity/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
)

type QuizHandler struct {
	quizService service.QuizService
	noteService service.NoteService
}

func NewQuizHandler(quizService service.QuizService, noteService service.NoteService) *QuizHandler {
	return &QuizHandler{
		quizService,
		noteService,
	}
}

func (h *QuizHandler) CreateQuiz(c *gin.Context) {
	uid, exists := c.Get("userID")
	userID, err := uuid.FromString(uid.(string))
	if !exists || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user not known"})
		return
	}

	var reqBody struct {
		Title   string   `json:"title" binding:"required"`
		NoteIDs []string `json:"noteIds" binding:"required"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	var notes []*model.Note
	for _, idStr := range reqBody.NoteIDs {
		noteID, err := uuid.FromString(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID format"})
			return
		}

		note, err := h.noteService.GetNoteByID(c.Request.Context(), noteID, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Note not found or unauthorized", "noteId": idStr})
			return
		}
		notes = append(notes, note)
	}

	quiz, err := h.quizService.CreateQuiz(c.Request.Context(), userID, reqBody.Title, notes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "unable to create quiz", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "quiz created successfully", "quiz": quiz})
}
