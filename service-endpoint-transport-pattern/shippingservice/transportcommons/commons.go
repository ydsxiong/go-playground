package transportcommons

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/shippingservice/cargo"
	"github.com/ydsxiong/go-playground/service-endpoint-transport-pattern/shippingservice/servicecommons"
)

func HandleResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(businessErrorer); ok && e.error() != nil {
		ProcessResponseError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type businessErrorer interface {
	error() error
}

// encode errors from business-logic
func ProcessResponseError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case cargo.ErrUnknown:
		w.WriteHeader(http.StatusNotFound)
	case servicecommons.ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
