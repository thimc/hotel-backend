package db

const (
	DBURI  = "mongodb://10.0.2.3:27017"
	DBNAME = "hotel-reservation"

	TestDBNAME = "hotel-reservation-test"
)

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}