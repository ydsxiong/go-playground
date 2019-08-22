package people

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAndShowEmployee(t *testing.T) {
	e1 := CreateEmployee("tom hudson", "Building-1, Paris", 35)

	created := reflect.TypeOf(e1).String()
	assert.Equal(t, "*people.Emp", created)

	assert.Equal(t, "tom hudson", ShowName(e1))
}
