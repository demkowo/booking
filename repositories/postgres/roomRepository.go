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
	CHECK_IF_ROOMS_TABLE_EXIST = "SELECT to_regclass('public.rooms')"
	CREATE_ROOMS_TABLE         = `CREATE TABLE IF NOT EXISTS public.rooms (
		id uuid NOT NULL,
		name varchar(255),
		created timestamptz NOT NULL,
		updated timestamptz NOT NULL
	);`

	ROOM_CREATE         = "INSERT INTO rooms (id, name, created, updated) VALUES ($1, $2, $3, $4)"
	ROOMS_FIND          = "SELECT id, name, created, updated FROM rooms ORDER BY name ASC"
	ROOMS_FIND_AVAILABE = `
	select
			r.id, r.name, r.created, r.updated
		from
			rooms r
		where r.id not in 
		(select room_id from reservations rr where $1 < rr.end_date and $2 > rr.start_date);
	`
	ROOM_GET_BY_ID                = "SELECT id, name, created, updated FROM rooms WHERE id = $1"
	ROOM_CHECK_IF_AVAILABLE_BY_ID = `
	select
			r.id, r.name, r.created, r.updated
		from
			rooms r
		where r.id=$1 
		and r.id not in (select room_id from reservations rr where $2 < rr.end_date and $3 > rr.start_date);
	`
	ROOM_UPDATE = "UPDATE rooms SET name=$1, updated=$2 WHERE id=$3"
)

type RoomRepo interface {
	CreateTableRooms() string

	Add(*model.Room) *errs.Error
	Find() ([]*model.Room, *errs.Error)
	FindAvailable(time.Time, time.Time) ([]*model.Room, *errs.Error)
	GetByID(uuid.UUID) (*model.Room, *errs.Error)
	CheckIfAvailableById(uuid.UUID, time.Time, time.Time) (bool, *errs.Error)
	Update(*model.Room) *errs.Error
}

type room struct {
	db *sql.DB
}

func NewRoom(db *sql.DB) RoomRepo {
	return &room{
		db: db,
	}
}

func (r *room) CreateTableRooms() string {
	log.Trace()

	rows, err := r.db.Query(CHECK_IF_ROOMS_TABLE_EXIST)
	if err != nil {
		log.Panicf("CHECK_IF_ROOMS_TABLE_EXIST failed: %v", err)
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
		return "Table rooms ready to go"
	}

	_, err = r.db.Exec(CREATE_ROOMS_TABLE)
	if err != nil {
		log.Panicf("CREATE_ROOMS_TABLE failed: %v", err)
	}

	return "Table rooms created, DB ready to go"
}

func (r *room) Add(room *model.Room) *errs.Error {
	log.Trace()

	created := time.Now()
	updated := created
	_, err := r.db.Exec(ROOM_CREATE, &room.Id, &room.Name, created, updated)
	if err != nil {
		log.Error("ROOM_CREATE failed", err)
		return errs.NewError("Failed to create room", 500, "Internal Server Error", []interface{}{})
	}

	return nil
}

func (r *room) Find() ([]*model.Room, *errs.Error) {
	log.Trace()

	rows, err := r.db.Query(ROOMS_FIND)
	if err != nil {
		log.Error("ROOMS_FIND failed", err)
		return nil, errs.NewError("Failed to find rooms", 500, "Internal Server Error", []interface{}{})
	}
	defer rows.Close()

	rooms := []*model.Room{}
	for rows.Next() {
		room := &model.Room{}

		err := rows.Scan(&room.Id, &room.Name, &room.Created, &room.Updated)
		if err != nil {
			log.Error("ROOMS_FIND rows.Scan failed", err)
			return nil, errs.NewError("Failed to scan rooms", 500, "Internal Server Error", []interface{}{})
		}

		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		log.Error("ROOMS_FIND rows.Err not nil", err)
		return nil, errs.NewError("Failed to find rooms", 500, "Internal Server Error", []interface{}{})
	}

	return rooms, nil
}

func (r *room) FindAvailable(start time.Time, end time.Time) ([]*model.Room, *errs.Error) {
	log.Trace()

	rows, err := r.db.Query(ROOMS_FIND_AVAILABE, start, end)
	if err != nil {
		log.Error("ROOMS_FIND_AVAILABE failed", err)
		return nil, errs.NewError("Failed to find available rooms", 500, "Internal Server Error", []interface{}{})
	}
	defer rows.Close()

	rooms := []*model.Room{}
	for rows.Next() {
		room := &model.Room{}

		err := rows.Scan(&room.Id, &room.Name, &room.Created, &room.Updated)
		if err != nil {
			log.Error("ROOMS_FIND_AVAILABE rows.Scan failed", err)
			return nil, errs.NewError("Failed to scan available rooms", 500, "Internal Server Error", []interface{}{})
		}

		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		log.Error("ROOMS_FIND_AVAILABE rows.Err not nil", err)
		return nil, errs.NewError("Failed to find available rooms", 500, "Internal Server Error", []interface{}{})
	}

	return rooms, nil
}

func (r *room) GetByID(id uuid.UUID) (*model.Room, *errs.Error) {
	log.Trace()

	row := r.db.QueryRow(ROOM_GET_BY_ID, id)
	room := &model.Room{}

	err := row.Scan(&room.Id, &room.Name, &room.Created, &room.Updated)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			log.Tracef("ROOM_GET_BY_ID room %s not found", id)
			return nil, errs.NewError("room not found", 404, "Not Found", nil)
		}
		log.Errorf("ROOM_GET_BY_ID failed: %v", err)
		return nil, errs.NewError("Failed to get room", 500, "Internal Server Error", []interface{}{})
	}

	return room, nil
}

func (r *room) CheckIfAvailableById(id uuid.UUID, start time.Time, end time.Time) (bool, *errs.Error) {
	log.Trace()

	row := r.db.QueryRow(ROOM_CHECK_IF_AVAILABLE_BY_ID, id, start, end)
	room := &model.Room{}

	err := row.Scan(&room.Id, &room.Name, &room.Created, &room.Updated)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			log.Tracef("ROOM_CHECK_IF_AVAILABLE_BY_ID %s not found", id)
			return false, nil
		}
		log.Errorf("ROOM_CHECK_IF_AVAILABLE_BY_ID failed: %v", err)
		return false, errs.NewError("Failed to check if room is available", 500, "Internal Server Error", []interface{}{})
	}

	return true, nil
}

func (r *room) Update(room *model.Room) *errs.Error {
	log.Trace()

	updated := time.Now()
	_, err := r.db.Exec(ROOM_UPDATE, room.Name, updated, room.Id)
	if err != nil {
		log.Error("ROOM_UPDATE failed", err)
		return errs.NewError("Failed to update room", 500, "Internal Server Error", []interface{}{})
	}

	return nil
}
