package entity

import "time"

type Person struct {
	Id        int
	Name      string
	BirthDate time.Time
}

type PersonRole struct {
	Id       int
	PersonId int
	MovieId  int
	Role     string
}

type Rating struct {
	Id   int
	Name string
}

type Studio struct {
	Id   int
	Name string
}

type Movie struct {
	Id       int
	Name     string
	Year     int
	RatingId int
	Gross    int
	StudioId int
}
