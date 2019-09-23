package lead_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/ydsxiong/playground/customerlead/lead"
)

func TestUnmarshallLead(t *testing.T) {

	testcases := []struct {
		name   string
		input  []byte
		output error
	}{
		{
			name: "Full data test",
			input: []byte(`{
					"first_name": "x",
					"last_name": "y",
					"email": "x@y.com",
					"company": "xyz",
					"postcode": "xxx",
					"terms_accepted": true
					}`), output: nil,
		},
		{
			name: "Missing first name test",
			input: []byte(`{
					"last_name": "y",
					"email": "x@y.com",
					"company": "xyz",
					"postcode": "xxx",
					"terms_accepted": true
					}`), output: errors.New("Required fields missing: [First name]"),
		},
		{
			name: "Missing last name and email test",
			input: []byte(`{
					"first_name": "y",
					"company": "xyz",
					"postcode": "xxx",
					"terms_accepted": true
					}`), output: errors.New("Required fields missing: [last_name, email]"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var lead lead.Lead
			got := json.Unmarshal(tc.input, &lead)
			if !areSameError(got, tc.output) {
				t.Errorf("Expected output: %v, but got: %v", tc.output, got)
			}
		})
	}
}

func areSameError(err1, err2 error) bool {
	return (err1 == nil && err2 == nil) ||
		(err1 != nil && err2 != nil && err1.Error() != err2.Error())
}
