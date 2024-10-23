package main

import (
	"context"
	"log"
	"time"

	pb "github.com/mccor2000/cinema/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
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
	updateResp, err := c.UpdateCinema(ctx, config)
	if err != nil {
		log.Fatalf("could not update cinema: %v", err)
	}
	log.Printf("Updated cinema: %v", updateResp)

	// Query available seats
	queryResp, err := c.QueryAvailableSeats(ctx, &pb.QueryRequest{})
	if err != nil {
		log.Fatalf("could not query available seats: %v", err)
	}
	log.Printf("Available seats: %v", queryResp.AvailableSeats)

	// Reserve seats
	seatsToReserve := []*pb.Seat{
		{Row: 0, Column: 0},
		{Row: 0, Column: 1},
	}
	reserveResp, err := c.ReserveSeats(ctx, &pb.ReservationRequest{Seats: seatsToReserve})
	if err != nil {
		log.Fatalf("could not reserve seats: %v", err)
	}
	log.Printf("Reservation response: %v, %s", reserveResp.Success, reserveResp.Message)

	// Cancel reservation
	cancelResp, err := c.CancelReservation(ctx, &pb.CancellationRequest{Seats: seatsToReserve})
	if err != nil {
		log.Fatalf("could not cancel reservation: %v", err)
	}
	log.Printf("Cancellation response: %v, %s", cancelResp.Success, cancelResp.Message)
}
