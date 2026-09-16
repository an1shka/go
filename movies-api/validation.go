package main

import (
	"errors"
	"strings"
	"time"
)

func (input *MovieInput) Validate() error {
	if strings.TrimSpace(input.Title) == "" {
		return errors.New("title is required")
	}

	currentYear := time.Now().Year()
	if input.Year < 1888 || input.Year > currentYear+5 {
		return errors.New("invalid year")
	}

	if input.Rating < 0.0 || input.Rating > 10.0 {
		return errors.New("rating must be between 0.0 and 10.0")
	}

	if strings.TrimSpace(input.Director) == "" {
		return errors.New("director is required")
	}

	return nil
}