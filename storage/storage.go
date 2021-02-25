package storage

import (
	"database/sql"

	"github.com/babolivier/london-covid-vaccination/common"
	"github.com/babolivier/london-covid-vaccination/config"

	_ "github.com/lib/pq"
)

// Database interfaces with the PostgreSQL database.
type Database struct {
	db    *sql.DB
	stats *statsStatements
}

// NewDatabase instantiates a new Database.
func NewDatabase(cfg *config.DatabaseConfig) (d *Database, err error) {
	d = new(Database)

	if d.db, err = sql.Open("postgres", cfg.ConnString()); err != nil {
		return nil, err
	}

	if d.stats, err = newStatsStatements(d.db); err != nil {
		return nil, err
	}

	return d, nil
}

// GetKnownFileNames retrieves the names of the daily data files that already exist in
// the database.
func (d *Database) GetKnownFileNames() (map[string]bool, error) {
	rawFileNames, err := d.stats.selectFileNames()
	if err != nil {
		return nil, err
	}

	// Convert the slice of strings into a map indexed on these strings for quicker
	// lookup.
	fileNames := make(map[string]bool)
	for _, fileName := range rawFileNames {
		fileNames[fileName] = true
	}

	return fileNames, nil
}

// GetAllStats retrieves all the stats stored in the database.
func (d *Database) GetAllStats() ([]*common.DailyStats, error) {
	return d.stats.selectAllStats()
}

// StoreStats stores the provided stats into the database.
func (d *Database) StoreStats(stats *common.DailyStats) error {
	return d.stats.insertStats(stats)
}
