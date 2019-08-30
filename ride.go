package wappa

import (
	"context"
	"net/http"
)

const rideEndpoint endpoint = `ride`

var rideFields = map[string]string{
	"id": "rideId",
}

// Ride is the  struct representing the ride
// entity in the API.
type Ride struct {
	Employee int `json:"employeeId"`
	TaxiType int `json:"taxiTypeId"`
	TaxiCategoryId int `json:"taxiCategoryID"`
	LatOrigin float64 `json:"latitudeOrigin"`
	LongOrigin float64 `json:"longitudeOrigin"`
	LatDestiny float64 `json:"latitudeDestiny"`
	LongDestiny float64 `json:"longitudeDestiny"`
	OriginRef string `json:"originReference"`
	ExternalID string `json:"externalID"`
}

type Passenger struct {
	ID int `json:"employeeId"`
	Name string `json:"name"`
	DDD string `json:"ddd"`
	Phone string `json:"phone"`
}

type Location struct {
	Lat float64 `json:"latitude"`
	Long float64 `json:"longitude"`
}

type Address struct {
	City string `json:"city"`
	State string `json:"state"`
	Country string `json:"country"`
	Address string `json:"address"`
	Location Location `json:"location"`
}

type Vehicle struct {
	Marker string `json:"marker"`
	Model string `json:"model"`
	Plate string `json:"plate"`
}

type Driver struct {
	Name string `json:"name"`
	DDD string `json:"ddd"`
	Phone string `json:"phone"`
	Vehicle Vehicle
}

type TravelInfo struct {
	Time string `json:"time"`
	TimeSec string `json:"timeSec"`
	DistanceKM float64 `json:"distanceKm"`
}

type RideInfo struct {
	// Possible values of Status:
	//  - searching-for-driver
	//  - driver-not-found
	//  - ride-cancelled
	//  - driver-found
	//  - waiting-for-driver
	//  - on-ride
	//  - ride-paid
	//  - ride-completed
	Status string `json:"status"`
	DriverLocation Location `json:"driverLocation"`
	ToOrigin TravelInfo `json:"toOrigin"`
	ToDestiny TravelInfo `json:"toDestiny"`
	// The agent that canceled the ride. 
	// Passenger = 1, Driver = 2, System = 3
	CancalledBy string `json:"cancelledBy"`
	// The reason that the ride was canceled for.
	CancalledReason string `json:"cancelledReason"`
	// The ride value, if available.
	rideValue float64 `json:"rideValue"`
	// The external ID provided when the ride was requested.
	externalID string `json:externalId"`
}

// DriverResult is the API response payload.
type RideResult struct {
	Result

	RideID int `json:"rideID"`
	Passenger Passenger `json:"passenger"`
	Origin Address `json:"origin"`
	Destiny Address `json:"destiny"`
	Driver Driver `json:"driver"`
}

// CancellationReasonResult represents the response of listing 
// the cancellation reason of rides.
type CancellationReasonResult struct {
	Reasons []Base `json:"reasons"`
}

// rideCancel contains the ride ID and reason ID.
// pulled off for testing.
type rideCancel struct {
	ID int `json:"rideId"`
	Reason int `json:"reasonId"`
}

// rideRate contains the ride ID and rating number.
// pulled off for testing.  
type rideRate struct {
	ID int `json:"rideId"`
	Rating int `json:"rating"`
}


// RideService is responsible for handling
// the requests to the ride resource.
type RideService struct {
	client requester
}

// Read returns the info of a ride.
func (rs *RideService) Read(ctx context.Context, f Filter) (*RideResult, error) {
	r := &RideResult{}

	if err := rs.client.Request(ctx, http.MethodGet, rideEndpoint.Query(f.Values(rideFields)), nil, r); err != nil {
		return nil, err
	}

	return r, nil
}

// Create creates a new ride in the API.
func (rs *RideService) Create(ctx context.Context, r *Ride) (*RideResult, error) {
	res := &RideResult{}

	if err := rs.client.Request(ctx, http.MethodPost, rideEndpoint, r, res); err != nil {
		return nil, err
	}

	return res, nil
}

// CancellationReason returns the list of possible reasons a user can choose when cancelling a ride.
func (rs *RideService) CancellationReason(ctx context.Context) (*CancellationReasonResult, error) {
	r := &CancellationReasonResult{}

	if err := rs.client.Request(ctx, http.MethodGet, indexEndpoint.Action(cancellationReason), nil, r); err != nil {
		return nil, err
	}

	return r, nil
}

//Cancel cancels a ride request.
func (rs *RideService) Cancel(ctx context.Context, ride int, reason int) (*Result, error) {
	r := &Result{}

	if err := rs.client.Request(ctx, http.MethodPost, rideEndpoint.Action(cancel), &rideCancel{ride, reason}, r); err != nil {
		return nil, err
	}

	return r, nil
}

// Rate created an experience rating of a ride.
func (rs *RideService) Rate(ctx context.Context, ride int, rating int) (*Result, error) {
	r := &Result{}

	if err := rs.client.Request(ctx, http.MethodPost, rideEndpoint.Action(rate), &rideRate{ride, rating}, r); err != nil {
		return nil, err
	}

	return r, nil
}
