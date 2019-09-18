package controller_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"text/template"
	"time"

	"github.com/ydsxiong/go-playground/boxoffice/controller"
	"github.com/ydsxiong/go-playground/boxoffice/model"
	"github.com/ydsxiong/go-playground/boxoffice/services/inprogress"
	"github.com/ydsxiong/go-playground/boxoffice/test"
)

func init() {
	controller.Init(
		5*time.Minute, 5,
		&test.MockDBService{},
		inprogress.NewInMemoryService(),
		template.Must(template.ParseGlob("../views/*")))
}

func TestDispatchingPages(t *testing.T) {
	scenarios := []struct {
		name         string
		funcHandler  http.HandlerFunc
		httpPath     string
		httpMethod   string
		cookie       *http.Cookie
		defaultGuest []*model.Guest
	}{
		{"homepage", controller.DispatchHomePage, "/", "GET", nil, nil},
		{"homepageSessionTimeout", controller.DispatchHomePage, "/", "GET", &http.Cookie{Name: "sessionexpired", Value: "true"}, nil},
		{"reservationform", controller.DispatchReservationForm, "/", "reservation", nil, nil},
		{"reservationwith3ticketsleft", controller.DispatchReservationForm, "/", "reservation", nil, []*model.Guest{{Name: "test1"}, {Name: "test2"}}},
	}

	var testDataPath = "../test/data"

	// for each given scenario above, do http send and receive to check the results
	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// set up http test server
			req, err := http.NewRequest(scenario.httpMethod, scenario.httpPath, nil)
			if err != nil {
				t.Fatal(err)
			}
			if scenario.cookie != nil {
				req.AddCookie(scenario.cookie)
			}
			if scenario.defaultGuest != nil { // adjust test case data to alter controller behaviour outcome accordingly
				controller.SwapGuestServiceWith(&test.MockDBService{scenario.defaultGuest})
			}
			res := httptest.NewRecorder()
			handler := http.HandlerFunc(scenario.funcHandler)

			// serve http request and check results
			handler.ServeHTTP(res, req)
			if status := res.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}
			got := res.Body.Bytes()

			wanted := test.LoadCachedGoldenFile(scenario.name, testDataPath, got)

			if !bytes.Equal(got, wanted) {
				t.Errorf("%s handler returned unexpected body: got %v want %v", scenario.name, res.Body.String(), string(wanted))
			}
		})
	}
}
