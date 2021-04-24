package pirsch

import (
	"database/sql"
	"time"
)

// Store is the database storage interface.
type Store interface {
	// SaveHits saves new hits.
	SaveHits([]Hit) error

	// Session returns the last session timestamp for given tenant, fingerprint, and maximum age.
	Session(sql.NullInt64, string, time.Time) (time.Time, error)

	// Run returns the results for given query.
	Run(*Query) ([]Stats, error)
}
