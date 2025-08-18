package handler

import (
	"net/http"

	"github.com/Neel-shetty/clarity/internal/config"
	"github.com/Neel-shetty/clarity/internal/model"
	"github.com/Neel-shetty/clarity/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
)

type QuizHandler struct {
	quizService service.QuizService
	noteService service.NoteService
	cfg         config.Config
}

func NewQuizHandler(quizService service.QuizService, noteService service.NoteService, cfg config.Config) *QuizHandler {
	return &QuizHandler{
		quizService,
		noteService,
		cfg,
	}
}

func (h *QuizHandler) CreateQuiz(c *gin.Context) {
	uid, exists := c.Get("userID")
	userID, err := uuid.FromString(uid.(string))
	if !exists || err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"message": "user not known"})
		return
	}
	var reqBody map[string]any

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}
	title, ok := reqBody["title"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid 'title' field"})
		return
	}
	fileName, ok := reqBody["fileName"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid 'filename' field"})
		return
	}
	contentType, ok := reqBody["contentType"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid 'filename' field"})
		return
	}

	note, presignedUrl, err := h.noteService.CreateNote(c.Request.Context(), userID, fileName, contentType)
	var notes []*model.Note
	notes = append(notes, note)
	quiz, qerr := h.quizService.CreateQuiz(c.Request.Context(), userID, title, notes)
	if qerr != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"message": "unable to create quiz", "error": err.Error()})
	}

	c.JSON(http.StatusOK,
		gin.H{"message": "quiz created successfully", "upload": presignedUrl, "quiz": quiz})

}
