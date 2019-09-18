package test

import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/ydsxiong/go-playground/boxoffice/model"
)

var update = flag.Bool("update", false, "update .golden files")

func LoadCachedGoldenFile(cacheName, path string, content []byte) []byte {
	goldenfile := filepath.Join(path, cacheName+".golden")
	// optionally cache or update the golden file once initially, in order for it to be available in for subsequent testings
	if *update {
		ioutil.WriteFile(goldenfile, content, 0644)
	}
	cached, _ := ioutil.ReadFile(goldenfile)
	return cached
}

type MockDBService struct {
	DefaultGuests []*model.Guest
}

func (ms *MockDBService) GetAllGuests() ([]*model.Guest, error) {
	return ms.DefaultGuests, nil
}
func (ms *MockDBService) GetGuestByName(guestname string) (*model.Guest, error) {
	return nil, nil
}
func (ms *MockDBService) SaveRegisteredGuest(name string) error {
	return nil
}

func (ms *MockDBService) AddGuestInProgress(guest *model.Guest, reservationTime time.Duration) error {
	return nil
}

func (ms *MockDBService) RemoveGuestFromInProgress(guest *model.Guest) error {
	return nil
}

func (ms *MockDBService) IsGuestInProcess(guest *model.Guest) (bool, time.Duration, error) {
	return false, 0, nil
}

func (ms *MockDBService) NumberOfGuestInProcess() (num int, err error) {
	return 0, nil
}
