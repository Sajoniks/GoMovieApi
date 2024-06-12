package repointerface

import (
	"context"
	"sajoniks.github.io/movieApi/internal/domain/agg"
	"sajoniks.github.io/movieApi/internal/infra/repo/query/interface"
)

type MovieRepository interface {
	GetMovieCast(ctx context.Context, q queryinterface.GetMovieCastById) (*agg.MovieCast, error)
	GetMovieById(ctx context.Context, q queryinterface.GetMovieById) (*agg.Movie, error)
	GetMovieByName(ctx context.Context, q queryinterface.GetMovieByName) (*agg.Movie, error)
	GetAllMovies(ctx context.Context, q queryinterface.GetAllMovies) ([]*agg.Movie, error)
	Update(ctx context.Context, movie *agg.Movie) error
	Delete(ctx context.Context, movie *agg.Movie) error
	InsertAll(ctx context.Context, movie []*agg.Movie) error
}
