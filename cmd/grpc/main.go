package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	h "github.com/mccor2000/cinema/pkg/handlers"
	pb "github.com/mccor2000/cinema/proto"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCinemaServiceServer(s, h.NewGrpcHandler("in_memory"))
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
