package model

type Seat struct {
	ID         int64
	Number     string
	TheaterID  int64
	SeatTypeID int64
}

type SeatType struct {
	ID          int64
	Description string
	Price       float64
}
