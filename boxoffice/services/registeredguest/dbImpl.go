package registeredguest

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/ydsxiong/go-playground/boxoffice/model"
)

type guestService struct {
	db *gorm.DB
}

func NewGuestService(gormdb *gorm.DB) RegisteredGuestService {
	return &guestService{gormdb}
}

func (svc *guestService) GetAllGuests() ([]*model.Guest, error) {

	guests := []*model.Guest{}
	if err := svc.db.Find(&guests).Error; err != nil {
		return nil, err
	}
	return guests, nil
}

func (svc *guestService) GetGuestByName(guestname string) (*model.Guest, error) {
	return findGuestByName(guestname, svc.db)
}

func (svc *guestService) SaveRegisteredGuest(name string) error {
	guest, err := findGuestByName(name, svc.db)
	if err != nil {
		return err
	}
	if guest == nil {
		guest = &model.Guest{Name: name, ExpiredAt: nil}
	} else {
		guest.ExpiredAt = nil
	}
	return svc.db.Save(guest).Error
}

func (svc *guestService) AddGuestInProgress(guest *model.Guest, reservationTime time.Duration) error {
	expiry := time.Now().Add(reservationTime)
	guest.ExpiredAt = &expiry
	return svc.db.Save(&guest).Error
}

func (s *guestService) RemoveGuestFromInProgress(guest *model.Guest) error {
	removeGuestFromDB(guest, s.db)
	return nil
}

func (s *guestService) IsGuestInProcess(guest *model.Guest) (bool, time.Duration, error) {
	if guest.ExpiredAt == nil {
		return false, 0, nil
	}
	remainingTime := guest.ExpiredAt.Sub(time.Now())

	if int(remainingTime.Seconds()) > 0 {
		return true, remainingTime, nil
	}
	removeGuestFromDB(guest, s.db)
	return false, 0, nil
}

// no longer needed, this api can be removed now.
func (s *guestService) NumberOfGuestInProcess() (num int, err error) {
	return 0, nil
}

func findGuestByName(guestname string, db *gorm.DB) (*model.Guest, error) {
	guest := model.Guest{}
	if err := db.First(&guest, model.Guest{Name: guestname}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &guest, nil
}

func removeGuestFromDB(guest *model.Guest, db *gorm.DB) {
	db.Unscoped().Delete(guest)
}
