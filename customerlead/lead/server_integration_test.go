package lead_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/ydsxiong/playground/customerlead/lead"
)

/**
this is a scenario with all end-to-end api calls being given valid user inputs
the full scenario cases are covered in unit test suite
*/
func TestAPICallFlow(t *testing.T) {
	// setup an empty store, and a new http server
	store := lead.NewInMemoryDataStore()
	server := mustMakeServer(t, store)

	// endpoint 1 request a token for a authorized user
	requestBodyInput := `{"user":"apiuser"}`
	tokenRes, _ := createAndServeGenerateTokenReqRes(server, bytes.NewBuffer([]byte(requestBodyInput)))

	// endpoint 2. grab the token from previous response, and create a couple of new leads into the store
	requestTokenHeader := map[string]string{}
	_ = json.Unmarshal(tokenRes.Body.Bytes(), &requestTokenHeader)

	requestBodyInput = `{"email":"one@abc.com", "first_name": "f", "last_name":"l", "terms_accepted": true}`
	createAndServeManageLeadReqRes(server, "/lead/new", bytes.NewBuffer([]byte(requestBodyInput)), requestTokenHeader)

	requestBodyInput = `{"email":"a@b.com", "first_name": "ff", "last_name":"l", "company": "xxx", "terms_accepted": true}`
	createAndServeManageLeadReqRes(server, "/lead/new", bytes.NewBuffer([]byte(requestBodyInput)), requestTokenHeader)

	// endpoint 3. query all available leads in the store
	requestBodyInput = ``
	findAllRes, _ := createAndServeManageLeadReqRes(server, "/leads", bytes.NewBuffer([]byte(requestBodyInput)), requestTokenHeader)

	// assert results
	all := make([]map[string]interface{}, 10)
	_ = json.Unmarshal(findAllRes.Body.Bytes(), &all)

	if len(all) != 2 {
		t.Errorf("wrong number of records returned, got %d, but wanted %d", len(all), 2)
	}

	assertLeadData(t, all[0], map[string]interface{}{
		"email": "one@abc.com", "first_name": "f", "last_name": "l", "terms_accepted": true,
	})

	assertLeadData(t, all[1], map[string]interface{}{
		"email": "a@b.com", "first_name": "ff", "last_name": "l", "company": "xxx", "terms_accepted": true,
	})

	// endpoint 4. query specific lead by its email
	requestBodyInput = ``
	queryByEmailRes, _ := createAndServeManageLeadReqRes(server, "/lead/one@abc.com", bytes.NewBuffer([]byte(requestBodyInput)), requestTokenHeader)

	// assert results
	specificLead := make(map[string]interface{})
	json.Unmarshal(queryByEmailRes.Body.Bytes(), &specificLead)

	assertLeadData(t, specificLead, map[string]interface{}{
		"email": "one@abc.com", "first_name": "f", "last_name": "l", "terms_accepted": true,
	})
}

func assertLeadData(t *testing.T, leadData, expected map[string]interface{}) {
	t.Helper()
	for k, v := range expected {
		if leadData[k] != v {
			t.Errorf("got %s, but wanted %s", leadData[k], v)
		}
	}

}
