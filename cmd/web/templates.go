package main

import (
	"snippetbox.micypac.io/internal/models"
)

type templateData struct {
	Snippet *models.Snippet
	Snippets []*models.Snippet
}
