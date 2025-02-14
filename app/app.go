package app

import (
	"database/sql"
	"log"
	"os"

	handler "github.com/demkowo/booking/handlers"
	"github.com/demkowo/booking/repositories/postgres"
	service "github.com/demkowo/booking/services"
	"github.com/demkowo/booking/utils/logger"
	"github.com/gin-gonic/gin"
)

const (
	portNumber = ":5000"
)

var (
	router       = gin.Default()
	dbConnection string
)

func init() {
	logger.Start.BasicConfig()
	dbConnection = os.Getenv("DB_DELMAJK")
}

func Start() {
	db, err := sql.Open("postgres", dbConnection)
	if err != nil {
		log.Panicf("sql.Open failed\n[%s]\n", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Panicf("db.Ping failed\n[%s]\n", err)
	}

	roomRepo := postgres.NewRoom(db)
	roomService := service.NewRoom(roomRepo)
	roomHandler := handler.NewRoom(roomService)
	roomRoutes(roomHandler)

	reservationRepo := postgres.NewReservation(db)
	reservationService := service.NewReservation(reservationRepo)
	reservationHandler := handler.NewReservation(reservationService)
	reservationRoutes(reservationHandler)

	roomHandler.CreateTableRooms()
	reservationHandler.CreateTableReservations()

	router.Run(portNumber)
}
