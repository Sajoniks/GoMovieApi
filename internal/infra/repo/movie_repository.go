package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"sajoniks.github.io/movieApi/internal/domain/agg"
	"sajoniks.github.io/movieApi/internal/infra/repo/interface"
	"sajoniks.github.io/movieApi/internal/infra/repo/query/interface"
	"time"
)

type movieRepoPostgres struct {
	conn *pgx.Conn
}

// NewMovieRepo @todo replace connString with DI
func NewMovieRepo(ctx context.Context, connString string) (repointerface.MovieRepository, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, errors.Join(errors.New("failed to init repository"), err)
	}
	return &movieRepoPostgres{
		conn: conn,
	}, nil
}

func (m *movieRepoPostgres) InsertAll(ctx context.Context, movies []*agg.Movie) error {
	qctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	tx, err := m.conn.Begin(qctx)
	if err != nil {
		return errors.Join(errors.New("insert all failed"), err)
	}
	b := &pgx.Batch{}
	for _, v := range movies {
		b.Queue(
			`INSERT INTO 
    			film (name, year, rating_id, gross, studio_id) 
			VALUES($1, $2, $3, $4, $5)`,
			v.Name, v.Year, v.Rating.Id, v.Gross, v.Studio.Id,
		)
	}
	br := m.conn.SendBatch(qctx, b)

	for range movies {
		_, brErr := br.Exec()
		if brErr != nil {
			br.Close()
			tx.Rollback(qctx)
			return errors.Join(errors.New("insert failed"), brErr)
		}
	}

	br.Close()
	err = tx.Commit(qctx)
	return err
}

func (m *movieRepoPostgres) GetMovieByName(ctx context.Context, q queryinterface.GetMovieByName) (*agg.Movie, error) {
	movie := &agg.Movie{}
	qctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	row := m.conn.QueryRow(qctx,
		`SELECT 
    			f.id, f.name, f.year, f.gross,
    			r.id, r.name,
    			s.id, s.name
    		FROM film f
    		LEFT JOIN rating r ON f.rating_id = r.id
    		LEFT JOIN studio s ON f.studio_id = s.id
    		WHERE f.name = $1`, q.GetMovieName(),
	)
	err := row.Scan(
		&movie.Id, &movie.Name, &movie.Year, &movie.Gross,
		&movie.Rating.Id, &movie.Rating.Name,
		&movie.Studio.Id, &movie.Studio.Name,
	)
	if err != nil {
		return nil, errors.Join(errors.New("query failed"), err)
	}
	return movie, nil
}

func (m *movieRepoPostgres) GetMovieById(ctx context.Context, q queryinterface.GetMovieById) (*agg.Movie, error) {
	movie := &agg.Movie{}
	qctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	row := m.conn.QueryRow(qctx,
		`SELECT
				f.id, f.name, f.year, f.gross,
				r.id, r.name,
				s.id, s.name
			FROM film f
			LEFT JOIN rating r ON f.rating_id = r.id
			LEFT JOIN studio s ON f.studio_id = s.id
			WHERE f.id = $1`, q.GetMovieId(),
	)
	err := row.Scan(
		&movie.Id, &movie.Name, &movie.Year, &movie.Gross,
		&movie.Rating.Id, &movie.Rating.Name,
		&movie.Studio.Id, &movie.Studio.Name,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("movie %d was not found", q.GetMovieId())
		} else {
			return nil, errors.Join(errors.New("query failed"), err)
		}
	}
	return movie, nil
}

func (m *movieRepoPostgres) Update(ctx context.Context, movie *agg.Movie) error {
	qctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	_, err := m.conn.Exec(qctx,
		`UPDATE film SET
				name = $1,
				year = $2,
				gross = $3,
				rating_id = $4,
				studio_id = $5
			WHERE id = $6`,
		movie.Name,
		movie.Year,
		movie.Gross,
		movie.Rating.Id,
		movie.Studio.Id,
		movie.Id,
	)
	if err != nil {
		return errors.Join(errors.New("update failed"), err)
	}
	return nil
}

func (m *movieRepoPostgres) Delete(ctx context.Context, movie *agg.Movie) error {
	qctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	_, err := m.conn.Exec(qctx, `DELETE FROM film WHERE id = $1`, movie.Id)
	if err != nil {
		return errors.Join(errors.New("delete failed"), err)
	}
	return nil
}

func (m *movieRepoPostgres) GetAllMovies(ctx context.Context, q queryinterface.GetAllMovies) ([]*agg.Movie, error) {
	movies := make([]*agg.Movie, 0)
	qctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	rows, err := m.conn.Query(qctx,
		`SELECT 
    				f.id, f.name, f.year, f.rating_id, r.name, f.gross, f.studio_id, s.name 
				FROM film f
				LEFT JOIN rating r ON r.id = f.rating_id
				LEFT JOIN studio s ON s.id = f.studio_id
				WHERE f.id > $1 AND ($2 = 0 OR f.studio_id = $2)
				ORDER BY f.id ASC
				LIMIT $3`,
		q.GetPage(),
		q.GetStudioId(),
		q.GetLimit(),
	)
	if err != nil {
		return nil, errors.Join(errors.New("query failed"), err)
	}
	for rows.Next() {
		movie := &agg.Movie{}
		scanErr := rows.Scan(
			&movie.Id,
			&movie.Name,
			&movie.Year,
			&movie.Rating.Id,
			&movie.Rating.Name,
			&movie.Gross,
			&movie.Studio.Id,
			&movie.Studio.Name,
		)
		if scanErr != nil {
			return nil, errors.Join(errors.New("read failed"), scanErr)
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

func (m *movieRepoPostgres) GetMovieCast(ctx context.Context, q queryinterface.GetMovieCastById) (*agg.MovieCast, error) {
	cast := &agg.MovieCast{}
	cast.Cast = make([]*agg.PersonRole, 0)
	qctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	rows, err := m.conn.Query(qctx,
		`SELECT
				pr.id, pr.role, p.id, p.name, p.birth_date
			FROM person p
			LEFT JOIN person_role pr 
			ON p.id = pr.person_id
			WHERE pr.film_id = $1`,
		q.GetMovieId(),
	)
	if err != nil {
		return nil, errors.Join(errors.New("query failed"), err)
	}
	for rows.Next() {
		role := &agg.PersonRole{}
		scanErr := rows.Scan(
			&role.Id,
			&role.Role,
			&role.Person.Id,
			&role.Person.Name,
			&role.Person.BirthDate,
		)
		if scanErr != nil {
			return nil, errors.Join(errors.New("read failed"), scanErr)
		}
		cast.Cast = append(cast.Cast, role)
	}

	return cast, nil
}
