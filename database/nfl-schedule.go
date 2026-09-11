package database

import (
	"database/sql"
	// "fmt"
)

type ScheduleRepository struct {
	DB *sql.DB
}

type Week struct {
	Type      string // 1=Preseason, 2=Regular Season, 3=Postseason
	Value     string
	StartDate string
	EndDate   string
}

func (r *ScheduleRepository) CreateTable() error {
	_, err := r.DB.Exec(`CREATE TABLE IF NOT EXISTS weeks (
	 type TEXT PRIMARY KEY,
	 value TEXT,
	 startDate TEXT,
	 endDate TEXT
 )`)

	return err
}

func (r *ScheduleRepository) Insert(row Week) error {
	_, err := r.DB.Exec(`INSERT INTO weeks (
		type,
		value,
		startDate,
		endDate) VALUES (?, ?, ?, ?)`,
		row.Type, row.Value, row.StartDate, row.EndDate)

	return err
}

func (r *ScheduleRepository) IsSeeded() bool {
	var seeded bool
	r.DB.QueryRow("SELECT 1 FROM weeks").Scan(&seeded)

	return seeded
}

func (r *ScheduleRepository) FindWeekFromTime(time string) (Week, error) {
	var week Week
	err := r.DB.QueryRow("SELECT * FROM weeks WHERE startDate < ? AND endDate > ?", time, time).Scan(
		&week.Type, &week.Value, &week.StartDate, &week.EndDate)

	if err != nil {
		return Week{}, err
	}

	return week, nil
}
