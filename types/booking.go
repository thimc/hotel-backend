package types

import (
	"time"
)

type Booking struct {
	ID         string    `bson:"_id,omitempty"        json:"id,omitempty"`
	UserID     string    `bson:"userID,omitempty"     json:"userID,omitempty"`
	RoomID     string    `bson:"roomID,omitempty"     json:"roomID,omitempty"`
	NumPersons int       `bson:"numPersons,omitempty" json:"numPersons,omitempty"`
	FromDate   time.Time `bson:"fromDate,omitempty"   json:"fromDate,omitempty"`
	UntilDate  time.Time `bson:"untilDate,omitempty"  json:"untilDate,omitempty"`
	Canceled   bool      `bson:"canceled"             json:"canceled"`
}
