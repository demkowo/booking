package app

import (
	handler "github.com/demkowo/booking/handlers"
	log "github.com/sirupsen/logrus"
)

func roomRoutes(h handler.Room) {
	log.Trace()

	rooms := router.Group("/api/v1/rooms")
	{
		rooms.POST("/add", h.Add)
		rooms.GET("/", h.Find)
		rooms.POST("/find-available", h.FindAvailable)
		rooms.GET("/:room_id", h.GetById)
		rooms.POST("/:room_id/availability-check", h.CheckIfAvailableById)
		rooms.PUT("/:room_id", h.Update)
	}
}
