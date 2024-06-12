package agg

import "sajoniks.github.io/movieApi/internal/domain/entity"

type PersonRole struct {
	Id     int
	Role   string
	Person entity.Person
}
