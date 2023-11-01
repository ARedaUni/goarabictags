package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/amrojjeh/arabic-tags/internal/models"
	"github.com/amrojjeh/arabic-tags/internal/speech"
	"github.com/google/uuid"
)

type templateData struct {
	Excerpt         models.Excerpt
	Type            string
	Form            any
	Error           string
	GrammaticalTags []string
	Host            string
}

func newTemplateData(r *http.Request) (templateData, error) {
	err := r.ParseForm()
	if err != nil {
		return templateData{}, err
	}
	return templateData{
		Error:           r.Form.Get("error"),
		GrammaticalTags: speech.GrammaticalTags,
		Host:            r.Host,
	}, nil
}

func JSONFunc(s any) (string, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}

	return string(b[:]), nil
}

func IdFunc(s uuid.UUID) string {
	return strings.ReplaceAll(s.String(), "-", "")
}

func (app *application) cacheTemplates() error {
	app.page = make(map[string]*template.Template)
	funcs := template.FuncMap{
		"json": JSONFunc,
		"id":   IdFunc,
	}

	names, err := filepath.Glob("./ui/html/pages/*")
	if err != nil {
		return err
	}

	for _, name := range names {
		baseName := filepath.Base(name)

		base := template.New(name).Funcs(funcs)

		base, err := base.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return err
		}

		partials, err := filepath.Glob("./ui/html/partials/*")
		if err != nil {
			return err
		}

		for _, name := range partials {
			base, err = base.ParseFiles(name)
			if err != nil {
				return err
			}
		}

		app.page[baseName], err = base.ParseFiles(name)
		if err != nil {
			return err
		}
		app.logger.Info("page cached", slog.String("name", baseName))
	}

	return nil
}

func (app *application) renderTemplate(w http.ResponseWriter, page string,
	code int, data templateData) {
	template, ok := app.page[page]
	if !ok {
		app.serverError(w, errors.New(
			fmt.Sprintf("Page %v does not exist", page)))
		return
	}

	buffer := bytes.Buffer{}
	err := template.ExecuteTemplate(&buffer, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Ignoring error as it's unlikely to occur
	w.WriteHeader(code)
	_, err = buffer.WriteTo(w)
}
