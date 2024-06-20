package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Server struct {
	addr     string
	database Database
}

func NewServer(addr string, database Database) *Server {
	return &Server{
		addr:     addr,
		database: database,
	}
}

func (s *Server) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/professor", makeHTTPHandleFunc(s.handleProfessor))
	router.HandleFunc("/professor/{id}", makeHTTPHandleFunc(s.handleProfessorId))

	log.Println("CM-Professors Service Running on port ", s.addr)

	http.ListenAndServe(s.addr, router)
}

func (s *Server) handleProfessor(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleListProfessors(w, r)
	case "POST":
		return s.handleCreateProfessor(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *Server) handleProfessorId(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetProfessorById(w, r)
	case "PUT":
		return s.handleUpdateProfessor(w, r)
	case "DELETE":
		return s.handleDeleteProfessor(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *Server) handleListProfessors(w http.ResponseWriter, r *http.Request) error {
	professors, err := s.database.ListProfessors()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, professors)
}
func (s *Server) handleGetProfessorById(w http.ResponseWriter, r *http.Request) error {
	id, err := getProfessorId(r)
	if err != nil {
		return err
	}
	professor, err := s.database.GetProfessorById(id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, professor)
}
func (s *Server) handleCreateProfessor(w http.ResponseWriter, r *http.Request) error {
	createProfessorRequest := new(CreateProfessorRequest)
	if err := json.NewDecoder(r.Body).Decode(createProfessorRequest); err != nil {
		return err
	}

	professor := NewProfessor(
		createProfessorRequest.FirstName,
		createProfessorRequest.LastName,
		createProfessorRequest.SecondLastName,
		createProfessorRequest.Degree,
		createProfessorRequest.Age)
	if err := s.database.CreateProfessor(professor); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, professor)
}
func (s *Server) handleUpdateProfessor(w http.ResponseWriter, r *http.Request) error {
	id, err := getProfessorId(r)
	if err != nil {
		return err
	}
	createProfessorRequest := new(CreateProfessorRequest)
	if err := json.NewDecoder(r.Body).Decode(createProfessorRequest); err != nil {
		return err
	}
	professor := NewProfessor(
		createProfessorRequest.FirstName,
		createProfessorRequest.LastName,
		createProfessorRequest.SecondLastName,
		createProfessorRequest.Degree,
		createProfessorRequest.Age)
	if err := s.database.UpdateProfessor(id, professor); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, professor)
}
func (s *Server) handleDeleteProfessor(w http.ResponseWriter, r *http.Request) error {
	id, err := getProfessorId(r)
	if err != nil {
		return err
	}
	if err := s.database.DeleteProfessor(id); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"professor deleteed": id})
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type ApiError struct {
	Error string `json:"error"`
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// Hanlde error
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func getProfessorId(r *http.Request) (int, error) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return id, fmt.Errorf("invalid id given: %s", idParam)
	}
	return id, nil
}
