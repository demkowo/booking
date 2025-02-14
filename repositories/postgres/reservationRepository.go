package postgres

import (
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	model "github.com/demkowo/booking/models"
	"github.com/demkowo/booking/utils/errs"

	_ "github.com/lib/pq"
)

const (
	CHECK_IF_RESERVATIONS_TABLE_EXIST = "SELECT to_regclass('public.reservations')"
	CREATE_RESERVATIONS_TABLE         = `CREATE TABLE public.reservations (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    start_date timestamptz NOT NULL,
    end_date timestamptz NOT NULL,
    room_id uuid NOT NULL,
	status INT NOT NULL DEFAULT 0,
    created timestamptz NOT NULL DEFAULT now(),
    updated timestamptz NOT NULL DEFAULT now(),
    deleted BOOLEAN NOT NULL DEFAULT FALSE, 
	CONSTRAINT reservations_pkey PRIMARY KEY (id)
	);`

	RESERVATION_CREATE          = "INSERT INTO reservations (id, user_id, start_date, end_date, room_id, status, created, updated, deleted) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
	RESERVATION_DELETE          = "UPDATE public.reservations SET deleted=TRUE, updated = $1 WHERE id = $2"
	RESERVATION_FIND            = "SELECT id, user_id, start_date, end_date, room_id, status, created, updated, deleted FROM reservations WHERE deleted = false ORDER BY updated DESC"
	RESERVATION_FIND_BY_ROOM_ID = "SELECT id, user_id, start_date, end_date, room_id, status, created, updated, deleted FROM reservations WHERE deleted = false AND room_id = $1 ORDER BY updated DESC"
	RESERVATION_GET_BY_ID       = "SELECT id, user_id, start_date, end_date, room_id, status, created, updated, deleted FROM reservations WHERE deleted = false AND id = $1"
	RESERVATION_UPDATE          = "UPDATE reservations SET start_date=$1, end_date=$2, room_id=$3, status=$4, updated=$5 WHERE id=$6"
)

type ReservationRepo interface {
	CreateTableReservations() string

	Add(*model.Reservation) *errs.Error
	Delete(string) *errs.Error
	Find() ([]*model.Reservation, *errs.Error)
	FindByRoomID(uuid.UUID) ([]*model.Reservation, *errs.Error)
	GetByID(uuid.UUID) (*model.Reservation, *errs.Error)
	Update(*model.Reservation) *errs.Error
}

type reservation struct {
	db *sql.DB
}

func NewReservation(db *sql.DB) ReservationRepo {
	return &reservation{
		db: db,
	}
}

func (r *reservation) CreateTableReservations() string {
	log.Trace()

	rows, err := r.db.Query(CHECK_IF_RESERVATIONS_TABLE_EXIST)
	if err != nil {
		log.Panicf("CHECK_IF_RESERVATIONS_TABLE_EXIST failed: %v", err)
	}
	defer rows.Close()

	var tableName sql.NullString
	for rows.Next() {
		err := rows.Scan(&tableName)
		if err != nil {
			log.Panicf("rows.Scan failed: %v", err)
		}
	}

	if tableName.Valid {
		return "Table reservations ready to go"
	}

	_, err = r.db.Exec(CREATE_RESERVATIONS_TABLE)
	if err != nil {
		log.Panicf("CREATE_RESERVATIONS_TABLE failed: %v", err)
	}

	return "Table reservations created, DB ready to go"
}

func (r *reservation) Add(reservation *model.Reservation) *errs.Error {
	log.Trace()

	_, err := r.db.Exec(RESERVATION_CREATE,
		&reservation.Id,
		&reservation.UserId,
		&reservation.StartDate,
		&reservation.EndDate,
		&reservation.RoomID,
		&reservation.Status,
		&reservation.Created,
		&reservation.Updated,
		&reservation.Deleted)
	if err != nil {
		log.Error("RESERVATION_CREATE failed", err)
		return errs.NewError("Failed to create reservation", 500, "Internal Server Error", []interface{}{})
	}

	return nil
}

func (r *reservation) Delete(id string) *errs.Error {
	log.Trace()

	updated := time.Now()
	_, err := r.db.Exec(RESERVATION_DELETE, updated, id)
	if err != nil {
		log.Error("RESERVATION_DELETE failed", err)
		return errs.NewError("Failed to delete reservation", 500, "Internal Server Error", []interface{}{err})
	}

	return nil
}

func (r *reservation) Find() ([]*model.Reservation, *errs.Error) {
	log.Trace()

	rows, err := r.db.Query(RESERVATION_FIND)
	if err != nil {
		log.Error("RESERVATION_FIND failed", err)
		return nil, errs.NewError("Failed to find reservation", 500, "Internal Server Error", []interface{}{})
	}
	defer rows.Close()

	reservations := []*model.Reservation{}
	for rows.Next() {
		reservation := &model.Reservation{}
		err := rows.Scan(&reservation.Id,
			&reservation.UserId,
			&reservation.StartDate,
			&reservation.EndDate,
			&reservation.RoomID,
			&reservation.Status,
			&reservation.Created,
			&reservation.Updated,
			&reservation.Deleted)
		if err != nil {
			log.Error("RESERVATION_FIND rows.Scan failed", err)
			return nil, errs.NewError("Failed to scan reservations", 500, "Internal Server Error", []interface{}{})
		}

		reservations = append(reservations, reservation)
	}

	if err := rows.Err(); err != nil {
		log.Error("RESERVATION_FIND rows.Err not nil", err)
		return nil, errs.NewError("Failed to find reservations", 500, "Internal Server Error", []interface{}{})
	}

	return reservations, nil
}

func (r *reservation) FindByRoomID(id uuid.UUID) ([]*model.Reservation, *errs.Error) {
	log.Trace()

	rows, err := r.db.Query(RESERVATION_FIND_BY_ROOM_ID, id)
	if err != nil {
		log.Error("RESERVATION_FIND_BY_ROOM_ID failed", err)
		return nil, errs.NewError("Failed to find reservations", 500, "Internal Server Error", []interface{}{})
	}
	defer rows.Close()

	reservations := []*model.Reservation{}
	for rows.Next() {
		reservation := &model.Reservation{}

		err := rows.Scan(&reservation.Id,
			&reservation.UserId,
			&reservation.StartDate,
			&reservation.EndDate,
			&reservation.RoomID,
			&reservation.Status,
			&reservation.Created,
			&reservation.Updated,
			&reservation.Deleted)
		if err != nil {
			log.Error("RESERVATION_FIND_BY_ROOM_ID rows.Scan failed", err)
			return nil, errs.NewError("Failed to scan reservations", 500, "Internal Server Error", []interface{}{})
		}

		reservations = append(reservations, reservation)
	}

	if err := rows.Err(); err != nil {
		log.Error("RESERVATION_FIND_BY_ROOM_ID rows.Err not nil", err)
		return nil, errs.NewError("Failed to find reservations", 500, "Internal Server Error", []interface{}{})
	}

	return reservations, nil
}

func (r *reservation) GetByID(id uuid.UUID) (*model.Reservation, *errs.Error) {
	log.Trace()

	row := r.db.QueryRow(RESERVATION_GET_BY_ID, id)
	reservation := &model.Reservation{}

	err := row.Scan(&reservation.Id,
		&reservation.UserId,
		&reservation.StartDate,
		&reservation.EndDate,
		&reservation.RoomID,
		&reservation.Status,
		&reservation.Created,
		&reservation.Updated,
		&reservation.Deleted)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			log.Tracef("RESERVATION_GET_BY_ID %s not found", id)
			return nil, errs.NewError("Reservation not found", 404, "Not Found", nil)
		}
		log.Errorf("RESERVATION_GET_BY_ID failed: %v", err)
		return nil, errs.NewError("Failed to get reservation", 500, "Internal Server Error", []interface{}{})
	}

	return reservation, nil
}

func (r *reservation) Update(reservation *model.Reservation) *errs.Error {
	log.Trace()

	updated := time.Now()
	_, err := r.db.Exec(RESERVATION_UPDATE, reservation.StartDate, reservation.EndDate, reservation.RoomID, reservation.Status, updated, reservation.Id)
	if err != nil {
		log.Error("RESERVATION_UPDATE failed", err)
		return errs.NewError("Failed to update reservation", 500, "Internal Server Error", []interface{}{})
	}

	return nil
}
