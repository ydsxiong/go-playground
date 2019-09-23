package lead_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ydsxiong/playground/customerlead/lead"
)

func TestFileSystemStore(t *testing.T) {
	initData := `[
		{"first_name": "f", "email": "one@abc.com", "last_name": "l", "terms_accepted": true},
		{"first_name": "ff", "email": "a@b.com","last_name": "l", "company": "xxx", "terms_accepted": true}]`

	fd, err, cleanStore := createTempFile(initData)
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}
	defer cleanStore()

	store, err := lead.NewFileSystemStore(fd)
	if err != nil {
		t.Fatalf("Problem with opening file store, %v", err)
	}

	t.Run("/leads from a reader", func(t *testing.T) {
		got, _ := store.FindAll()
		byteData, _ := json.Marshal(got)
		all := make([]map[string]interface{}, 10)
		json.Unmarshal(byteData, &all)

		assertLeadData(t, all[0], map[string]interface{}{
			"email": "one@abc.com", "first_name": "f", "last_name": "l", "terms_accepted": true,
		})

		assertLeadData(t, all[1], map[string]interface{}{
			"email": "a@b.com", "first_name": "ff", "last_name": "l", "company": "xxx", "terms_accepted": true,
		})
	})

	t.Run("get specific lead by its email", func(t *testing.T) {

		got, _ := store.FindByEmail("one@abc.com")
		byteData, _ := json.Marshal(got)
		specificLead := make(map[string]interface{})
		json.Unmarshal(byteData, &specificLead)

		assertLeadData(t, specificLead, map[string]interface{}{
			"email": "one@abc.com", "first_name": "f", "last_name": "l", "terms_accepted": true,
		})
	})

	t.Run("can work with an empty file", func(t *testing.T) {
		fd, _, cleanUpStore := createTempFile("")
		defer cleanUpStore()

		_, err := lead.NewFileSystemStore(fd)

		if err != nil {
			t.Fatalf("Problem with opening file store, %v", err)
		}
	})
}

func TestTape_Write(t *testing.T) {
	fd, err, clean := createTempFile("123456")
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}
	defer clean()

	tape := lead.NewTapeFile(fd)

	tape.Write([]byte("abc"))

	fd.Seek(0, 0)
	newFileContents, _ := ioutil.ReadAll(fd)

	got := string(newFileContents)
	want := "abc"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func createTempFile(initialData string) (*os.File, error, func()) {

	fd, err := ioutil.TempFile(".", "filestore")

	if err != nil {
		return nil, err, nil
	}

	fd.Write([]byte(initialData))

	removeFile := func() {
		fd.Close()
		os.Remove(fd.Name())
	}

	return fd, nil, removeFile
}
