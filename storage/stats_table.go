package storage

import (
	"database/sql"

	"github.com/babolivier/london-covid-vaccination/common"
)

const statsSchema = `
CREATE TABLE IF NOT EXISTS stats (
	file_name TEXT PRIMARY KEY,
	pub_date DATE,
	first_dose BIGINT NOT NULL,
	second_dose BIGINT NOT NULL
);
`

const selectFileNamesSQL = `
SELECT file_name FROM stats
`

const selectAllStatsSQL = `
SELECT file_name, pub_date, first_dose, second_dose from stats
`

const insertStatsSQL = `
INSERT INTO stats (file_name, pub_date, first_dose, second_dose)
VALUES ($1, $2, $3, $4)
`

type statsStatements struct {
	selectFileNameStmt *sql.Stmt
	selectAllStatsStmt *sql.Stmt
	insertStatsStmt    *sql.Stmt
}

func newStatsStatements(db *sql.DB) (s *statsStatements, err error) {
	s = new(statsStatements)

	if _, err = db.Exec(statsSchema); err != nil {
		return
	}

	if s.selectFileNameStmt, err = db.Prepare(selectFileNamesSQL); err != nil {
		return
	}

	if s.selectAllStatsStmt, err = db.Prepare(selectAllStatsSQL); err != nil {
		return
	}

	if s.insertStatsStmt, err = db.Prepare(insertStatsSQL); err != nil {
		return
	}

	return
}

func (s *statsStatements) selectFileNames() ([]string, error) {
	files := make([]string, 0)

	rows, err := s.selectFileNameStmt.Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var filename string
		if err = rows.Scan(&filename); err != nil {
			return nil, err
		}
		files = append(files, filename)
	}

	return files, nil
}

func (s *statsStatements) selectAllStats() ([]*common.DailyStats, error) {
	stats := make([]*common.DailyStats, 0)

	rows, err := s.selectAllStatsStmt.Query()
	if err != nil {
		return nil, err
	}

	var ds *common.DailyStats
	for rows.Next() {
		ds = new(common.DailyStats)
		if err = rows.Scan(
			&ds.FileName,
			&ds.PubDate,
			&ds.FirstDose,
			&ds.SecondDose,
		); err != nil {
			return nil, err
		}

		stats = append(stats, ds)
	}

	return stats, nil
}

func (s *statsStatements) insertStats(stats *common.DailyStats) error {
	_, err := s.insertStatsStmt.Exec(
		stats.FileName,
		stats.PubDate,
		stats.FirstDose,
		stats.SecondDose,
	)
	return err
}
