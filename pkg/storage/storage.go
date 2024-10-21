package storage

type Storage interface {
	// Manage cinemas
	Get() (*Cinema)
	UpdateConf( conf CinemaConfig) (*Cinema, error)
	// Manage seats
	UpdateAvailableSeats(seats [][2]int32) (*Cinema, error)
	// helpers
	Lock()
	Unlock()
}

type CinemaConfig struct {
	Rows        int32
	Columns     int32
	MinDistance int32
}

type Cinema struct {
	Conf           CinemaConfig
	AvailableSeats [][2]int32
}
