package utils_test

import (
	"strings"
	"testing"

	"github.com/ydsxiong/go-playground/utils"
)

func TestGetDataIntegration(t *testing.T) {
	got := utils.GetData(utils.GetXMLFromCommand())
	want := "HAPPY NEW YEAR!"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestGetdata(t *testing.T) {
	input := strings.NewReader(`
<payload>
    <message>Happy new year!</message>
</payload>`)

	got := utils.GetData(input)
	want := "HAPPY NEW YEAR!"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
