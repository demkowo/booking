package handler

import (
	"net/http"
	"time"

	model "github.com/demkowo/booking/models"
	service "github.com/demkowo/booking/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Room interface {
	CreateTableRooms()

	Add(*gin.Context)
	Find(*gin.Context)
	FindAvailable(*gin.Context)
	GetById(*gin.Context)
	CheckIfAvailableById(*gin.Context)
	Update(*gin.Context)
}

type room struct {
	service service.Room
}

func NewRoom(service service.Room) Room {
	log.Trace()

	return &room{
		service: service,
	}
}

func (h *room) CreateTableRooms() {
	log.Trace()

	res := h.service.CreateTableRooms()
	log.Info(res)
}

func (h *room) Add(c *gin.Context) {
	log.Trace()

	var input struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Errorf("Failed to bind JSON input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	room := &model.Room{
		Id:      uuid.New(),
		Name:    input.Name,
		Created: time.Now(),
		Updated: time.Now(),
	}

	if err := h.service.Add(room); err != nil {
		log.Errorf("Failed to add room: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add room",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"room": room})
}

func (h *room) Find(c *gin.Context) {
	log.Trace()

	rooms, err := h.service.Find()
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "list rooms failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

func (h *room) FindAvailable(c *gin.Context) {
	log.Trace()

	var input struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	if e := c.ShouldBindJSON(&input); e != nil {
		log.Errorf("Failed to bind JSON input: %v", e)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	startDate, e := time.Parse("2006-01-02", input.StartDate)
	if e != nil {
		log.Error(e)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid start date",
		})
		return
	}

	endDate, e := time.Parse("2006-01-02", input.EndDate)
	if e != nil {
		log.Error(e)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid end date",
		})
		return
	}

	rooms, err := h.service.FindAvailable(startDate, endDate)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "list available rooms failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

func (h *room) GetById(c *gin.Context) {
	log.Trace()

	idStr := c.Param("room_id")

	if idStr == "" {
		log.Error("empty room id")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "room id can't be empty",
		})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid room id",
		})
		return
	}

	room, e := h.service.GetByID(id)
	if e != nil {
		log.Error(e)
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"room": room})
}

func (h *room) CheckIfAvailableById(c *gin.Context) {
	log.Trace()

	id, e := uuid.Parse(c.Param("room_id"))
	if e != nil {
		log.Errorf("Failed to parse room_id: %v", e)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid room_id",
		})
		return
	}

	var input struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	if e := c.ShouldBindJSON(&input); e != nil {
		log.Errorf("Failed to bind JSON input: %v", e)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	startDate, e := time.Parse("2006-01-02", input.StartDate)
	if e != nil {
		log.Error(e)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid start date",
		})
		return
	}

	endDate, e := time.Parse("2006-01-02", input.EndDate)
	if e != nil {
		log.Error(e)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid end date",
		})
		return
	}

	res, err := h.service.CheckIfAvailableById(id, startDate, endDate)
	if err != nil {
		log.Errorf("checking if room is available failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	if !res {
		c.JSON(http.StatusOK, gin.H{"response": "room is not available"})
	}

	c.JSON(http.StatusOK, gin.H{"response": "room is available"})
}

func (h *room) Update(c *gin.Context) {
	log.Trace()

	id, err := uuid.Parse(c.Param("room_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid room id",
		})
		return
	}

	var input struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Errorf("Failed to bind JSON input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	room := &model.Room{
		Id:      id,
		Name:    input.Name,
		Updated: time.Now(),
	}

	if err := h.service.Update(room); err != nil {
		log.Errorf("Failed to update room: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update room",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"room": room})
}
