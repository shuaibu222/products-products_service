package main

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

const mongoURL = "mongodb://mongo:27017"

// connect with MongoDB
func init() {
	credential := options.Credential{
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	}
	clientOpts := options.Client().ApplyURI(mongoURL).SetAuth(credential)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Println("Error connecting to MongoDB")
		return
	}

	collection = client.Database("products").Collection("products")

	// collection instance
	log.Println("Collections instance is ready")
}

type ProductsEntry struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	Motto       string             `bson:"motto" json:"motto"`
	Link        string             `bson:"link" json:"link"`
	Category    string             `bson:"category" json:"category"`
	Description string             `bson:"description" json:"description"`
	ImageLink   string             `bson:"image_link" json:"image_link"`
}

func (u *ProductsEntry) Insert(entry ProductsEntry, ctx context.Context) (*mongo.InsertOneResult, error) {
	product, err := collection.InsertOne(ctx, entry)
	if err != nil {
		log.Println("Error inserting into products:", err)
	}

	return product, nil
}

func (u *ProductsEntry) GetAllProducts(ctx context.Context) ([]*ProductsEntry, error) {

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Finding all products error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*ProductsEntry

	for cursor.Next(ctx) {
		var item ProductsEntry

		err := cursor.Decode(&item)
		if err != nil {
			log.Print("Error decoding log into slice:", err)
			return nil, err
		} else {
			products = append(products, &item)
		}
	}

	return products, nil
}

func (u *ProductsEntry) GetOneProductById(id string, ctx context.Context) (*ProductsEntry, error) {
	var product *ProductsEntry
	Id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Failed to convert id", err)
		return nil, err
	}

	filter := bson.M{"_id": Id}
	err = collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		log.Println("Failed to get product", err)
		return nil, err
	}

	return product, nil
}

func (u *ProductsEntry) UpdateProductById(id string, productToUpdate ProductsEntry, ctx context.Context) (*mongo.UpdateResult, error) {
	Id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Failed to convert id", err)
		return nil, err
	}

	updatedCount, err := collection.UpdateOne(ctx, bson.M{"_id": Id}, bson.M{"$set": productToUpdate})
	if err != nil {
		log.Println("Failed to update the product: ", err)
		return nil, err
	}
	return updatedCount, nil
}

func (u *ProductsEntry) DeleteProduct(id string, ctx context.Context) (*mongo.DeleteResult, error) {
	Id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Failed to convert id", err)
		return nil, err
	}

	deleteCount, err := collection.DeleteOne(ctx, bson.M{"_id": Id})
	if err != nil {
		log.Println("Failed to delete product", err)
		return nil, err
	}

	return deleteCount, nil
}
