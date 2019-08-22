package servicecommons

import (
	"errors"
)

// ErrInvalidArgument is returned when one or more arguments are invalid.
var ErrInvalidArgument = errors.New("invalid argument")

type ServicePath int

var (
	BookingBasePath  = "/booking/v1/"
	HandlingBasePath = "/handling/v1/"
	TrackingBasePath = "/tracking/v1/"

	cargosPath   = "cargos"
	cargosIdPath = "cargos/{id}"
)

const (
	CargosPath ServicePath = iota
	LoadCargoPath
	RequestRoutesPath
	AssignToRoutePath
	ChangeDestinationPath
	ListLocationsPath
	RegisterIncidentPath
	TrackCargoPath
)

func (sp ServicePath) String() string {
	return [...]string{
		BookingBasePath + cargosPath,
		BookingBasePath + cargosIdPath,
		BookingBasePath + cargosIdPath + "/request_routes",
		BookingBasePath + cargosIdPath + "/assign_to_route",
		BookingBasePath + cargosIdPath + "/change_destination",
		BookingBasePath + "locations",
		HandlingBasePath + "incidents",
		TrackingBasePath + cargosIdPath,
	}[sp]
}
