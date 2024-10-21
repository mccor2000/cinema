package service

import (
	"fmt"
	"testing"

	"github.com/mccor2000/cinema/pkg/storage"
)

func TestReserveSeats(t *testing.T) {
	fmt.Println("TestReserveSeats")
	// Given a 5x5 cinema with min distance is 1
	s := NewService(storage.NewInMemoryStorage())
	fmt.Println("TestReserveSeats")
	conf := storage.CinemaConfig{
		Rows:        5,
		Columns:     5,
		MinDistance: 1,
	}
	_, err := s.UpdateCinema(conf)
	fmt.Println("TestReserveSeats")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Reserve 3 seats
	seats := [][2]int32{{0, 0}, {0, 1}, {0, 2}}
	err = s.ReserveSeats(seats)
	fmt.Println("TestReserveSeats")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	availableSeats, err := s.GetAvailableSeats()
	fmt.Println("TestReserveSeats")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// Expected result: 
	// x is reserved, v is available, - is not available
	//
	// x  x  x	-  v
	// -	-  -	v  v
	// v	v  v	v  v
	// v	v  v	v  v
	// v	v  v	v  v
	//
	// There should be 18 available seats
	if len(availableSeats) != 18 {
		t.Errorf("unexpected number of available seats: %v", availableSeats)
	}
	// First available seat should be [0, 4]
	if availableSeats[0][0] != 0 || availableSeats[0][1] != 4 {
		t.Errorf("unexpected available seats: %v", availableSeats)
	}
	// Second available seat should be [1, 3]
	if availableSeats[1][0] != 1 || availableSeats[1][1] != 3 {
		t.Errorf("unexpected available seats: %v", availableSeats)
	}

	// Reserve 2 more seats
	seats = [][2]int32{{1, 3}, {1, 4}}
	err = s.ReserveSeats(seats)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	availableSeats, err = s.GetAvailableSeats()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// Expected result:
	// x is reserved, v is available, - is not available
	//
	// x  x  x  -  -
	// -  -  -  x  x
	// v  v  v  -  -
	// v  v  v  v  v
	// v  v  v  v  v
	//
	// There should be 13 available seats
	if len(availableSeats) != 13 {
		t.Errorf("unexpected number of available seats: %v", availableSeats)
	}
	// First available seat should be [2, 0]
	if availableSeats[0][0] != 2 || availableSeats[0][1] != 0 {
		t.Errorf("unexpected available seats: %v", availableSeats)
	}
}
