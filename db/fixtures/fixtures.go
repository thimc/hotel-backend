package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/thimc/hotel-backend/db"
	"github.com/thimc/hotel-backend/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddBooking(store *db.Store, userID, roomID string, from, until time.Time) *types.Booking {
	booking := &types.Booking{
		UserID:    userID,
		RoomID:    roomID,
		FromDate:  from,
		UntilDate: until,
	}
	insertedBooking, err := store.Booking.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}

func AddRoom(store *db.Store, size string, price float64, hotelID string) *types.Room {
	room := &types.Room{
		Size:    size,
		Price:   price,
		HotelID: hotelID,
	}
	insertedRoom, err := store.Room.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func AddHotel(store *db.Store, name, location string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	roomIDs := rooms
	if rooms == nil {
		roomIDs = []primitive.ObjectID{}
	}
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    roomIDs,
		Rating:   rating,
	}
	insertedHotel, err := store.Hotel.Insert(context.Background(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func AddUser(store *db.Store, firstName, lastName string, isAdmin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: firstName,
		LastName:  lastName,
		Email:     fmt.Sprintf("%s@%s.com", firstName, lastName),
		Password:  fmt.Sprintf("%s_%s", firstName, lastName),
	})
	if err != nil {
		log.Fatal(err)
	}

	user.IsAdmin = isAdmin
	_, err = store.User.InsertUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
	return user
}
