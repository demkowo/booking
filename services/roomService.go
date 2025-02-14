package service

import (
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	model "github.com/demkowo/booking/models"
	"github.com/demkowo/booking/utils/errs"
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

type Room interface {
	CreateTableRooms() string

	Add(*model.Room) *errs.Error
	Find() ([]*model.Room, *errs.Error)
	FindAvailable(time.Time, time.Time) ([]*model.Room, *errs.Error)
	GetByID(uuid.UUID) (*model.Room, *errs.Error)
	CheckIfAvailableById(uuid.UUID, time.Time, time.Time) (bool, *errs.Error)
	Update(*model.Room) *errs.Error
}

type room struct {
	repo RoomRepo
}

func NewRoom(repo RoomRepo) Room {
	log.Trace()

	return &room{
		repo: repo,
	}
}

func (s *room) CreateTableRooms() string {
	log.Trace()

	return s.repo.CreateTableRooms()
}

func (s *room) Add(room *model.Room) *errs.Error {
	log.Trace()

	if err := s.repo.Add(room); err != nil {
		return err
	}

	return nil
}

func (s *room) Find() ([]*model.Room, *errs.Error) {
	log.Trace()

	rooms, err := s.repo.Find()
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *room) FindAvailable(start time.Time, end time.Time) ([]*model.Room, *errs.Error) {
	log.Trace()

	rooms, err := s.repo.FindAvailable(start, end)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *room) GetByID(id uuid.UUID) (*model.Room, *errs.Error) {
	log.Trace()

	return s.repo.GetByID(id)
}

func (s *room) CheckIfAvailableById(id uuid.UUID, start time.Time, end time.Time) (bool, *errs.Error) {
	log.Trace()

	return s.repo.CheckIfAvailableById(id, start, end)
}

func (s *room) Update(room *model.Room) *errs.Error {
	log.Trace()

	return s.repo.Update(room)
}
