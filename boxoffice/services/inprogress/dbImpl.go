package inprogress

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/ydsxiong/go-playground/boxoffice/model"
)

const MAGIC_TIME_LAYOUT = "2006-01-02 15:04:05"

type dbService struct {
	db *gorm.DB
}

func NewDBService(gdb *gorm.DB) InProgressGuestService {
	return &dbService{db: gdb}
}

func (s *dbService) AddGuestInProgress(guest *model.Guest, reservationTime time.Duration) error {
	_ = time.Now().Add(reservationTime)
	// TODO
	return nil
}

func (s *dbService) RemoveGuestFromInProgress(guest *model.Guest) error {
	// TODO
	return nil
}

func (s *dbService) IsGuestInProcess(guest *model.Guest) (bool, time.Duration, error) {
	// TODO
	return false, 0, nil
}

func (s *dbService) NumberOfGuestInProcess() (num int, err error) {
	// TODO
	return 0, nil
}
