package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/thimc/hotel-backend/api"
	"github.com/thimc/hotel-backend/api/errors"
	"github.com/thimc/hotel-backend/db"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		if apiError, ok := err.(errors.Error); ok {
			return c.Status(apiError.Code).JSON(apiError.Message)
		}
		apiErr := errors.NewError(http.StatusInternalServerError, err.Error())
		return c.Status(apiErr.Code).JSON(apiErr)
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "listen port of the api server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	var (
		userStore    = db.NewMongoUserStore(client)
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		bookingStore = db.NewMongoBookingStore(client)

		store = &db.Store{
			User:    userStore,
			Hotel:   hotelStore,
			Room:    roomStore,
			Booking: bookingStore,
		}

		authHandler    = api.NewAuthHandler(userStore)
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)

		app   = fiber.New(config)
		auth  = app.Group("/")
		v1    = app.Group("/v1", api.JWTAuthentication(userStore))
		admin = v1.Group("/admin", api.AdminAuth)
	)

	// auth handler
	auth.Post("/auth", authHandler.HandleAuthenticate)
	fmt.Println()

	// user handlers
	v1.Post("/user", userHandler.HandlePostUser)
	v1.Get("/user", userHandler.HandleGetUsers)
	v1.Get("/user/:id", userHandler.HandleGetUser)
	v1.Put("/user/:id", userHandler.HandlePutUser)
	v1.Delete("/user/:id", userHandler.HandleDeleteUser)

	// hotel handlers
	v1.Get("/hotel", hotelHandler.HandleGetHotels)
	v1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	v1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// room handlers
	v1.Get("/room", roomHandler.HandleGetRooms)
	v1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	// booking handlers
	v1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	v1.Post("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	// admin handlers
	admin.Get("/booking", bookingHandler.HandleGetBookings)

	if err := app.Listen(*listenAddr); err != nil {
		log.Fatal(err)
	}
}
