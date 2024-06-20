package main

import (
	"time"
)

type Professor struct {
	ID             int       `json:"id"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	SecondLastName string    `json:"secondLastName"`
	Age            int       `json:"age"`
	Degree         string    `json:"degree"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type CreateProfessorRequest struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	SecondLastName string `json:"secondLastName"`
	Age            int    `json:"age"`
	Degree         string `json:"degree"`
}

func NewProfessor(firstName, lastName, secondLastName, degree string, age int) *Professor {
	return &Professor{
		FirstName:      firstName,
		LastName:       lastName,
		SecondLastName: secondLastName,
		Age:            age,
		Degree:         degree,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}
}
