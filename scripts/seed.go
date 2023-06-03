package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/thimc/hotel-backend/api"
	"github.com/thimc/hotel-backend/db"
	"github.com/thimc/hotel-backend/db/fixtures"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	userStore := db.NewMongoUserStore(client)
	bookingStore := db.NewMongoBookingStore(client)
	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)
	db := &db.Store{
		User:    userStore,
		Booking: bookingStore,
		Hotel:   hotelStore,
		Room:    roomStore,
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

	for i := 0; i < 100; i++ {
		hotelName := fmt.Sprintf("Hotel_%d", i)
		hotelLocation := fmt.Sprintf("Location_%d", i)
		fixtures.AddHotel(db, hotelName, hotelLocation, rand.Intn(5)+1, nil)
	}
}
