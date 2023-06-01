package types

type Hotel struct {
	ID       string           `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string           `bson:"name"          json:"name"`
	Location string           `bson:"location"      json:"location"`
	Rooms    []map[string]any `bson:"rooms"         json:"rooms"`
	Rating   int              `bson:"rating"        json:"rating"`
}

type Room struct {
	ID      string  `bson:"_id,omitempty" json:"id,omitempty"`
	Seaside bool    `bson:"seaside"       json:"seaside"`
	Size    string  `bson:"size"          json:"size"`
	Price   float64 `bson:"price"         json:"price"`
	HotelID string  `bson:"hotelID"       json:"hotelID"`
}
