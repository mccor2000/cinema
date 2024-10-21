package storage

import (
	"sync"
)

// Storage is an interface for storing cinema data.
type InMemoryStorage struct {
	cinema	 *Cinema
	mu         sync.RWMutex
}

// NewInMemoryStorage creates a new instance of InMemoryStorage.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		cinema: &Cinema{},
		mu: sync.RWMutex{},
	}
}

// Get returns a cinema from the storage.
func (s *InMemoryStorage) Get() *Cinema {
	return s.cinema
}

// Update updates the configuration of a cinema.
func (s *InMemoryStorage) UpdateConf(conf CinemaConfig) (*Cinema, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cinema.Conf = conf

	// reset available seats
	s.cinema.AvailableSeats = make([][2]int32, 0)
	for i := int32(0); i < conf.Rows; i++ {
		for j := int32(0); j < conf.Columns; j++ {
			s.cinema.AvailableSeats = append(s.cinema.AvailableSeats, [2]int32{i, j})
		}
	}

	return s.cinema, nil
}

// UpdateSeats updates the seats of a cinema.
func (s *InMemoryStorage) UpdateAvailableSeats(seats [][2]int32) (*Cinema, error) {
	// Let the service handle the locking when dealing with concurrent requests.
	// s.mu.Lock()
	// defer s.mu.Unlock()
	//
	s.cinema.AvailableSeats = seats

	return s.cinema, nil
}

// Helper functions to lock and unlock the storage.
func (s *InMemoryStorage) Lock() {
	s.mu.Lock()
}
func (s *InMemoryStorage) Unlock() {
	s.mu.Unlock()
}
