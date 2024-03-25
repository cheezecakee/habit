package models

import (
	"database/sql"
	"errors"
	"time"
)

type Habit struct {
  ID      int
  Title   string
  Created time.Time
  UserID  int
}

type HabitModel struct {
  DB *sql.DB
}

func (m *HabitModel) Insert(title string, userID int) (int, error) {
  stmt := `INSERT INTO habits (title, created, user_id) VALUES(?, UTC_TIMESTAMP(),  ?)`

  // Use the Exec() method on the embedded connection pool to execute the statement.
  // The first parameter is the SQL statement, followed by the
  // title, content and expiry values for the placeholder parameters. 
  // This method returns a sql.Result type, which contains some basic
  // information about what happened when the statement was executed.
  result, err := m.DB.Exec(stmt, title, userID)
  if err != nil {
    return 0, err
  }

  // Use the LastInsertId() method on the result to get the ID of our
  // newly inserted record in the snippets table.
  id, err := result.LastInsertId()
  if err != nil {
    return 0, err
  }

  // The ID returned has the type int64, so we convert it to an int type
  // before returning.
  return int(id), nil
}

func (m *HabitModel) Get(id, userID int) (*Habit, error) {
  stmt := `SELECT * FROM habits WHERE id = ? AND user_id = ?`

  row := m.DB.QueryRow(stmt, id, userID)

  h := &Habit{}

  err := row.Scan(&h.Title, &h.Created, &h.ID,  &h.UserID)
  if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
      return nil, ErrNoRecord
    } else {
      return nil, err
    }
  }

  return h, nil
}

func (m *HabitModel) List(userID int) ([]*Habit, error) {
  stmt := `SELECT id, title FROM habits WHERE user_id = ? ORDER BY title ASC`

  rows, err := m.DB.Query(stmt, userID)
  if err != nil {
    return nil, err
  }

  defer rows.Close()

  habits := []*Habit{}
  
  for rows.Next() {
    h := &Habit{}

    err = rows.Scan(&h.ID, &h.Title)
    if err != nil {
      return nil, err
    }

    habits = append(habits, h)
  }

  if err = rows.Err(); err != nil {
    return nil, err
  }

  return habits, nil
}
 
