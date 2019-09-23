package lead_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ydsxiong/playground/customerlead/lead"
)

var validTokenHeader = map[string]string{"x-access-token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiYXBpdXNlciJ9.Yl17RXkS81JfiThymEo8lSE6xwABpquPcBxCBpYI-O8"}

func TestGeneratetoken(t *testing.T) {
	store := lead.NewInMemoryDataStore()
	server := mustMakeServer(t, store)

	validTokenString, _ := json.Marshal(validTokenHeader)
	testcases := []struct {
		name    string
		reqBody io.Reader
		resCode int
		resBody []byte
	}{
		{
			"generate token test 1",
			bytes.NewBuffer([]byte(``)),
			http.StatusBadRequest,
			[]byte(`{"error":"EOF"}`),
		},
		{
			"generate token test 2",
			bytes.NewBuffer([]byte(`{}`)),
			http.StatusBadRequest,
			[]byte(`{"error":"Unrecognized user"}`),
		},
		{
			"generate token test 2",
			bytes.NewBuffer([]byte(`{"user":"apiuser"}`)),
			http.StatusOK,
			[]byte(validTokenString),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			res, _ := createAndServeGenerateTokenReqRes(server, tc.reqBody)

			assertStatusCode(t, res, tc.resCode)
			assertResponseBody(t, res, tc.resBody)
		})
	}
}

func TestCreateNewLead(t *testing.T) {
	store := lead.NewInMemoryDataStore()
	lead1 := lead.Lead{Email: "one@abc.com", Fname: "abc"}
	store.Save(lead1)
	server := mustMakeServer(t, store)

	testcases := []struct {
		name      string
		path      string
		reqHeader map[string]string
		reqBody   io.Reader
		resCode   int
		resBody   []byte
		nRecords  int
	}{
		{
			"create new lead test 1",
			"/lead/new",
			map[string]string{},
			bytes.NewBuffer([]byte(``)),
			http.StatusForbidden,
			[]byte(`{"error":"A valid access token could not be found found in the request header!"}`),
			1,
		},
		{
			"create new lead test 2",
			"/lead/new",
			validTokenHeader,
			bytes.NewBuffer([]byte(`{}`)),
			http.StatusBadRequest,
			[]byte(`{"error":"Required fields missing: [first_name last_name email terms_accepted]"}`),
			1,
		},
		{
			"create new lead test 3",
			"/lead/new",
			validTokenHeader,
			bytes.NewBuffer([]byte(`{"email":"a@b.com", "first_name": "f", "last_name":"l", "terms_accepted": true}`)),
			http.StatusAccepted,
			[]byte(`{"email":"a@b.com"}`),
			2,
		},
		{
			"find all leads test 1",
			"/leads",
			map[string]string{},
			bytes.NewBuffer([]byte(``)),
			http.StatusForbidden,
			[]byte(`{"error":"A valid access token could not be found found in the request header!"}`),
			2,
		},
		{
			"find all leads test 2",
			"/leads",
			validTokenHeader,
			bytes.NewBuffer([]byte(``)),
			http.StatusOK,
			[]byte(`[{"email":"one@abc.com"}, {"email":"a@b.com"}]`),
			2,
		},
		{
			"find by email test 1",
			"/lead/one@abc.com",
			map[string]string{},
			bytes.NewBuffer([]byte(``)),
			http.StatusForbidden,
			[]byte(`{"error":"A valid access token could not be found found in the request header!"}`),
			2,
		},
		{
			"find by email test 2",
			"/lead/foo@bar.com",
			validTokenHeader,
			bytes.NewBuffer([]byte(``)),
			http.StatusNotFound,
			[]byte(`{"error":"No lead data found!"}`),
			2,
		},
		{
			"find by email test 3",
			"/lead/one@abc.com",
			validTokenHeader,
			bytes.NewBuffer([]byte(``)),
			http.StatusOK,
			[]byte(`{"email":"one@abc.com"}`),
			2,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			res, _ := createAndServeManageLeadReqRes(server, tc.path, tc.reqBody, tc.reqHeader)

			assertStatusCode(t, res, tc.resCode)
			assertResponseBody(t, res, tc.resBody)
			assertServerStore(t, store, tc.nRecords)
		})
	}
}

func createAndServeGenerateTokenReqRes(server *lead.LeadServer, body io.Reader) (*httptest.ResponseRecorder, *http.Request) {
	req, _ := http.NewRequest(http.MethodPost, "/generatetoken", body)

	res := httptest.NewRecorder()

	server.ServeHTTP(res, req)
	return res, req
}

func createAndServeManageLeadReqRes(server *lead.LeadServer, path string, body io.Reader, headers map[string]string) (*httptest.ResponseRecorder, *http.Request) {
	req, _ := http.NewRequest(http.MethodPost, path, body)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res := httptest.NewRecorder()

	server.ServeHTTP(res, req)
	return res, req
}

func assertStatusCode(t *testing.T, res *httptest.ResponseRecorder, expectedCode int) {
	t.Helper()
	if res.Code != expectedCode {
		t.Errorf("wrong status code, got %d, but wanted %d", res.Code, expectedCode)
	}
}

func assertResponseBody(t *testing.T, res *httptest.ResponseRecorder, expectedBody []byte) {
	t.Helper()
	got := make([]map[string]string, 10)
	wanted := make([]map[string]string, 10)
	_ = json.Unmarshal(res.Body.Bytes(), &got)
	_ = json.Unmarshal(expectedBody, &wanted)
	if got[0]["error"] != wanted[0]["error"] {
		t.Errorf("wrong response body, got %s, but wanted %s", res.Body.String(), string(expectedBody))
	}
	_, expectedEmail := wanted[0]["email"]
	_, gotEmail := got[0]["email"]
	if len(got) != len(wanted) ||
		(!expectedEmail && gotEmail) ||
		(expectedEmail && (!gotEmail || got[0]["email"] != wanted[0]["email"])) {
		t.Errorf("wrong response body, got %s, but wanted %s", res.Body.String(), string(expectedBody))
	}
}

func assertServerStore(t *testing.T, store lead.LeadStore, expectedRecords int) {
	t.Helper()
	total, _ := store.FindAll()
	if len(total) != expectedRecords {
		t.Errorf("wrong number of records in the store, got %d, but wanted %d", len(total), expectedRecords)
	}
}

func mustMakeServer(t *testing.T, store lead.LeadStore) *lead.LeadServer {
	t.Helper()
	server, err := lead.NewLeadServer(store)
	if err != nil {
		t.Fatal("problem creating lead server", err)
	}
	return server
}
