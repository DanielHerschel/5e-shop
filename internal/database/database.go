package database

import (
	"5e-shop/internal/domain"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() map[string]string

	CreateShop(ctx context.Context, shop domain.Shop) (err error)
	GetShop(ctx context.Context, id primitive.ObjectID) (domain.Shop, error)
	UpdateShop(ctx context.Context, update domain.Shop) (err error)
	DeleteShop(ctx context.Context, update domain.Shop) (err error)
}

type service struct {
	db *mongo.Client
}

var (
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	database = os.Getenv("DB_DATABASE")

	shopCollection = os.Getenv("DB_SHOP_COLLECTION")
)

func New() Service {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port)))

	if err != nil {
		log.Fatal(err)

	}
	return &service{
		db: client,
	}
}

// ----------------- Create funcs -----------------

func (s *service) CreateShop(ctx context.Context, shop domain.Shop) (err error) {
	_, err = s.db.Database(database).Collection(shopCollection).InsertOne(ctx, shop)
	return
}

// ----------------- Get funcs -----------------

func (s *service) GetShop(ctx context.Context, id primitive.ObjectID) (res domain.Shop, err error) {
	filter := bson.D{{Key: "_id", Value: id}}
	err = s.db.Database(database).Collection(shopCollection).FindOne(ctx, filter).Decode(&res)
	return
}

func (s *service) GetShopsFromCampaign(ctx context.Context, campaignId primitive.ObjectID) (res []domain.Shop, err error) {
	filter := bson.D{{Key: "campaignId", Value: campaignId}}
	cur, err := s.db.Database(database).Collection(shopCollection).Find(ctx, filter)
	err = cur.All(ctx, &res)
	return
}

// ----------------- Update funcs -----------------

func (s *service) UpdateShop(ctx context.Context, update domain.Shop) (err error) {
	filter := bson.D{{Key: "_id", Value: update.Id}}
	_, err = s.db.Database(database).Collection(shopCollection).UpdateOne(ctx, filter, update)
	return
}

// ----------------- Delete funcs -----------------
func (s *service) DeleteShop(ctx context.Context, update domain.Shop) (err error) {
	filter := bson.D{{Key: "_id", Value: update.Id}}
	_, err = s.db.Database(database).Collection(shopCollection).DeleteOne(ctx, filter)
	return
}

// ----------------- Misc funcs -----------------

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.Ping(ctx, nil)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}
