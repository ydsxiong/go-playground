package registeredguest

import (
	"github.com/ydsxiong/go-playground/boxoffice/model"
	"github.com/ydsxiong/go-playground/boxoffice/services/inprogress"
)

type RegisteredGuestService interface {
	GetAllGuests() ([]*model.Guest, error)
	GetGuestByName(guestname string) (*model.Guest, error)
	SaveRegisteredGuest(name string) error
	inprogress.InProgressGuestService
}
