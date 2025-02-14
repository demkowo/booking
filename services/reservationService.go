package service

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	model "github.com/demkowo/booking/models"
	"github.com/demkowo/booking/utils/errs"
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

type Reservation interface {
	CreateTableReservations() string

	Add(*model.Reservation) *errs.Error
	Delete(string) *errs.Error
	Find() ([]*model.Reservation, *errs.Error)
	FindByRoomID(uuid.UUID) ([]*model.Reservation, *errs.Error)
	GetByID(uuid.UUID) (*model.Reservation, *errs.Error)
	Update(*model.Reservation) *errs.Error
}

type reservation struct {
	repo ReservationRepo
}

func NewReservation(repo ReservationRepo) Reservation {
	log.Trace()

	return &reservation{
		repo: repo,
	}
}

func (s *reservation) CreateTableReservations() string {
	log.Trace()

	return s.repo.CreateTableReservations()
}

func (s *reservation) Add(reservation *model.Reservation) *errs.Error {
	log.Trace()

	if err := s.repo.Add(reservation); err != nil {
		return err
	}

	return nil
}

func (s *reservation) Delete(id string) *errs.Error {
	log.Trace()

	if id == "" {
		return errs.NewError("ID is required", 400, "Bad Request", nil)
	}

	return s.repo.Delete(id)
}

func (s *reservation) Find() ([]*model.Reservation, *errs.Error) {
	log.Trace()

	res, err := s.repo.Find()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *reservation) FindByRoomID(id uuid.UUID) ([]*model.Reservation, *errs.Error) {
	log.Trace()

	res, err := s.repo.FindByRoomID(id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *reservation) GetByID(id uuid.UUID) (*model.Reservation, *errs.Error) {
	log.Trace()

	return s.repo.GetByID(id)
}

func (s *reservation) Update(reservation *model.Reservation) *errs.Error {
	log.Trace()

	return s.repo.Update(reservation)
}
