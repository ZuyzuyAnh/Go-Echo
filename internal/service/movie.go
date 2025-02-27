package service

import (
	"echo-demo/internal/dto"
	"echo-demo/internal/model"
	"echo-demo/internal/repository"
)

type MovieService struct {
	MovieRepository *repository.MovieRepository
}

func NewMovieService(ur *repository.MovieRepository) *MovieService {
	return &MovieService{
		MovieRepository: ur,
	}
}

func (m *MovieService) CreateMovie(req *dto.CreateMovieReq) (*dto.MovieResp, error) {
	movie := model.Movie{
		Title:         req.Title,
		Description:   req.Description,
		Duration:      req.Duration,
		CoverURL:      req.CoverURL,
		BackgroundURL: req.BackgroundURL,
	}

	err := m.MovieRepository.InsertMovie(nil, &movie)
	if err != nil {
		return nil, err
	}

	return &dto.MovieResp{
		ID:            movie.ID,
		Title:         movie.Title,
		Description:   movie.Description,
		Duration:      movie.Duration,
		CoverURL:      movie.CoverURL,
		BackgroundURL: movie.BackgroundURL,
	}, nil
}
