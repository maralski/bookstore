package web

import (
	"context"

	"github.com/yourusername/bookstore/proto"
	"github.com/yourusername/bookstore/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	proto.UnimplementedWebServiceServer
	db *mongo.Database
}

func NewServer(db *mongo.Database) *Server {
	return &Server{db: db}
}

func (s *Server) GetBook(ctx context.Context, req *proto.GetBookRequest) (*proto.Book, error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	var book proto.Book
	err = s.db.Collection("books").FindOne(ctx, bson.M{"_id": id}).Decode(&book)
	if err != nil {
		return nil, err
	}

	return &book, nil
}
