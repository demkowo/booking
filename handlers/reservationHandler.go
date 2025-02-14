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

type Reservation interface {
	CreateTableReservations()

	Add(*gin.Context)
	Delete(*gin.Context)
	Find(*gin.Context)
	FindByRoomID(*gin.Context)
	GetById(*gin.Context)
	Update(*gin.Context)
}

type reservation struct {
	service service.Reservation
}

func NewReservation(service service.Reservation) Reservation {
	log.Trace()

	return &reservation{
		service: service,
	}
}

func (h *reservation) CreateTableReservations() {
	log.Trace()

	h.service.CreateTableReservations()
}

func (h *reservation) Add(c *gin.Context) {
	log.Trace()

	var input struct {
		UserId    string `json:"user_id"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		RoomID    string `json:"room_id"`
		Status    int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Errorf("Failed to bind JSON input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	userId, err := uuid.Parse(input.UserId)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid user id",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid start date",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid end date",
		})
		return
	}

	roomId, err := uuid.Parse(input.RoomID)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid room id",
		})
		return
	}

	reservation := &model.Reservation{
		Id:        uuid.New(),
		UserId:    userId,
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomId,
		Status:    input.Status,
		Created:   time.Now(),
		Updated:   time.Now(),
		Deleted:   false,
	}

	if err := h.service.Add(reservation); err != nil {
		log.Errorf("Failed to create reservation: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create reservation",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reservation": reservation})
}

func (h *reservation) Delete(c *gin.Context) {
	log.Trace()

	if err := h.service.Delete(c.Param("reservation_id")); err != nil {
		log.Errorf("Failed to delete reservation: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete reservation",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reservation deleted successfully"})
}

func (h *reservation) Find(c *gin.Context) {
	log.Trace()

	reservations, err := h.service.Find()
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "list reservation failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reservations": reservations})
}

func (h *reservation) FindByRoomID(c *gin.Context) {
	log.Trace()

	roomId, e := uuid.Parse(c.Param("room_id"))
	if e != nil {
		log.Error(e)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid room id",
		})
		return
	}

	rooms, err := h.service.FindByRoomID(roomId)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "finding rooms by id failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

func (h *reservation) GetById(c *gin.Context) {
	log.Trace()

	idStr := c.Param("reservation_id")

	if idStr == "" {
		log.Error("empty reservation ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "reservation ID can't be empty",
		})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid reservation ID",
		})
		return
	}

	reservation, e := h.service.GetByID(id)
	if e != nil {
		log.Error(e)
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reservation": reservation})
}

func (h *reservation) Update(c *gin.Context) {
	log.Trace()

	idStr := c.Param("reservation_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Errorf("Invalid reservation id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation ID"})
		return
	}

	var input struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		RoomID    string `json:"room_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Errorf("Failed to bind JSON input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid start date",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid end date",
		})
		return
	}

	roomId, err := uuid.Parse(input.RoomID)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid room id",
		})
		return
	}

	reservation := &model.Reservation{
		Id:        id,
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomId,
	}

	if err := h.service.Update(reservation); err != nil {
		log.Errorf("Failed to update reservation: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update reservation",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reservation": reservation})
}
