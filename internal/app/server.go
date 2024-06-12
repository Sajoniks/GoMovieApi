package app

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"sajoniks.github.io/movieApi/internal/infra/repo"
	"sajoniks.github.io/movieApi/pkg/proto-gen"
)

type movieApiServer struct {
	proto_gen.UnimplementedMovieApiServer
	repo *repo.MovieRepository
}

func NewMovieApiServer() (proto_gen.MovieApiServer, error) {
	r, err := repo.NewMovieRepository("postgres://movie_api:123@localhost:5432/movie_api")
	if err != nil {
		return nil, errors.Join(errors.New("failed to create server"), err)
	}
	return &movieApiServer{
		repo: r,
	}, nil
}

type getFilmsRequestPageToken struct {
	StudioId int `json:"studio_id"`
	Offset   int `json:"offset"`
}

func (s *movieApiServer) GetAllMovies(r *proto_gen.GetFilmsRequest, stream proto_gen.MovieApi_GetAllMoviesServer) error {
	pageTokenData := getFilmsRequestPageToken{
		StudioId: 0,
		Offset:   0,
	}
	if len(r.PageToken) > 0 {
		decodedToken, decodeErr := base64.StdEncoding.DecodeString(r.PageToken)
		if decodeErr != nil {
			return errors.Join(errors.New("page token corrupted"), decodeErr)
		}

		decodeErr = json.Unmarshal(decodedToken, &pageTokenData)
		if decodeErr != nil {
			return errors.Join(errors.New("page token corrupted"), decodeErr)
		}
	}

	req := repo.GetAllMoviesRequest{
		StudioId: pageTokenData.StudioId,
		Pagination: repo.Pagination{
			Size: int(r.PageSize),
			Page: pageTokenData.Offset,
		},
	}
	movies, getErr := s.repo.GetAllMovies(context.Background(), req)
	if getErr != nil {
		return errors.Join(errors.New("could not get movies"), getErr)
	}

	pageTokenData.Offset += len(movies)
	films := make([]*proto_gen.GetFilmsFilm, len(movies))
	for i, _ := range movies {
		films[i] = &proto_gen.GetFilmsFilm{}
	}

	var nextPageToken string
	if len(movies) > 0 {
		b, _ := json.Marshal(&pageTokenData)
		nextPageToken = base64.StdEncoding.EncodeToString(b)
	}
	sendErr := stream.Send(&proto_gen.GetFilmsResponse{
		Films:         films,
		NextPageToken: nextPageToken,
	})

	return sendErr
}
