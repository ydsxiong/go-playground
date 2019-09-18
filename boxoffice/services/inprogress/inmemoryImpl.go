package inprogress

import (
	"sync"
	"time"

	"github.com/ydsxiong/go-playground/boxoffice/model"
)

type basicService struct {
	inprogress map[string]*inprogressWrapper
	mux        sync.Mutex
}

type inprogressWrapper struct {
	timer    *time.Timer
	expireAt time.Time
}

func NewInMemoryService() InProgressGuestService {
	return &basicService{inprogress: make(map[string]*inprogressWrapper)}
}

func (bs *basicService) AddGuestInProgress(guest *model.Guest, reservationTime time.Duration) error {
	bs.mux.Lock()
	defer bs.mux.Unlock()

	bs.inprogress[guest.Name] = &inprogressWrapper{time.NewTimer(reservationTime), time.Now().Add(reservationTime)}
	return nil
}

func (bs *basicService) RemoveGuestFromInProgress(guest *model.Guest) error {
	bs.mux.Lock()
	defer bs.mux.Unlock()
	reservation, existing := bs.inprogress[guest.Name]
	if existing {
		running := isRunning(reservation.timer)
		if running {
			reservation.timer.Stop()
		}
		delete(bs.inprogress, guest.Name)
	}
	return nil
}

func (bs *basicService) IsGuestInProcess(guest *model.Guest) (bool, time.Duration, error) {
	bs.mux.Lock()
	defer bs.mux.Unlock()

	reservation, existing := bs.inprogress[guest.Name]
	if !existing {
		return false, 0, nil
	}
	running := isRunning(reservation.timer)

	if !running {
		delete(bs.inprogress, guest.Name)
	}
	return running, reservation.expireAt.Sub(time.Now()), nil
}

func isRunning(reservation *time.Timer) bool {
	select {
	case <-reservation.C:
		return false
	default:
		return true
	}
}

func (bs *basicService) NumberOfGuestInProcess() (num int, err error) {
	num = 0
	toremove := []string{}

	bs.mux.Lock()
	defer bs.mux.Unlock()

	for guest, reservation := range bs.inprogress {
		running := isRunning(reservation.timer)

		if !running {
			toremove = append(toremove, guest)
		} else {
			num++
		}
	}

	for _, guest := range toremove {
		delete(bs.inprogress, guest)
	}
	err = nil
	return
}
