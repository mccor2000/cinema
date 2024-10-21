package service

import (
	"errors"
	"math"

	"github.com/mccor2000/cinema/pkg/storage"
)

type Service struct {
	storage storage.Storage
}

func NewService(s storage.Storage) *Service {
	return &Service{
		storage: s,
	}
}

func (s *Service) UpdateCinema(conf storage.CinemaConfig) (*storage.Cinema, error) {
	cinema, err := s.storage.UpdateConf(conf)
	if err != nil {
		return nil, err
	}

	// Initialize available seats after updating configuration
	availableSeats := s.initializeAvailableSeats(cinema.Conf)
	_, err = s.storage.UpdateAvailableSeats(availableSeats)
	if err != nil {
		return nil, err
	}

	return cinema, nil
}

func (s *Service) GetAvailableSeats() ([][2]int32, error) {
	cinema := s.storage.Get()
	return cinema.AvailableSeats, nil
}

func (s *Service) ReserveSeats(seats [][2]int32) error {
	cinema := s.storage.Get()
	s.storage.Lock()
	defer s.storage.Unlock()

	// Check if all seats are valid and available
	for _, seat := range seats {
		if !s.isValidSeat(seat, cinema.Conf) {
			return errors.New("invalid seat")
		}
		if !s.isSeatAvailable(seat, cinema.AvailableSeats) {
			return errors.New("seat not available")
		}
	}

	// Update available seats
	newAvailableSeats := s.updateAvailableSeats(cinema.AvailableSeats, seats, cinema.Conf.MinDistance)

	if _, err := s.storage.UpdateAvailableSeats(newAvailableSeats); err != nil {
		return err
	}

	return nil
}

func (s *Service) CancelSeats(seats [][2]int32) error {
	cinema := s.storage.Get()
	s.storage.Lock()
	defer s.storage.Unlock()

	// Check if all seats are valid
	for _, seat := range seats {
		if !s.isValidSeat(seat, cinema.Conf) {
			return errors.New("invalid seat")
		}
	}

	// Add cancelled seats back to available seats
	newAvailableSeats := append(cinema.AvailableSeats, seats...)

	// Remove duplicates and seats that are too close to occupied seats
	newAvailableSeats = s.filterAvailableSeats(newAvailableSeats, cinema.Conf.MinDistance)

	if _, err := s.storage.UpdateAvailableSeats(newAvailableSeats); err != nil {
		return err
	}

	return nil
}

func (s *Service) manhattanDistance(seat1, seat2 [2]int32) int32 {
	return int32(math.Abs(float64(seat1[0]-seat2[0])) + math.Abs(float64(seat1[1]-seat2[1])) - 1)
}

func (s *Service) isValidSeat(seat [2]int32, config storage.CinemaConfig) bool {
	return seat[0] >= 0 && seat[0] < config.Rows && seat[1] >= 0 && seat[1] < config.Columns
}

func (s *Service) isSeatAvailable(seat [2]int32, availableSeats [][2]int32) bool {
	for _, availableSeat := range availableSeats {
		if seat == availableSeat {
			return true
		}
	}
	return false
}

func (s *Service) initializeAvailableSeats(config storage.CinemaConfig) [][2]int32 {
	var seats [][2]int32
	for row := int32(0); row < config.Rows; row++ {
		for col := int32(0); col < config.Columns; col++ {
			seats = append(seats, [2]int32{row, col})
		}
	}
	return seats
}

func (s *Service) updateAvailableSeats(availableSeats, reservedSeats [][2]int32, minDistance int32) [][2]int32 {
	var newAvailableSeats [][2]int32

	for _, availableSeat := range availableSeats {
		if !s.isSeatReserved(availableSeat, reservedSeats) && !s.isSeatTooClose(availableSeat, reservedSeats, minDistance) {
			newAvailableSeats = append(newAvailableSeats, availableSeat)
		}
	}

	return newAvailableSeats
}

func (s *Service) isSeatReserved(seat [2]int32, reservedSeats [][2]int32) bool {
	for _, reservedSeat := range reservedSeats {
		if seat == reservedSeat {
			return true
		}
	}
	return false
}

func (s *Service) isSeatTooClose(seat [2]int32, reservedSeats [][2]int32, minDistance int32) bool {
	for _, reservedSeat := range reservedSeats {
		if s.manhattanDistance(seat, reservedSeat) < minDistance {
			return true
		}
	}
	return false
}

func (s *Service) filterAvailableSeats(seats [][2]int32, minDistance int32) [][2]int32 {
	var filteredSeats [][2]int32
	occupiedSeats := make(map[[2]int32]bool)

	for _, seat := range seats {
		if !occupiedSeats[seat] && !s.isSeatTooClose(seat, filteredSeats, minDistance) {
			filteredSeats = append(filteredSeats, seat)
			occupiedSeats[seat] = true
		}
	}

	return filteredSeats
}
