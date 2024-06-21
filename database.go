package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Database interface {
	ListProfessors() ([]*Professor, error)
	GetProfessorById(int) (*Professor, error)
	CreateProfessor(*Professor) error
	DeleteProfessor(int) error
	UpdateProfessor(int, *Professor) error
}

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB() (*PostgresDB, error) {
	// connStr := "user=postgres dbname=postgres password=lemm2301 sslmode=disable"
	connStr := "user= " + os.Getenv("DATABASE_USER") + " dbname=" + os.Getenv("DATABASE_NAME") + " password=" + os.Getenv("DATABASE_PASSWORD") + " sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresDB{
		db: db,
	}, nil
}

func (db *PostgresDB) createProfessorTable() error {
	query := `CREATE TABLE IF NOT EXISTS professor(
		id serial primary key,
		first_name varchar(255) not null,
		last_name varchar(255) not null,
		second_last_name varchar(255) not null,
		age int not null,
		degree varchar(255) not null,
		created_at timestamp,
		updated_at timestamp
	)`

	_, err := db.db.Exec(query)
	return err
}

func (db *PostgresDB) ListProfessors() ([]*Professor, error) {
	rows, err := db.db.Query("SELECT * FROM professor")
	if err != nil {
		return nil, err
	}

	professors := []*Professor{}
	for rows.Next() {
		professor, err := scanToProfessor(rows)
		if err != nil {
			return nil, err
		}
		professors = append(professors, professor)
	}
	return professors, nil
}
func (db *PostgresDB) GetProfessorById(id int) (*Professor, error) {
	rows, err := db.db.Query("SELECT * FROM professor WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanToProfessor(rows)
	}
	return nil, fmt.Errorf("professor with id %d not found", id)
}
func (db *PostgresDB) CreateProfessor(professor *Professor) error {
	query := `INSERT INTO professor (first_name, last_name, second_last_name, age, degree, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)`
	resp, err := db.db.Query(
		query,
		professor.FirstName,
		professor.LastName,
		professor.SecondLastName,
		professor.Age,
		professor.Degree,
		professor.CreatedAt,
		professor.UpdatedAt)

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)
	return nil

}
func (db *PostgresDB) DeleteProfessor(id int) error {
	_, err := db.db.Query("DELETE FROM professor WHERE id=$1", id)
	return err
}
func (db *PostgresDB) UpdateProfessor(id int, professor *Professor) error {
	query := `UPDATE professor SET first_name=$1, last_name=$2, second_last_name=$3, age=$4, degree=$5 WHERE id=$6`
	resp, err := db.db.Query(
		query,
		professor.FirstName,
		professor.LastName,
		professor.SecondLastName,
		professor.Age,
		professor.Degree,
		id)

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)
	return nil
}
func (db *PostgresDB) Init() error {
	return db.createProfessorTable()
}

func scanToProfessor(rows *sql.Rows) (*Professor, error) {
	professor := new(Professor)
	err := rows.Scan(
		&professor.ID,
		&professor.FirstName,
		&professor.LastName,
		&professor.SecondLastName,
		&professor.Age,
		&professor.Degree,
		&professor.CreatedAt,
		&professor.UpdatedAt)

	return professor, err
}
