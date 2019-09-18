package inprogress

import (
	"time"

	"github.com/ydsxiong/go-playground/boxoffice/model"
)

type InProgressGuestService interface {
	AddGuestInProgress(guest *model.Guest, reservationTime time.Duration) error
	RemoveGuestFromInProgress(guest *model.Guest) error
	IsGuestInProcess(guest *model.Guest) (bool, time.Duration, error)
	NumberOfGuestInProcess() (num int, err error)
}
