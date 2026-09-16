package main

import (
	"errors"
	"sync"
)

var (
	ErrNotFound = errors.New("movie not found")
)

type MemoryStorage struct {
	mu     sync.RWMutex
	movies map[int]Movie
	nextID int
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		movies: make(map[int]Movie),
		nextID: 1,
	}
}

func (s *MemoryStorage) GetAll() []Movie {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Movie, 0, len(s.movies))
	for _, m := range s.movies {
		result = append(result, m)
	}
	return result
}

func (s *MemoryStorage) GetByID(id int) (Movie, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	movie, ok := s.movies[id]
	if !ok {
		return Movie{}, ErrNotFound
	}
	return movie, nil
}

func (s *MemoryStorage) Create(input MovieInput) Movie {
	s.mu.Lock()
	defer s.mu.Unlock()

	movie := Movie{
		ID:       s.nextID,
		Title:    input.Title,
		Year:     input.Year,
		Rating:   input.Rating,
		Director: input.Director,
	}
	s.movies[s.nextID] = movie
	s.nextID++

	return movie
}

func (s *MemoryStorage) Update(id int, input MovieInput) (Movie, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.movies[id]; !ok {
		return Movie{}, ErrNotFound
	}

	updated := Movie{
		ID:       id,
		Title:    input.Title,
		Year:     input.Year,
		Rating:   input.Rating,
		Director: input.Director,
	}
	s.movies[id] = updated

	return updated, nil
}

func (s *MemoryStorage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.movies[id]; !ok {
		return ErrNotFound
	}

	delete(s.movies, id)
	return nil
}