package models

import (
	"database/sql"
	"time"
)

type HabitLog struct {
  ID        int
  UserID    int
  HabitID   int
  Date      time.Time
  IsCompleted bool
}

type HabitLogModel struct {
  DB *sql.DB
}

func (m *HabitLogModel) Insert(habitID int, userID int, date time.Time, isCompleted bool) error {
  stmt := `INSERT INTO habit_logs (habit_id, user_id, date, is_completed) VALUES (?, ?, ?, ?)`
  _, err := m.DB.Exec(stmt, habitID, userID, date, isCompleted)
  return err
}

func (m *HabitLogModel) Get(habitID, userID int, date time.Time) ([]*HabitLog, error) {
  stmt := `SELECT id, habit_id, user_id, date, is_completed FROM habit_logs WHERE habit_id = ? AND user_id = ? AND date = ?`
  rows, err := m.DB.Query(stmt, habitID, userID, date)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  logs := []*HabitLog{}
  for rows.Next() {
    log := &HabitLog{}
    err := rows.Scan(&log.ID, &log.HabitID, &log.UserID, &log.Date, &log.IsCompleted)
    if err != nil {
      return nil, err
    }
    logs = append(logs, log)
  }
  return logs, nil
}

func (m *HabitLogModel) List(userID int) ([]*HabitLog, error) {
  stmt := `SELECT id, habit_id, date, is_completed FROM habit_logs WHERE user_id = ?`

  rows, err := m.DB.Query(stmt, userID)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  habitsLog := []*HabitLog{}

  for rows.Next() {
    hl := &HabitLog{}

    err := rows.Scan(&hl.ID, &hl.HabitID, &hl.Date, &hl.IsCompleted)
    if err != nil {
      return nil, err
    }
    habitsLog = append(habitsLog, hl)
  }

  if err := rows.Err(); err != nil {
    return nil, err
  }
    
  return habitsLog, nil
}

// If log exists toggles the log is_completed to true of false without losing the logged in data
func (m *HabitLogModel) Update(habitID, userID int, date time.Time, isCompleted bool) error {
  stmt := `UPDATE habit_logs SET is_completed = ? WHERE habit_id = ? AND user_id = ? AND date = ?`
  _, err := m.DB.Exec(stmt, isCompleted, habitID, userID, date.Format("2006-01-02"))
  if err != nil {
    return err
  }

  return nil
}

// Checks if the log exists in the DB
func (m *HabitLogModel) HasLog(habitID, userID int, date time.Time) (bool, error) {
  stmt := `SELECT EXISTS (SELECT 1 FROM habit_logs WHERE habit_id = ? AND user_id = ? AND DATE(date) = ? LIMIT 1)`
  var exists bool
  err := m.DB.QueryRow(stmt, habitID, userID, date.Format("2006-01-02")).Scan(&exists)
  if err != nil {
    return false, err
  }
  return exists, nil
}

// Checks the logs is_completed status
func (m *HabitLogModel) Status(habitID, userID int, date time.Time) (bool, error) {
  stmt := `SELECT is_completed FROM habit_logs WHERE habit_id = ? AND user_id = ? AND date = ?`
  var is_completed bool
  err := m.DB.QueryRow(stmt, habitID, userID, date).Scan(&is_completed)
  if err != nil {
    return false, err
  }

  return is_completed, nil
}
