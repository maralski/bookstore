package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yourusername/bookstore/proto"
	"google.golang.org/grpc"
)

type HTTPServer struct {
	webClient       proto.WebServiceClient
	purchaseClient  proto.PurchaseServiceClient
	searchClient    proto.SearchServiceClient
	browseClient    proto.BrowseServiceClient
}

func NewHTTPServer(grpcServerAddr string) (*HTTPServer, error) {
	conn, err := grpc.Dial(grpcServerAddr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %v", err)
	}

	return &HTTPServer{
		webClient:       proto.NewWebServiceClient(conn),
		purchaseClient:  proto.NewPurchaseServiceClient(conn),
		searchClient:    proto.NewSearchServiceClient(conn),
		browseClient:    proto.NewBrowseServiceClient(conn),
	}, nil
}

func (s *HTTPServer) Start(port string) error {
	r := mux.NewRouter()

	r.HandleFunc("/book/{id}", s.handleGetBook).Methods("GET")
	r.HandleFunc("/book/{id}/purchase", s.handlePurchaseBook).Methods("POST")
	r.HandleFunc("/books/search", s.handleSearchBooks).Methods("GET")
	r.HandleFunc("/books/browse", s.handleBrowseBooks).Methods("GET")

	log.Printf("HTTP server listening on :%s", port)
	return http.ListenAndServe(":"+port, r)
}

func (s *HTTPServer) handleGetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["id"]

	book, err := s.webClient.GetBook(context.Background(), &proto.GetBookRequest{Id: bookID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(book)
}

func (s *HTTPServer) handlePurchaseBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["id"]

	var request struct {
		Quantity int32 `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := s.purchaseClient.PurchaseBook(context.Background(), &proto.PurchaseBookRequest{
		BookId:   bookID,
		Quantity: request.Quantity,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

func (s *HTTPServer) handleSearchBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	response, err := s.searchClient.SearchBooks(context.Background(), &proto.SearchBooksRequest{Query: query})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

func (s *HTTPServer) handleBrowseBooks(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 {
		pageSize = 10
	}

	response, err := s.browseClient.BrowseBooks(context.Background(), &proto.BrowseBooksRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}
