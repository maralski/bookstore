package purchase

import (
	"context"

	"github.com/yourusername/bookstore/proto"
	"github.com/yourusername/bookstore/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	proto.UnimplementedPurchaseServiceServer
	db *mongo.Database
}

func NewServer(db *mongo.Database) *Server {
	return &Server{db: db}
}

func (s *Server) PurchaseBook(ctx context.Context, req *proto.PurchaseBookRequest) (*proto.PurchaseBookResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.BookId)
	if err != nil {
		return nil, err
	}

	// Implement purchase logic here
	// For simplicity, we'll just check if the book exists

	var book proto.Book
	err = s.db.Collection("books").FindOne(ctx, bson.M{"_id": id}).Decode(&book)
	if err != nil {
		return &proto.PurchaseBookResponse{Success: false}, nil
	}

	// Generate a new order ID
	orderId := primitive.NewObjectID().Hex()

	return &proto.PurchaseBookResponse{
		Success: true,
		OrderId: orderId,
	}, nil
}
