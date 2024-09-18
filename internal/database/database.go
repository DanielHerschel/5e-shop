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

	CreateUser(ctx context.Context, item domain.User) (err error)
	GetUser(ctx context.Context, id primitive.ObjectID) (domain.User, error)
	UpdateUser(ctx context.Context, update domain.User) (err error)

	GetUserCampaigns(ctx context.Context, userId primitive.ObjectID) (res []domain.Campaign, err error)

	CreateCharacter(ctx context.Context, character domain.Character) (err error)
	GetCharacter(ctx context.Context, id primitive.ObjectID) (domain.Character, error)
	UpdateCharacter(ctx context.Context, update domain.Character) (err error)
	DeleteCharacter(ctx context.Context, characterId primitive.ObjectID) (err error)

	CreateItem(ctx context.Context, item domain.Item) (err error)
	GetItem(ctx context.Context, id primitive.ObjectID) (domain.Item, error)
	UpdateItem(ctx context.Context, update domain.Item) (err error)
	DeleteItem(ctx context.Context, itemId primitive.ObjectID) (err error)

	CreateShop(ctx context.Context, shop domain.Shop) (err error)
	GetShop(ctx context.Context, id primitive.ObjectID) (domain.Shop, error)
	UpdateShop(ctx context.Context, update domain.Shop) (err error)
	DeleteShop(ctx context.Context, shopId primitive.ObjectID) (err error)

	GetCampaignShops(ctx context.Context, campaignId primitive.ObjectID) (res []domain.Shop, err error)
	GetCampaignCurrentShop(ctx context.Context, campaignId primitive.ObjectID) (res domain.Shop, err error)

	CreateCampaign(ctx context.Context, campaign domain.Campaign) (err error)
	GetCampaign(ctx context.Context, id primitive.ObjectID) (res domain.Campaign, err error)
	UpdateCampaign(ctx context.Context, update domain.Campaign) (err error)
	DeleteCampaign(ctx context.Context, campaignId primitive.ObjectID) (err error)
}

type service struct {
	db *mongo.Client

	userCollection      *mongo.Collection
	characterCollection *mongo.Collection
	itemCollection      *mongo.Collection
	shopCollection      *mongo.Collection
	campaignCollection  *mongo.Collection
}

var (
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	database = os.Getenv("DB_DATABASE")

	userCollection      = os.Getenv("DB_USER_COLLECTION")
	characterCollection = os.Getenv("DB_CHARACTER_COLLECTION")
	itemCollection      = os.Getenv("DB_ITEM_COLLECTION")
	shopCollection      = os.Getenv("DB_SHOP_COLLECTION")
	campaignCollection  = os.Getenv("DB_CAMPAIGN_COLLECTION")
)

func New() Service {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port)))
	if err != nil {
		log.Fatal(err)
	}

	userCol := client.Database(database).Collection(userCollection)
	characterCol := client.Database(database).Collection(characterCollection)
	itemCol := client.Database(database).Collection(itemCollection)
	shopCol := client.Database(database).Collection(shopCollection)
	campaignCol := client.Database(database).Collection(campaignCollection)

	return &service{
		db:                  client,
		userCollection:      userCol,
		characterCollection: characterCol,
		itemCollection:      itemCol,
		shopCollection:      shopCol,
		campaignCollection:  campaignCol,
	}

}

// ----------------- Create funcs -----------------

func (s *service) CreateUser(ctx context.Context, user domain.User) (err error) {
	_, err = s.userCollection.InsertOne(ctx, user)
	return
}

func (s *service) CreateCharacter(ctx context.Context, character domain.Character) (err error) {
	_, err = s.characterCollection.InsertOne(ctx, character)
	return
}

func (s *service) CreateItem(ctx context.Context, item domain.Item) (err error) {
	_, err = s.itemCollection.InsertOne(ctx, item)
	return
}

func (s *service) CreateShop(ctx context.Context, shop domain.Shop) (err error) {
	_, err = s.shopCollection.InsertOne(ctx, shop)
	return
}

func (s *service) CreateCampaign(ctx context.Context, campagin domain.Campaign) (err error) {
	_, err = s.campaignCollection.InsertOne(ctx, campagin)
	return
}

// ----------------- Get funcs -----------------

func (s *service) GetUser(ctx context.Context, id primitive.ObjectID) (res domain.User, err error) {
	filter := bson.D{{Key: "_id", Value: id}}
	err = s.userCollection.FindOne(ctx, filter).Decode(&res)
	return
}

func (s *service) GetUserCampaigns(ctx context.Context, userId primitive.ObjectID) (res []domain.Campaign, err error) {
	filter := bson.D{{Key: "ownerId", Value: userId}}
	cur, err := s.campaignCollection.Find(ctx, filter)
	if err != nil {
		return
	}

	err = cur.All(ctx, &res)
	return
}

func (s *service) GetCharacter(ctx context.Context, id primitive.ObjectID) (res domain.Character, err error) {
	filter := bson.D{{Key: "_id", Value: id}}
	err = s.characterCollection.FindOne(ctx, filter).Decode(&res)
	return
}

func (s *service) GetUserCharacters(ctx context.Context, userId primitive.ObjectID) (res []domain.Character, err error) {
	filter := bson.D{{Key: "ownerId", Value: userId}}
	cur, err := s.characterCollection.Find(ctx, filter)
	if err != nil {
		return
	}

	err = cur.All(ctx, &res)
	return
}

func (s *service) GetItem(ctx context.Context, id primitive.ObjectID) (res domain.Item, err error) {
	filter := bson.D{{Key: "_id", Value: id}}
	err = s.itemCollection.FindOne(ctx, filter).Decode(&res)
	return
}

func (s *service) GetShopItems(ctx context.Context, shopId primitive.ObjectID) (res []domain.Item, err error) {
	filter := bson.D{{Key: "shopId", Value: shopId}}
	cur, err := s.itemCollection.Find(ctx, filter)
	if err != nil {
		return
	}

	err = cur.All(ctx, &res)
	return
}

func (s *service) GetShop(ctx context.Context, id primitive.ObjectID) (res domain.Shop, err error) {
	filter := bson.D{{Key: "_id", Value: id}}
	err = s.shopCollection.FindOne(ctx, filter).Decode(&res)
	return
}

func (s *service) GetCampaignShops(ctx context.Context, campaignId primitive.ObjectID) (res []domain.Shop, err error) {
	filter := bson.D{{Key: "campaignId", Value: campaignId}}
	cur, err := s.shopCollection.Find(ctx, filter)
	if err != nil {
		return
	}

	err = cur.All(ctx, &res)
	return
}

func (s *service) GetCampaign(ctx context.Context, id primitive.ObjectID) (res domain.Campaign, err error) {
	filter := bson.D{{Key: "_id", Value: id}}
	err = s.campaignCollection.FindOne(ctx, filter).Decode(&res)
	return
}

// ----------------- Update funcs -----------------

func (s *service) UpdateUser(ctx context.Context, update domain.User) (err error) {
	filter := bson.D{{Key: "_id", Value: update.Id}}
	_, err = s.userCollection.UpdateOne(ctx, filter, update)
	return
}

func (s *service) UpdateCharacter(ctx context.Context, update domain.Character) (err error) {
	filter := bson.D{{Key: "_id", Value: update.Id}}
	_, err = s.characterCollection.UpdateOne(ctx, filter, update)
	return
}

func (s *service) UpdateItem(ctx context.Context, update domain.Item) (err error) {
	filter := bson.D{{Key: "_id", Value: update.Id}}
	_, err = s.itemCollection.UpdateOne(ctx, filter, update)
	return
}

func (s *service) UpdateShop(ctx context.Context, update domain.Shop) (err error) {
	filter := bson.D{{Key: "_id", Value: update.Id}}
	_, err = s.shopCollection.UpdateOne(ctx, filter, update)
	return
}

func (s *service) UpdateCampaign(ctx context.Context, update domain.Campaign) (err error) {
	filter := bson.D{{Key: "_id", Value: update.Id}}
	_, err = s.campaignCollection.UpdateOne(ctx, filter, update)
	return
}

// ----------------- Delete funcs -----------------

func (s *service) DeleteCharacter(ctx context.Context, characterId primitive.ObjectID) (err error) {
	filter := bson.D{{Key: "_id", Value: characterId}}
	_, err = s.characterCollection.DeleteOne(ctx, filter)
	return
}

func (s *service) DeleteItem(ctx context.Context, itemId primitive.ObjectID) (err error) {
	filter := bson.D{{Key: "_id", Value: itemId}}
	_, err = s.itemCollection.DeleteOne(ctx, filter)
	return
}

func (s *service) DeleteShop(ctx context.Context, shopId primitive.ObjectID) (err error) {
	filter := bson.D{{Key: "_id", Value: shopId}}
	_, err = s.shopCollection.DeleteOne(ctx, filter)
	return
}

func (s *service) DeleteCampaign(ctx context.Context, campaignId primitive.ObjectID) (err error) {
	filter := bson.D{{Key: "_id", Value: campaignId}}
	_, err = s.shopCollection.DeleteOne(ctx, filter)
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
