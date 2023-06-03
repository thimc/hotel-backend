package db

type Pagination struct {
	Limit int64
	Page  int64
}

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

const (
	ENV_DB_NAME = "MONGODB_DB_NAME"
	ENV_DB_URI  = "MONGODB_DB_URI"

	ENV_TEST_DB_URI = "MONGODB_TEST_DB_URI"

	ENV_LISTEN_ADDRESS = "HTTP_LISTEN_ADDRESS"
	ENV_JWT_SECRET     = "JWT_SECRET"
)
