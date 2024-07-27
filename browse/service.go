package browse

import (
	"context"

	"github.com/yourusername/bookstore/proto"
	"github.com/yourusername/bookstore/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	proto.UnimplementedBrowseServiceServer
	db *mongo.Database
}

func NewServer(db *mongo.Database) *Server {
	return &Server{db: db}
}

func (s *Server) BrowseBooks(ctx context.Context, req *proto.BrowseBooksRequest) (*proto.BrowseBooksResponse, error) {
	skip := int64((req.Page - 1) * req.PageSize)
	limit := int64(req.PageSize)

	opts := options.Find().SetSkip(skip).SetLimit(limit)
	cursor, err := s.db.Collection("books").Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var books []*proto.Book
	if err = cursor.All(ctx, &books); err != nil {
		return nil, err
	}

	totalBooks, err := s.db.Collection("books").CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	return &proto.BrowseBooksResponse{
		Books:      books,
		TotalBooks: int32(totalBooks),
	}, nil
}
