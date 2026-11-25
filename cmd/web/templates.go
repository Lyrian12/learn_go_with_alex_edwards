package main

import "gojourney/internal/models"

type templateData struct {
	CurrentYear int
	Snippets    []*models.Snippet
	Snippet     *models.Snippet
	Form        any
	Flash       string
}
