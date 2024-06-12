package app

import (
	"context"
	"errors"
	"sajoniks.github.io/movieApi/internal/infra/repo"
	"sajoniks.github.io/movieApi/internal/infra/repo/interface"
	"sajoniks.github.io/movieApi/pkg/proto-gen"
)

type movieApiServer struct {
	proto_gen.UnimplementedMovieApiServer
	repo repointerface.MovieRepository
}

func NewMovieApiServer() (proto_gen.MovieApiServer, error) {
	r, err := repo.NewMovieRepo(context.Background(), "postgres://movie_api:123@localhost:5432/movie_api")
	if err != nil {
		return nil, errors.Join(errors.New("failed to create server"), err)
	}
	return &movieApiServer{
		repo: r,
	}, nil
}

func (s *movieApiServer) GetAllMovies(r *proto_gen.GetFilmsRequest, stream proto_gen.MovieApi_GetAllMoviesServer) error {
	return nil
}
