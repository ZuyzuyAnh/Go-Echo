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

func (m *MovieService) GetMovieByID(id int64) (*dto.MovieResp, error) {
	movie, err := m.MovieRepository.GetMovieByID(nil, id)
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

func (m *MovieService) UpdateMovie(req *dto.CreateMovieReq, id int64) error {
	movie := model.Movie{
		ID: id,
	}

	if req.Title != "" {
		movie.Title = req.Title
	}
	if req.Description != "" {
		movie.Description = req.Description
	}
	if req.Duration != 0 {
		movie.Duration = req.Duration
	}
	if req.CoverURL != "" {
		movie.CoverURL = req.CoverURL
	}
	if req.BackgroundURL != "" {
		movie.BackgroundURL = req.BackgroundURL
	}

	err := m.MovieRepository.UpdateMovie(nil, &movie)

	return err
}

func (m *MovieService) DeleteMovie(id int64) error {
	err := m.MovieRepository.DeleteMovie(nil, id)
	return err
}
