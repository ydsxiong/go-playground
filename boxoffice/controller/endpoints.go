package controller

import (
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ydsxiong/go-playground/boxoffice/model"
	"github.com/ydsxiong/go-playground/boxoffice/services/inprogress"
	"github.com/ydsxiong/go-playground/boxoffice/services/registeredguest"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/sub"
)

var totalTicketsAvailable int
var reservationTime time.Duration
var guestService registeredguest.RegisteredGuestService
var inProgressService inprogress.InProgressGuestService
var page *template.Template

func Init(max_reservation_time time.Duration,
	total_tickets_available int,
	guestSvc registeredguest.RegisteredGuestService,
	inProgressSvc inprogress.InProgressGuestService,
	pageViews *template.Template) {

	totalTicketsAvailable = total_tickets_available
	reservationTime = max_reservation_time
	guestService = guestSvc
	inProgressService = inProgressSvc
	page = pageViews
}

func SwapGuestServiceWith(newService registeredguest.RegisteredGuestService) {
	guestService = newService
}

func DispatchHomePage(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	msg := checkAndRemoveCookieForSessionTimeout(w, r)
	if msg != "" {
		data["message"] = msg
	} else {
		msg = checkAndRemoveCookieForUnavailable(w, r)
		if msg != "" {
			data["message"] = msg
		}
	}

	available, reserved, registeredGuests, err := findNumberOfTicketsAvailable(true)
	if err != nil {
		handleInternalError(w, "")
		return
	}

	data["guests"] = registeredGuests
	data["remaining"] = available
	data["inreservation"] = reserved

	page.ExecuteTemplate(w, "Guests", data)
}

func findNumberOfTicketsAvailable(checkInprogress bool) (int, int, []*model.Guest, error) {

	allGuests, err := guestService.GetAllGuests()

	if err != nil {
		return 0, 0, nil, err
	}

	registeredGuests := guestFilter(allGuests, func(v *model.Guest) bool {
		return v.ExpiredAt == nil
	})

	available := totalTicketsAvailable - len(registeredGuests)

	reserved := len(allGuests) - len(registeredGuests)

	if checkInprogress {
		available -= reserved
	}
	return available, reserved, registeredGuests, nil
}

func DispatchReservationForm(w http.ResponseWriter, r *http.Request) {
	sendNewReservationForm(w, nil)
}

func sendNewReservationForm(w http.ResponseWriter, msg *string) {
	available, reserved, _, err := findNumberOfTicketsAvailable(true)
	if err != nil {
		handleInternalError(w, "")
		return
	}

	data := make(map[string]interface{})
	data["remaining"] = available
	data["reserved"] = reserved
	if msg != nil {
		data["msg"] = msg
	}

	page.ExecuteTemplate(w, "New", data)
}

func MakeReservation(w http.ResponseWriter, r *http.Request) {

	handleMsg := func(msg string, code int) {
		sendNewReservationForm(w, &msg)
	}

	name := r.FormValue("guestname")

	// first, sanitize and check the new input
	if name == "" {
		handleMsg("Error: please specify your name to reserve ticket for!", http.StatusBadRequest)
		return
	}

	existingGuest, err := guestService.GetGuestByName(name)
	if err != nil {
		handleInternalError(w, "")
		return
	}

	if existingGuest != nil && existingGuest.ExpiredAt == nil {
		handleMsg(name+" is an already registered guest, one guest can only reserve one ticket!", http.StatusNotAcceptable)
		return
	}

	// next, ensure that one guest can reserve only once while it's still in progress, before they can reserve another one
	tokenCookie := pullTokenFromClientCookie(r)
	if tokenCookie != nil {
		claims := validateToken(tokenCookie)
		if claims != nil {
			// check first the authenticated guest to see if anything in progress
			guest := existingGuest
			if guest == nil {
				guest = &model.Guest{Name: claims.GuestName}
			}
			running, remainingTime, _ := guestService.IsGuestInProcess(guest)
			if running {
				dispatchReservationStillInProgress(w, claims.GuestName, name, remainingTime)
				return
			} else if claims.GuestName != name { // check for this incoming other guest to see if anything in progress
				running, remainingTime, _ = guestService.IsGuestInProcess(&model.Guest{Name: name})
				if running {
					dispatchReservationStillInProgress(w, name, name, remainingTime)
					return
				}
				// otherwise let this incoming other guest to replace the existing one and start up a new reservation
			}
			// otherwise, let the same, existing guest start over again
		}
	} else if existingGuest != nil {
		running, remainingTime, _ := guestService.IsGuestInProcess(existingGuest)
		if running {
			dispatchReservationStillInProgress(w, name, name, remainingTime)
			return
		}
	}

	if !checkAvailability(w, r, true) {
		return
	}
	// finally: proceed with accepting this guest to the reversation in-progress list
	guestToReserve := existingGuest
	if guestToReserve == nil {
		guestToReserve = &model.Guest{Name: name}
	}
	err = guestService.AddGuestInProgress(guestToReserve, reservationTime)
	if err != nil {
		handleInternalError(w, "")
		return
	}

	tknstr, expiry := generateToken(&Claims{
		GuestName:      name,
		StandardClaims: jwt.StandardClaims{},
	})
	pushTokenIntoClientCookie(w, tknstr, expiry)

	dispatchNewReservationConfirmation(w, name)
}

func dispatchNewReservationConfirmation(w http.ResponseWriter, name string) {
	data := make(map[string]interface{})
	data["name"] = name
	data["greetings"] = "Nice to you meet you " + name + "!"
	data["remaining"] = fmt.Sprintf("We reserved your ticket for %d minutes", int(reservationTime.Minutes()))
	page.ExecuteTemplate(w, "Reservation", data)
}

func dispatchReservationStillInProgress(w http.ResponseWriter, name, other string, remainingTime time.Duration) {
	secondsLeft := int(remainingTime.Seconds())
	min, sec := secondsLeft/60, secondsLeft%60
	data := make(map[string]interface{})
	data["remaining"] = fmt.Sprintf("%s, you have %d minutes %d seconds left in reservation", name, min, sec)
	if name != other {
		data["warning"] = fmt.Sprintf("Can not start a new reservation for %s, while reservation for %s is still in progress", other, name)
	}
	page.ExecuteTemplate(w, "Reservation", data)
}

func checkAvailability(w http.ResponseWriter, req *http.Request, checkInprogress bool) bool {
	remaingTickets, reservedTickets, _, err := findNumberOfTicketsAvailable(checkInprogress)
	if err != nil {
		handleInternalError(w, "")
		return false
	}
	if (checkInprogress && remaingTickets == 0) || (!checkInprogress && reservedTickets == 0) {
		pushUnavailableCookie(w)
		redirectBackHome(w, req)
		return false
	}
	return true
}

func DoAuth(w http.ResponseWriter, req *http.Request) bool {
	tokenCookie := pullTokenFromClientCookie(req)
	if tokenCookie == nil {
		pushSessionTimeoutCookie(w)
		redirectBackHome(w, req)
		return false
	}
	claims := validateToken(tokenCookie)
	if claims == nil {
		redirectBackHome(w, req)
		return false
	}

	pushValidClaimIntoContext(req, claims)
	return true
}

func PayWithCard(w http.ResponseWriter, r *http.Request) {
	claims := retriveValidClaimsFromContext(r)
	if claims == nil {
		// if guest has been timed out, do not let they go ahead with payment, send them back to start over.
		redirectBackHome(w, r)
		return
	}

	if !checkAvailability(w, r, false) {
		return
	}

	// proceed with checking out the reservation for guest
	// Token is created using Checkout or Elements! Get the payment token ID submitted by the form:
	var payment paymentResult
	chanPayment := make(chan paymentResult)
	go processPayment(r.FormValue("stripeToken"), chanPayment)
	select {
	case payment = <-chanPayment:
	case <-time.After(5 * time.Second):
		handleInternalError(w, "Service was temporarily unavailable, please go back to try again later")
		return
	}
	if payment.err != nil {
		handleInternalError(w, payment.err.Error()+";  please go back to try again")
		return
	}

	// payment now done, so save the confirmed guest into the db
	if err := guestService.SaveRegisteredGuest(claims.GuestName); err != nil {
		// in the event of db saving failure, need to canx the charge
		chanCancel := make(chan error)
		go cancelPayment(payment.charge.ID, chanCancel)
		select {
		case e := <-chanCancel:
			if e != nil {
				// TODO: if canx failed, notify the customer to contact the system to resolve the pending charges.
				handleInternalError(w, e.Error()+";  please contact customr service to resolve the issue")
				return
			}
		case <-time.After(5 * time.Second):
			handleInternalError(w, "Service was temporarily unavailable, please contact customr service to resolve the issue")
			return
		}
	}

	// clear out and expire the token, and emove it off the in progress list
	clearoutTokenFromCookie(w)

	// finally dispatch the successful registeration confirmation to the guest
	data := make(map[string]interface{})
	data["name"] = claims.GuestName
	page.ExecuteTemplate(w, "Success", data)
}

func redirectBackHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func handleInternalError(w http.ResponseWriter, errMsg string) {
	var data map[string]interface{}
	if errMsg != "" {
		data = make(map[string]interface{})
		data["error"] = errMsg
	} else {
		data = nil
	}
	page.ExecuteTemplate(w, "ErrorPage", data)
}

type paymentResult struct {
	charge *stripe.Charge
	err    error
}

func processPayment(token string, result chan paymentResult) {
	stripe.Key = "sk_test_NJkFUrt4czgQdKvyHIMW3O9I007l9IMGx9"

	params := &stripe.ChargeParams{
		Amount:      stripe.Int64(999),
		Currency:    stripe.String(string(stripe.CurrencyGBP)),
		Description: stripe.String("Test charge"),
	}
	params.SetSource(token)
	charge, err := charge.New(params)
	result <- paymentResult{charge, err}
}

func cancelPayment(paymentId string, result chan error) {
	_, e := sub.Cancel(paymentId, nil)
	result <- e
}

func guestFilter(vs []*model.Guest, f func(*model.Guest) bool) []*model.Guest {
	vsf := make([]*model.Guest, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
