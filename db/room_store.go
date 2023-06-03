package db

import (
	"context"
	"os"

	"github.com/thimc/hotel-backend/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	GetRoomByID(context.Context, string) (*types.Room, error)
	GetRooms(context.Context, map[string]any) ([]*types.Room, error)
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	dbName := os.Getenv(ENV_DB_NAME)
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(dbName).Collection("rooms"),
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) GetRoomByID(ctx context.Context, id string) (*types.Room, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var room types.Room
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&room); err != nil {
		return nil, err
	}
	return &room, nil
}

func (s *MongoRoomStore) GetRooms(ctx context.Context, filter map[string]any) ([]*types.Room, error) {
	resp, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var rooms []*types.Room
	if err := resp.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	resp, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = resp.InsertedID.(primitive.ObjectID).Hex()

	hid, _ := primitive.ObjectIDFromHex(room.HotelID)
	oid, _ := primitive.ObjectIDFromHex(room.ID)

	// update the hotel with this id
	filter := bson.M{"_id": hid}
	update := bson.M{"$push": bson.M{"rooms": oid}}
	if err := s.HotelStore.Update(ctx, filter, update); err != nil {
		return nil, err
	}

	return room, nil
}
