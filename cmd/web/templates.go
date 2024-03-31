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
  HabitsLog       *models.HabitLog
  HabitsLogs      []*models.HabitLog
  Form            any
  Flash           string
  IsAuthenticated bool
  CSRFToken       string
}

func humanDate(t time.Time) string {
  return t.Format("02 Jan 2006 at 15:04")
}

func days() []int {
    // Get the current date
    now := time.Now()

    // Calculate the day of the week for today (0 is Sunday, 1 is Monday, etc.)
    todayWeekday := int(now.Weekday())

    // Calculate the number of days to the previous Sunday
    daysToPreviousSunday := (todayWeekday + 6) % 7

    // Calculate the date of the previous Sunday
    previousSunday := now.AddDate(0, 0, -daysToPreviousSunday)

    // Initialize the slice to store the days of the current week
    var weekDays []int
    for i := 0; i < 7; i++ {
        // Add the day of the week to the slice
        weekDays = append(weekDays, previousSunday.Day())

        // Move to the next day (Monday, Tuesday, ..., Sunday)
        previousSunday = previousSunday.AddDate(0, 0, 1)
    }

    return weekDays
}

type LogInfo struct {
  Exists      bool
  IsCompleted bool
}

func hasLog(HabitsLogs []*models.HabitLog, habitID, day int) LogInfo {
  for _, log := range HabitsLogs {
    if log.HabitID == habitID && log.Date.Day() == day {
      return LogInfo{Exists: true, IsCompleted: log.IsCompleted}
    }
  }
  return LogInfo{Exists: false}
}

func isFutureDate(day int) bool {
  // Get the current date
  currentDay := time.Now().Day()

  // Compare the given date with the current date
  return day > currentDay
}

var functions = template.FuncMap{
  "humanDate":    humanDate,
  "days":         days,
  "hasLog":       hasLog,
  "isFutureDate": isFutureDate,
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
