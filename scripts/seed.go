package main

import (
	"context"
	"fmt"
	"hotel/api"
	"hotel/db"
	"hotel/db/fixtures"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var ctx = context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)
	userStore := db.NewMongoUserStore(client)
	bookingStore := db.NewMongoBookingStore(client)
	db := &db.Store{
		Hotel:   hotelStore,
		Room:    roomStore,
		User:    userStore,
		Booking: bookingStore,
	}

	user := fixtures.AddUser(db, "James", "Smith", false)
	fmt.Println("User:", api.CreateTokenFromUser(user))

	admin := fixtures.AddUser(db, "admin", "admin", true)
	fmt.Println("Admin:", api.CreateTokenFromUser(admin))

	hotel := fixtures.AddHotel(db, "Hotel", "Somewhere", 4, nil)
	fmt.Println("Hotel:", hotel)

	room := fixtures.AddRoom(db, "large", 199.99, hotel.ID)
	fmt.Println("Room:", room)

	booking := fixtures.AddBooking(db, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))
	fmt.Println("Booking:", booking)
}
