package handler

import (
	"context"
	"log"

	srv "github.com/mccor2000/cinema/pkg/service"
	str "github.com/mccor2000/cinema/pkg/storage"
	pb "github.com/mccor2000/cinema/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcHandler struct {
	pb.UnimplementedCinemaServiceServer
	srv *srv.Service
}

func NewGrpcHandler() *GrpcHandler {
	return &GrpcHandler{
		srv: srv.NewService(str.NewInMemoryStorage()),
	}
}

func (h *GrpcHandler) UpdateCinema(ctx context.Context, in *pb.CinemaConfig) (*pb.CinemaResponse, error) {
	cinema, err := h.srv.UpdateCinema(str.CinemaConfig{
		Rows:        in.Rows,
		Columns:     in.Columns,
		MinDistance: in.MinDistance,
	})
	if err != nil {
		log.Printf("failed to update cinema: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update cinema: %v", err)
	}
	return &pb.CinemaResponse{
		Rows:        cinema.Conf.Rows,
		Columns:     cinema.Conf.Columns,
		MinDistance: cinema.Conf.MinDistance,
	}, nil
}

func (h *GrpcHandler) QueryAvailableSeats(ctx context.Context, in *pb.QueryRequest) (*pb.QueryResponse, error) {
	seats, err := h.srv.GetAvailableSeats()
	if err != nil {
		log.Printf("failed to get available seats: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get available seats: %v", err)
	}
	return &pb.QueryResponse{AvailableSeats: toPbSeats(seats)}, nil
}

func (h *GrpcHandler) ReserveSeats(ctx context.Context, in *pb.ReservationRequest) (*pb.ReservationResponse, error) {
	seats := toSeats(in.Seats)
	err := h.srv.ReserveSeats(seats)
	if err != nil {
		log.Printf("failed to reserve seats: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to reserve seats: %v", err)
	}
	return &pb.ReservationResponse{
		Success: true,
		Message: "Seats reserved successfully",
	}, nil
}

func (h *GrpcHandler) CancelReservation(ctx context.Context, in *pb.CancellationRequest) (*pb.CancellationResponse, error) {
	seats := toSeats(in.Seats)
	err := h.srv.CancelSeats(seats)
	if err != nil {
		log.Printf("failed to cancel reservation: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to cancel reservation: %v", err)
	}
	return &pb.CancellationResponse{
		Success: true,
		Message: "Reservation cancelled successfully",
	}, nil
}

// helper to convert a list of [2]int32 to a list of pb.Seat and vice versa
func toPbSeats(seats [][2]int32) []*pb.Seat {
	res := make([]*pb.Seat, 0)
	for _, seat := range seats {
		res = append(res, &pb.Seat{Row: seat[0], Column: seat[1]})
	}
	return res
}

func toSeats(seat []*pb.Seat) [][2]int32 {
	res := make([][2]int32, 0)
	for _, s := range seat {
		res = append(res, [2]int32{s.Row, s.Column})
	}
	return res
}
