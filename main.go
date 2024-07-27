package main

import (
	"log"
	"net"

	"github.com/maralski/bookstore/proto"
	"github.com/maralski/bookstore/web"
	"github.com/maralski/bookstore/purchase"
	"github.com/maralski/bookstore/search"
	"github.com/maralski/bookstore/browse"
	"github.com/maralski/bookstore/database"
	"google.golang.org/grpc"
)

func main() {
	// Initialize database connection
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register services
	proto.RegisterWebServiceServer(grpcServer, web.NewServer(db))
	proto.RegisterPurchaseServiceServer(grpcServer, purchase.NewServer(db))
	proto.RegisterSearchServiceServer(grpcServer, search.NewServer(db))
	proto.RegisterBrowseServiceServer(grpcServer, browse.NewServer(db))

	// Start gRPC server
	grpcAddr := ":50051"
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("gRPC server listening on %s", grpcAddr)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Create and start HTTP server
	httpServer, err := web.NewHTTPServer(grpcAddr)
	if err != nil {
		log.Fatalf("Failed to create HTTP server: %v", err)
	}
	if err := httpServer.Start("8080"); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
