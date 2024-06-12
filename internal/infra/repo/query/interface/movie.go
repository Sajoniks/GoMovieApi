package queryinterface

type GetMovieCastById interface {
	GetMovieId() int
}

type GetMovieById interface {
	GetMovieId() int
}

type GetMovieByName interface {
	GetMovieName() string
}

type GetAllMovies interface {
	GetStudioId() int
	GetPage() int
	GetLimit() int
}
