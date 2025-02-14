package app

import (
	handler "github.com/demkowo/booking/handlers"
	log "github.com/sirupsen/logrus"
)

func reservationRoutes(h handler.Reservation) {
	log.Trace()

	reservations := router.Group("/api/v1/reservations")
	{
		reservations.POST("/add", h.Add)
		reservations.DELETE("/:reservation_id", h.Delete)
		reservations.GET("/", h.Find)
		reservations.GET("/find/:room_id", h.FindByRoomID)
		reservations.GET("/:reservation_id", h.GetById)
		reservations.PUT("/:reservation_id", h.Update)
	}
}
