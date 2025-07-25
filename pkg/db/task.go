package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (?, ?, ?, ?)
	`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func Tasks(limit int) ([]*Task, error) {
	rows, err := DB.Query(`
		SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var (
			id      int64
			date    string
			title   string
			comment string
			repeat  string
		)
		if err := rows.Scan(&id, &date, &title, &comment, &repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &Task{
			ID:      strconv.FormatInt(id, 10),
			Date:    date,
			Title:   title,
			Comment: comment,
			Repeat:  repeat,
		})
	}
	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, nil
}

func TasksSearch(query string, limit int) ([]*Task, error) {
	pattern := "%" + query + "%"
	rows, err := DB.Query(`
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE title LIKE ? OR comment LIKE ? OR date LIKE ?
		ORDER BY date
		LIMIT ?
	`, pattern, pattern, pattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var (
			id      int64
			date    string
			title   string
			comment string
			repeat  string
		)
		if err := rows.Scan(&id, &date, &title, &comment, &repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &Task{
			ID:      strconv.FormatInt(id, 10),
			Date:    date,
			Title:   title,
			Comment: comment,
			Repeat:  repeat,
		})
	}
	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, nil
}

func UpdateTask(task *Task) error {
	query := `
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?
	`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("no task found with given ID")
	}

	return nil
}

func GetTask(id string) (*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

	var task Task
	err := DB.QueryRow(query, id).Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task not found")
	}
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func DeleteTask(id string) error {
	_, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
	return err
}

func UpdateDate(next string, id string) error {
	_, err := DB.Exec("UPDATE scheduler SET date = ? WHERE id = ?", next, id)
	return err
}

func PutNextDate(current string, repeat string) (string, error) {
	t, err := time.Parse("20060102", current)
	if err != nil {
		return "", fmt.Errorf("failed to parse date: %v", err)
	}

	parts := strings.Fields(repeat)
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid repeat format: %s", repeat)
	}

	switch parts[0] {
	case "d", "w":
		n, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("failed to parse number: %v", err)
		}
		days := n
		if parts[0] == "w" {
			days = n * 7
		}
		t = t.AddDate(0, 0, days)
		return t.Format("20060102"), nil

	case "m":
		if len(parts) < 3 {
			return "", fmt.Errorf("format 'm' must contain days and months")
		}

		dayStrs := strings.Split(parts[1], ",")
		monthStrs := strings.Split(parts[2], ",")

		var days []int
		for _, d := range dayStrs {
			day, err := strconv.Atoi(strings.TrimSpace(d))
			if err != nil {
				return "", fmt.Errorf("invalid day: %v", d)
			}
			days = append(days, day)
		}

		var months []int
		for _, m := range monthStrs {
			month, err := strconv.Atoi(strings.TrimSpace(m))
			if err != nil {
				return "", fmt.Errorf("invalid month: %v", m)
			}
			months = append(months, month)
		}

		for i := 0; i < 36; i++ {
			t = t.AddDate(0, 1, 0)
			month := int(t.Month())

			if !contains(months, month) {
				continue
			}

			daysInMonth := daysIn(t.Year(), t.Month())
			for _, d := range days {
				day := d
				if d < 0 {
					day = daysInMonth + d + 1
				}
				if day <= 0 || day > daysInMonth {
					continue
				}

				candidate := time.Date(t.Year(), t.Month(), day, 0, 0, 0, 0, time.UTC)
				if candidate.After(t) {
					return candidate.Format("20060102"), nil
				}
			}
		}

		return "", fmt.Errorf("could not find next date")

	default:
		return "", fmt.Errorf("unsupported repeat type: %s", parts[0])
	}
}

func contains(arr []int, x int) bool {
	for _, a := range arr {
		if a == x {
			return true
		}
	}
	return false
}

func daysIn(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
