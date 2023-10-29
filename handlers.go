package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shuaibu222/p-products/products"
	"go.mongodb.org/mongo-driver/bson"
)

// here should be the grpc handlers
type ProductsServer struct {
	products.UnimplementedProductsServiceServer
}

func (p *ProductsServer) CreateProducts(ctx context.Context, req *products.ProductRequest) (*products.ProductResponse, error) {
	input := req.GetProductEntry()

	// create the product
	productEntry := ProductsEntry{
		Name:        input.Name,
		Motto:       input.Motto,
		Link:        input.Link,
		Category:    input.Category,
		Description: input.Description,
		ImageLink:   input.ImageLink,
	}

	pro, err := productEntry.Insert(productEntry, ctx)
	if err != nil {
		log.Printf("Error inserting product: %v", err)
	}

	log.Println("Product inserted successfully: ", pro)

	var entry ProductsEntry
	err = collection.FindOne(ctx, bson.M{"_id": pro.InsertedID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	insertedProduct := &products.Product{
		Id:          fmt.Sprintf("%v", entry.ID),
		Name:        entry.Name,
		Motto:       entry.Motto,
		Link:        entry.Link,
		Category:    entry.Category,
		Description: entry.Description,
		ImageLink:   entry.ImageLink,
	}
	// return response
	res := &products.ProductResponse{Response: insertedProduct}
	return res, nil
}

func (u *ProductsServer) GetAllProducts(req *products.NoParams, stream products.ProductsService_GetAllProductsServer) error {
	productEntry := ProductsEntry{}

	productsResult, err := productEntry.GetAllProducts(context.Background())
	if err != nil {
		log.Println("Failed to get products", err)
	}

	for _, productEntry := range productsResult {
		productResponse := &products.Product{
			Id:          productEntry.ID.Hex(),
			Name:        productEntry.Name,
			Motto:       productEntry.Motto,
			Link:        productEntry.Link,
			Category:    productEntry.Category,
			Description: productEntry.Description,
			ImageLink:   productEntry.ImageLink,
		}

		if err := stream.Send(productResponse); err != nil {
			log.Println("Error sending product to the client:", err)
			return err
		}
	}

	// Return the response
	return nil
}

func (u *ProductsServer) GetProductById(ctx context.Context, req *products.ProductId) (*products.ProductResponse, error) {
	productEntry := ProductsEntry{}

	result, err := productEntry.GetOneProductById(req.Id, ctx)
	if err != nil {
		log.Println(err)
	}

	productResult := &products.ProductResponse{
		Response: &products.Product{
			Id:          result.ID.String(),
			Name:        result.Name,
			Motto:       result.Motto,
			Link:        result.Link,
			Category:    result.Category,
			Description: result.Description,
			ImageLink:   result.ImageLink,
		},
	}

	return productResult, nil
}

func (u *ProductsServer) UpdateProduct(ctx context.Context, req *products.Product) (*products.Count, error) {
	productEntry := ProductsEntry{
		Name:        req.Name,
		Motto:       req.Motto,
		Link:        req.Link,
		Category:    req.Category,
		Description: req.Description,
		ImageLink:   req.ImageLink,
	}

	updateCount, err := productEntry.UpdateProductById(req.Id, productEntry, ctx)
	if err != nil {
		log.Println(err)
	}

	log.Println(updateCount.ModifiedCount)

	productResult := &products.Count{
		Count: fmt.Sprint(updateCount.ModifiedCount),
	}

	return productResult, nil
}

func (u *ProductsServer) DeleteProduct(ctx context.Context, req *products.ProductId) (*products.Count, error) {
	productEntry := ProductsEntry{}

	deletedCount, err := productEntry.DeleteProduct(req.Id, ctx)
	if err != nil {
		log.Println(err)
	}

	res := &products.Count{
		Count: fmt.Sprint(deletedCount.DeletedCount),
	}

	return res, nil
}
