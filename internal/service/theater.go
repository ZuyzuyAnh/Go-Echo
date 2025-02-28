package service

import (
	"echo-demo/internal/dto"
	"echo-demo/internal/model"
	"echo-demo/internal/repository"
	"fmt"
)

type TheaterService struct {
	TheaterRepository *repository.TheaterRepository
	SeatRepository    *repository.SeatRepository
}

func NewTheaterService(tr *repository.TheaterRepository, sr *repository.SeatRepository) *TheaterService {
	return &TheaterService{
		TheaterRepository: tr,
		SeatRepository:    sr,
	}
}

func (s *TheaterService) CreateTheater(req *dto.CreateUpdateTheaterRequest) (*dto.TheaterResponse, error) {
	tx, err := s.TheaterRepository.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	theater := model.Theater{
		Name: req.Name,
	}

	if err = s.TheaterRepository.CreateTheater(tx, &theater); err != nil {
		return nil, err
	}

	for i := 0; i < req.Rows; i++ {
		for j := 1; j <= req.Cols; j++ {
			rowLetter := string(rune('A' + i))
			number := fmt.Sprintf("%s%d", rowLetter, j)

			seat := model.Seat{
				Number:    number,
				TheaterID: theater.ID,
			}

			if err = s.SeatRepository.CreateSeat(tx, &seat); err != nil {
				return nil, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &dto.TheaterResponse{
		ID:   theater.ID,
		Name: theater.Name,
	}, nil
}

func (s *TheaterService) GetTheaterByID(id int64) (*dto.TheaterResponse, error) {
	theater, err := s.TheaterRepository.GetTheaterByID(nil, id)
	if err != nil {
		return nil, err
	}

	return &dto.TheaterResponse{
		ID:   theater.ID,
		Name: theater.Name,
	}, nil
}

func (s *TheaterService) ListTheaters() ([]*dto.TheaterResponse, error) {
	theaters, err := s.TheaterRepository.ListTheaters(nil)
	if err != nil {
		return nil, err
	}

	var responses []*dto.TheaterResponse
	for _, t := range theaters {
		responses = append(responses, &dto.TheaterResponse{
			ID:   t.ID,
			Name: t.Name,
		})
	}

	return responses, nil
}

func (s *TheaterService) UpdateTheater(req *dto.CreateUpdateTheaterRequest, id int64) error {
	theater := model.Theater{
		ID:   id,
		Name: req.Name,
	}

	return s.TheaterRepository.UpdateTheater(nil, &theater)
}

func (s *TheaterService) DeleteTheater(id int64) error {
	return s.TheaterRepository.DeleteTheater(nil, id)
}

func (s *TheaterService) UpdateSeatType(req *dto.UpdateSeatTypeRequest) error {
	return s.SeatRepository.BatchUpdateSeatType(nil, req.SeatTypeID, req.SeatIDs)
}
