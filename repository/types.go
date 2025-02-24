// This file contains types that are used in the repository layer.
package repository

import "time"

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type Estate struct {
	ID        string
	Length    int
	Width     int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Tree struct {
	ID        string
	EstateID  string
	X         int
	Y         int
	Height    int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Stats struct {
	TotalTrees int
	MaxHeight  int
	MinHeight  int
	Median     int
}
