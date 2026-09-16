package main

import (
	"errors"
	"strings"
)

// validateMovie проверяет корректность всех полей MovieInput
func validateMovie(input MovieInput) error {
	if strings.TrimSpace(input.Title) == "" {
		return errors.New("title cannot be empty")
	}

	if strings.TrimSpace(input.Director) == "" {
		return errors.New("director cannot be empty")
	}

	if input.Year <= 1888 {
		return errors.New("year must be greater than 1888")
	}

	if input.Year > 2100 {
		return errors.New("year must not be greater than 2100")
	}

	if input.Rating < 0.0 || input.Rating > 10.0 {
		return errors.New("rating must be between 0 and 10")
	}

	return nil
}