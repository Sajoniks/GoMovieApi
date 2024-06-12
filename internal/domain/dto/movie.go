package dto

type GetMovieCastByIdDto struct {
	Id int
}

func (g *GetMovieCastByIdDto) GetMovieId() int {
	return g.Id
}

type GetAllMoviesDto struct {
	StudioId int
	Page     int
	PageSize int
}

func (g *GetAllMoviesDto) GetStudioId() int {
	return g.StudioId
}

func (g *GetAllMoviesDto) GetPage() int {
	return g.Page
}

func (g *GetAllMoviesDto) GetLimit() int {
	return g.PageSize
}

type GetMovieDto struct {
	Id   int
	Name string
}

func (g *GetMovieDto) GetMovieName() string {
	return g.Name
}

func (g *GetMovieDto) GetMovieId() int {
	return g.Id
}
