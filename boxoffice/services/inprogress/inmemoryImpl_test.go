package inprogress_test

import (
	"testing"
	"time"

	"github.com/ydsxiong/go-playground/boxoffice/model"
	"github.com/ydsxiong/go-playground/boxoffice/services/inprogress"
)

func TestGuestInProgress(t *testing.T) {

	inProgressService := inprogress.NewInMemoryService()

	inProgressService.AddGuestInProgress(&model.Guest{Name: "mark"}, 10*time.Second)
	inProgressService.AddGuestInProgress(&model.Guest{Name: "baker"}, 20*time.Second)

	expected := 2
	if got, _ := inProgressService.NumberOfGuestInProcess(); got != expected {
		t.Errorf("expected %d. Got %d instead", expected, got)
	}

	yes := true
	if got, _, _ := inProgressService.IsGuestInProcess(&model.Guest{Name: "mark"}); got != yes {
		t.Errorf("expected %t. Got %t instead", yes, got)
	}

	time.Sleep(12 * time.Second)

	yes = false
	if got, _, _ := inProgressService.IsGuestInProcess(&model.Guest{Name: "mark"}); got != yes {
		t.Errorf("expected %t. Got %t instead", yes, got)
	}

	expected = 1
	if got, _ := inProgressService.NumberOfGuestInProcess(); got != expected {
		t.Errorf("expected %d. Got %d instead", expected, got)
	}

	yes = true
	if got, _, _ := inProgressService.IsGuestInProcess(&model.Guest{Name: "baker"}); got != yes {
		t.Errorf("expected %t. Got %t instead", yes, got)
	}

	expected = 0
	inProgressService.RemoveGuestFromInProgress(&model.Guest{Name: "baker"})
	if got, _ := inProgressService.NumberOfGuestInProcess(); got != expected {
		t.Errorf("expected %d. Got %d instead", expected, got)
	}
}
