package repo

import (
	"context"
	"fmt"
	"sajoniks.github.io/movieApi/internal/domain/agg"
	"sajoniks.github.io/movieApi/internal/domain/dto"
	"sajoniks.github.io/movieApi/internal/domain/entity"
	"sajoniks.github.io/movieApi/internal/infra/repo/interface"
	"testing"
)

var repo repointerface.MovieRepository

func TestMain(m *testing.M) {
	var err error
	repo, err = NewMovieRepo(context.Background(), "postgres://movie_api:123@localhost:5432/movie_api")
	if err != nil {
		panic(err)
	}
	m.Run()

}

func Test_MovieRepo_GetMovieById(t *testing.T) {
	r := &dto.GetMovieDto{Id: 1}
	m, err := repo.GetMovieById(context.Background(), r)
	if err != nil {
		t.Fatalf("want movie with id %d, got error %v", r.GetMovieId(), err)
	}
	if m == nil {
		t.Fatal("want movie, got nil")
	}
	if m.Id != r.GetMovieId() {
		t.Errorf("want movie with id %d, got %d", r.GetMovieId(), m.Id)
	}
}

func Test_MovieRepo_GetMovieByName(t *testing.T) {
	r := &dto.GetMovieDto{Name: "Django Unchained"}
	m, err := repo.GetMovieByName(context.Background(), r)
	if err != nil {
		t.Fatalf("did not want an error %v", err)
	}
	if m == nil {
		t.Fatalf("got nil result")
	}
	if m.Name != r.GetMovieName() {
		t.Errorf("want movie with name %q, got %q", r.GetMovieName(), m.Name)
	}
}

func Test_MovieRepo_GetMovieCastById(t *testing.T) {
	r := &dto.GetMovieCastByIdDto{Id: 1}
	cast, err := repo.GetMovieCast(context.Background(), r)
	if err != nil {
		t.Fatalf("did not want an error %v", err)
	}

	if cast == nil {
		t.Fatalf("got nil result")
	}

	if len(cast.Cast) == 0 {
		t.Errorf("want result length != 0, got %d", len(cast.Cast))
	}

	for _, v := range cast.Cast {
		fmt.Println(v)
	}
}

func Test_MovieRepo_GetAllMovies(t *testing.T) {
	r := &dto.GetAllMoviesDto{
		StudioId: 0,
		Page:     1,
		PageSize: 2,
	}
	movies, err := repo.GetAllMovies(context.Background(), r)
	if err != nil {
		t.Fatalf("did not want an error %v", err)
	}

	if movies == nil {
		t.Fatalf("got nil result")
	}

	if len(movies) == 0 {
		t.Errorf("want result length != 0, got %d", len(movies))
	}
	if len(movies) > r.GetLimit() {
		t.Errorf("want limit %d, got %d", r.GetLimit(), len(movies))
	}

	for _, v := range movies {
		fmt.Println(v)
		if v.Studio.Id != r.GetStudioId() && r.GetStudioId() != 0 {
			t.Errorf("got film with studio_id %d, want %d", v.Studio.Id, r.GetStudioId())
		}
	}
}

func Test_MovieRepo_InsertAll(t *testing.T) {
	movies := []*agg.Movie{
		{Name: "Movie 1", Year: 2018, Gross: 10000, Rating: entity.Rating{Id: 1}, Studio: entity.Studio{Id: 1}},
		{Name: "Movie 2", Year: 2018, Gross: 10000, Rating: entity.Rating{Id: 2}, Studio: entity.Studio{Id: 2}},
		{Name: "Movie 3", Year: 2018, Gross: 10000, Rating: entity.Rating{Id: 3}, Studio: entity.Studio{Id: 3}},
	}
	err := repo.InsertAll(context.Background(), movies)
	if err != nil {
		t.Fatalf("did not want an error %v", err)
	}

	m1, _ := repo.GetMovieByName(context.Background(), &dto.GetMovieDto{Name: movies[0].Name})
	m2, _ := repo.GetMovieByName(context.Background(), &dto.GetMovieDto{Name: movies[1].Name})
	m3, _ := repo.GetMovieByName(context.Background(), &dto.GetMovieDto{Name: movies[2].Name})
	movies = []*agg.Movie{m1, m2, m3}
	for i, v := range movies {
		if v == nil {
			t.Errorf("insert movie %d failed", i)
		}
	}
}

func Test_MovieRepo_UpdateMovie(t *testing.T) {
	r := &dto.GetMovieDto{Id: 1}
	m, err := repo.GetMovieById(context.Background(), r)
	if err != nil {
		t.Fatalf("did not want an error %v", err)
	}
	m.Name = "Bollywood Unlimited"
	m.Gross = 0
	err = repo.Update(context.Background(), m)
	if err != nil {
		t.Fatalf("did not want an error %v", err)
	}
	m, err = repo.GetMovieById(context.Background(), r)
	if err != nil {
		t.Fatalf("did not want an error %v", err)
	}
	if m.Name != "Bollywood Unlimited" {
		t.Errorf("wanted name = %q, got %q", "Bollywood Unlimited", m.Name)
	}
	if m.Gross != 0 {
		t.Errorf("wanted gross = %d, got %d", 0, m.Gross)
	}
}

func Test_MovieRepo_DeleteMovie(t *testing.T) {
	r := &dto.GetMovieDto{Id: 1}
	m, err := repo.GetMovieById(context.Background(), r)
	if err != nil {
		t.Fatalf("did not want an error %v", err)
	}
	err = repo.Delete(context.Background(), m)
	if err != nil {
		t.Errorf("did not want an error %v", err)
	}
	m, err = repo.GetMovieById(context.Background(), r)
	if err == nil {
		t.Errorf("wanted a error, got nil")
	}
}
