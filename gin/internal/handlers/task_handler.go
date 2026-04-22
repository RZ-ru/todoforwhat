package handlers

import (
	"net/http"
	"todo/gin/internal/models"
	"todo/gin/internal/services"

	"github.com/gin-gonic/gin"
)

type TaskRequest struct {
	Description string `json:"description"`
}

type TaskHandler struct {
	service *services.TaskService
}

func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req TaskRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		// 400
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if req.Description == "" {
		// 400
		c.JSON(http.StatusBadRequest, gin.H{"error": "description is required"})
		return
	}
	/*
		task := models.Task{
			Description: req.Description,
			Done:        false,
			CreatedAt:   time.Now(),
		}
	*/
	task, err := h.service.CreateTask(c.Request.Context(), req.Description)
	if err != nil {
		// 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}
	// 201
	c.JSON(http.StatusCreated, task)
}
func (h *TaskHandler) GetTasks(c *gin.Context) {
	tasks, err := h.service.GetTasks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}
func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	task, err := h.service.GetTaskByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}

	var req TaskRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		// 400
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if req.Description == "" {
		// 400
		c.JSON(http.StatusBadRequest, gin.H{"error": "description is required"})
		return
	}

	task, err := h.service.UpdateTask(
		c.Request.Context(),
		id,
		req.Description,
	)

	if err != nil {
		// 404
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to update task"})
		return
	}
	// 200
	c.JSON(http.StatusOK, task)
}
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, ok := getID(c)
	if !ok {
		return
	}
	err := h.service.DeleteTask(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
func getID(c *gin.Context) (models.ID, bool) {
	id := models.ID(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return "", false
	}
	return id, true
}
