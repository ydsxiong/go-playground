package lead

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type LeadServer struct {
	store LeadStore
	http.Handler
}

type httpHandlerMiddleware func(next http.HandlerFunc) http.HandlerFunc

var NoLeadFoundErr = errors.New("No lead data found!")

func NewLeadServer(store LeadStore) (*LeadServer, error) {
	server := new(LeadServer)

	server.store = store

	router := mux.NewRouter()
	// initially generate once a token for a regestered api user
	router.HandleFunc("/generatetoken", server.generateToken).Methods(http.MethodPost)
	// authorize all api calls that access lead data resource.
	router.HandleFunc("/lead/new", tokenAuthoringMiddleWare(server.createNew)).Methods(http.MethodPost)
	router.HandleFunc("/leads", tokenAuthoringMiddleWare(server.findAll))
	router.HandleFunc("/lead/{email}", tokenAuthoringMiddleWare(server.findByEmail))

	server.Handler = router

	return server, nil
}

/**
generate a token for registered user to use for their access to lead data resource via api call
*/
func (s *LeadServer) generateToken(w http.ResponseWriter, r *http.Request) {
	var claim claims
	err := json.NewDecoder(r.Body).Decode(&claim)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !claim.isValid() {
		respondError(w, http.StatusBadRequest, "Unrecognized user")
		return
	}
	token := generateApiUserToken(&claim)
	if token == nil {
		respondError(w, http.StatusInternalServerError, "Unable to generate a valid token, try again later")
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"x-access-token": *token})
}

func (s *LeadServer) createNew(w http.ResponseWriter, r *http.Request) {
	// if needed: claim := retriveValidClaimsFromContext(r)
	var lead Lead
	err := json.NewDecoder(r.Body).Decode(&lead)
	if err != nil {
		// the lead data will be valided inside lead's customerized unmarshaller
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	existingLead, err := s.store.FindByEmail(lead.Email)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	// email is unique among leads
	if existingLead != nil {
		respondError(w, http.StatusNotAcceptable, "The emails specified already exists in the system")
		return
	}
	err = s.store.Save(lead)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondJSON(w, http.StatusAccepted, lead)
	}
}

func (s *LeadServer) findAll(w http.ResponseWriter, r *http.Request) {
	// if needed: claim := retriveValidClaimsFromContext(r)
	leads, err := s.store.FindAll()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondJSON(w, http.StatusOK, leads)
	}
}

func (s *LeadServer) findByEmail(w http.ResponseWriter, r *http.Request) {
	// if needed: claim := retriveValidClaimsFromContext(r)
	vars := mux.Vars(r)
	email := vars["email"]

	result, err := s.store.FindByEmail(email)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err != nil || result == nil {
		respondError(w, http.StatusNotFound, NoLeadFoundErr.Error())
		return
	}
	respondJSON(w, http.StatusOK, result)
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, errCode int, message string) {
	w.WriteHeader(errCode)
	json.NewEncoder(w).Encode(errorWrapper{Error: message})
}

type errorWrapper struct {
	Error string `json:"error"`
}
