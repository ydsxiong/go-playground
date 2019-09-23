package lead_test

import (
	"sync"
	"testing"

	"github.com/ydsxiong/playground/customerlead/lead"
)

func TestSaveAndFetchSequencially(t *testing.T) {

	store := lead.NewInMemoryDataStore()

	lead1 := lead.Lead{Email: "one@abc.com", Fname: "abc"}
	lead2 := lead.Lead{Email: "two@abc.com", Fname: "def"}

	store.Save(lead1)
	store.Save(lead2)

	got, _ := store.FindAll()

	assertEqual(t, got[0], lead1)

	assertEqual(t, got[1], lead2)

	got2, _ := store.FindByEmail("two@abc.com")

	assertEqual(t, *got2, lead2)
}

func assertEqual(t *testing.T, l1, l2 lead.Lead) {
	t.Helper()
	for l1.Email != l2.Email || l1.Fname != l2.Fname {
		t.Errorf("expected leads %v, but got %v", l1, l2)
	}
}

func TestSaveAndFetchSimultaneously(t *testing.T) {
	store := lead.NewInMemoryDataStore()

	lead1 := lead.Lead{Email: "one@abc.com", Fname: "abc"}
	lead2 := lead.Lead{Email: "two@abc.com", Fname: "def"}

	store.Save(lead1)

	var wg sync.WaitGroup

	// fire up some random simultaneus writes and reads on the shared slice of data
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			store.Save(lead2)
		}()
	}
	for i := 0; i < 300; i++ {
		go func() {
			store.FindByEmail("xxx")
		}()
	}
	for i := 0; i < 300; i++ {
		go func() {
			store.FindByEmail("one@abc.com")
		}()
	}

	wg.Wait()

	// once all writes done, check the final total:
	got, err := store.FindAll()
	if err != nil {
		t.Errorf("unexpected err: %v", err)
	} else if len(got) != 101 {
		t.Errorf("expected leads number %d, but got %d", 101, len(got))
	}
}
