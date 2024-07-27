package search

import (
	"context"

	"github.com/yourusername/bookstore/proto"
	"github.com/yourusername/bookstore/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	proto.UnimplementedSearchServiceServer
	db *mongo.Database
}

func NewServer(db *mongo.Database) *Server {
	return &Server{db: db}
}

func (s *Server) SearchBooks(ctx context.Context, req *proto.SearchBooksRequest) (*proto.SearchBooksResponse, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": req.Query, "$options": "i"}},
			{"author": bson.M{"$regex": req.Query, "$options": "i"}},
		},
	}

	cursor, err := s.db.Collection("books").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var books []*proto.Book
	if err = cursor.All(ctx, &books); err != nil {
		return nil, err
	}

	return &proto.SearchBooksResponse{Books: books}, nil
}
