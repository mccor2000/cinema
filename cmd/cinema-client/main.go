package main

import (
	"context"
	"log"
	"time"

	pb "github.com/mccor2000/cinema/proto"
	"google.golang.org/grpc"
)

func main() {
  opts := []grpc.DialOption{}
  conn, err := grpc.NewClient("localhost:50051", opts...)

	if err != nil {
		log.Fatalf("Failed to not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCinemaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Example usage
	config := &pb.CinemaConfig{Rows: 5, Columns: 5, MinDistance: 2}

	// Query available seats
	queryResp, err := c.QueryAvailableSeats(ctx, &pb.QueryRequest{Config: config, GroupSize: 2})
	if err != nil {
		log.Fatalf("could not query available seats: %v", err)
	}
	log.Printf("Available seats: %v", queryResp.AvailableSeats)

	// Reserve seats
	seatsToReserve := []*pb.Seat{
		{Row: 0, Column: 0},
		{Row: 0, Column: 1},
	}
	reserveResp, err := c.ReserveSeats(ctx, &pb.ReservationRequest{Config: config, Seats: seatsToReserve})
	if err != nil {
		log.Fatalf("could not reserve seats: %v", err)
	}
	log.Printf("Reservation response: %v, %s", reserveResp.Success, reserveResp.Message)

	// Cancel reservation
	cancelResp, err := c.CancelReservation(ctx, &pb.CancellationRequest{Config: config, Seats: seatsToReserve})
	if err != nil {
		log.Fatalf("could not cancel reservation: %v", err)
	}
	log.Printf("Cancellation response: %v, %s", cancelResp.Success, cancelResp.Message)
}
