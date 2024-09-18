package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID `json:"id" bson:"_id"` // UUID
	Username string             `json:"username" bson:"username"`
}

type Campaign struct {
	Id   primitive.ObjectID `json:"id" bson:"_id"` // UUID
	Name string             `json:"name" bson:"name"`

	OwnerId    primitive.ObjectID `json:"ownerId" bson:"ownerId"`       // UUID user
	ActiveShop primitive.ObjectID `json:"activeShop" bson:"activeShop"` // UUID shop
}

type Character struct {
	Id   string `json:"id" bson:"_id"` // UUID
	Name string `json:"name" bson:"name"`

	IsInUse    bool
	Balance    Balance
	CampaignId primitive.ObjectID `json:"campaignId" bson:"campaignId"` // UUID campaign
	OwnerId    primitive.ObjectID `json:"ownerId" bson:"ownerId"`       // UUID user
}

type Shop struct {
	Id   primitive.ObjectID `json:"id" bson:"_id"` // UUID
	Name string             `json:"name" bson:"name"`

	CampaignId primitive.ObjectID `json:"campaignId" bson:"campaignId"` // UUID campaign
}

type Item struct {
	Id   primitive.ObjectID `json:"id" bson:"_id"` // UUID
	Name string             `json:"name" bson:"name"`

	ShopId primitive.ObjectID `json:"shopId" bson:"shopId"` // UUID shop
	Cost   Balance            `json:"cost" bson:"cost"`
}

// Not a collection
type Balance struct {
	Copper   uint `json:"copper" bson:"copper"`
	Silver   uint `json:"silver" bson:"silver"`
	Electrum uint `json:"electrum" bson:"electrum"`
	Gold     uint `json:"gold" bson:"gold"`
	Platinum uint `json:"platinum" bson:"platinum"`
}
