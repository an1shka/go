package main

type Movie struct {
	ID       int     `json:"id"`
	Title    string  `json:"title"`
	Year     int     `json:"year"`
	Rating   float64 `json:"rating"`
	Director string  `json:"director"`
}

type MovieInput struct {
	Title    string  `json:"title"`
	Year     int     `json:"year"`
	Rating   float64 `json:"rating"`
	Director string  `json:"director"`
}
