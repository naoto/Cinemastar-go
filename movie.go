package main

import "time"

type Movie struct {
	Name      string    `json:"name"`
	Directory bool      `json:"Directory"`
	Due       time.Time `json:"due"`
	Path      string    `json:"path"`
	Thumbnail string    `json:"thumbnail"`
}

type Movies []Movie
