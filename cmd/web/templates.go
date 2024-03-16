package main

import (
  "html/template"
  "path/filepath"
  "time"

  "habits.cheezecake.net/internal/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
type templateData struct {
  CurrentYear     int
  Habit           *models.Habit
  Habits          []*models.Habit
  Form            any
  Flash           string
  IsAuthenticated bool
  CSRFToken       string
}

func humanDate(t time.Time) string {
  return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
  "humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
  // Initialize a new map to act as the cache.
  cache := map[string]*template.Template{}

  pages, err := filepath.Glob("./ui/html/pages/*.html")
  if err != nil {
    return nil, err
  }

  for _, page := range pages {
    // Extract the file name (like 'home.tmpl') from the full filepath
    // and assign it to the name variable.
    name := filepath.Base(page)

    // Parse the base template file into a template set.
    ts, err := template.New(name).Funcs(functions).ParseFiles("ui/html/base.html")
    if err != nil{
      return nil, err
    }

    // Call ParseGlob() *on this template set* to add any partials.
    ts, err = ts.ParseGlob("./ui/html/partials/*.html")
    if err != nil {
      return nil, err
    }

    // Call ParseFiles() *on this template set* to add the page template.
    ts, err = ts.ParseFiles(page)
    if err != nil {
      return nil, err
    }

    cache[name] = ts
  }

  return cache, nil
}


