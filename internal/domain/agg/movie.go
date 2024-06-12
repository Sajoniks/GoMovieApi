package agg

import "sajoniks.github.io/movieApi/internal/domain/entity"

type Movie struct {
	Id     int
	Name   string
	Year   int
	Gross  int
	Rating entity.Rating
	Studio entity.Studio
}
