// Package locationstore keeps where the user has been: named places, visits
// to them, trips between them, the raw location points underneath, and the
// imports that brought them in.
//
// Places, visits and trips are cited elsewhere by their ids (place_000001,
// visit_000001, trip_000001), which kanban-store's registry claims, so a
// journal timeline or a note can link one as a chip.
package locationstore

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrNotFound = errors.New("not found")
	ErrInvalid  = errors.New("invalid")
	ErrConflict = errors.New("conflict")
)

// Store owns the database.
type Store struct {
	db      *sql.DB
	dataDir string
	// writeMutex serializes every write; the next id is MAX(seq)+1 read inside
	// the write's transaction.
	writeMutex sync.Mutex
}

// Open creates the data directory (0700) if needed, opens the database and
// applies the schema.
func Open(dataDir string) (*Store, error) {
	if dataDir == "" {
		return nil, errors.New("open store: no data directory given")
	}
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}
	if err := os.Chmod(dataDir, 0o700); err != nil {
		return nil, err
	}
	dbPath := filepath.Join(dataDir, "location-store.db")
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("set WAL: %w", err)
	}
	if _, err := db.Exec(schemaSQL); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	if err := os.Chmod(dbPath, 0o600); err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db, dataDir: dataDir}, nil
}

// Close releases the database handle.
func (s *Store) Close() error { return s.db.Close() }

// DataDir is the directory the database lives in.
func (s *Store) DataDir() string { return s.dataDir }

func now() int64 { return time.Now().Unix() }

type queryer interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

func (s *Store) inTransaction(work func(tx *sql.Tx) error) error {
	s.writeMutex.Lock()
	defer s.writeMutex.Unlock()
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := work(tx); err != nil {
		return err
	}
	return tx.Commit()
}

// nextID mints the next id of a table: prefix_000001.
func nextID(tx *sql.Tx, table, prefix string) (string, int64, error) {
	var maxSeq sql.NullInt64
	if err := tx.QueryRow(`SELECT MAX(seq) FROM ` + table).Scan(&maxSeq); err != nil {
		return "", 0, err
	}
	seq := maxSeq.Int64 + 1
	return fmt.Sprintf("%s_%06d", prefix, seq), seq, nil
}

func encodeStrings(values []string) string {
	if values == nil {
		values = []string{}
	}
	encoded, err := json.Marshal(values)
	if err != nil {
		panic(err)
	}
	return string(encoded)
}

func decodeStrings(encoded, what string) ([]string, error) {
	var values []string
	if err := json.Unmarshal([]byte(encoded), &values); err != nil {
		return nil, fmt.Errorf("%s has tags that are not a JSON array: %w", what, err)
	}
	return values, nil
}

func normalizeTags(tags []string) ([]string, error) {
	out := []string{}
	seen := map[string]bool{}
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			return nil, fmt.Errorf("%w: a tag is empty", ErrInvalid)
		}
		if !seen[tag] {
			seen[tag] = true
			out = append(out, tag)
		}
	}
	return out, nil
}

func nullableFloat(value *float64) any {
	if value == nil {
		return nil
	}
	return *value
}

func nullableString(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}

func floatPointer(value sql.NullFloat64) *float64 {
	if !value.Valid {
		return nil
	}
	return &value.Float64
}

func stringPointer(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func sameFloat(a, b *float64) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}

func sameString(a, b *string) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}

func validLatitudeLongitude(latitude, longitude float64) bool {
	return latitude >= -90 && latitude <= 90 && longitude >= -180 && longitude <= 180
}

// checkCoordinatePair refuses half a coordinate or one out of range.
func checkCoordinatePair(what string, latitude, longitude *float64) error {
	if (latitude == nil) != (longitude == nil) {
		return fmt.Errorf("%w: %s latitude and longitude come together", ErrInvalid, what)
	}
	if latitude != nil && !validLatitudeLongitude(*latitude, *longitude) {
		return fmt.Errorf("%w: %s latitude or longitude out of range", ErrInvalid, what)
	}
	return nil
}

func checkTimeRange(startedAt, endedAt int64) error {
	if startedAt <= 0 || endedAt < startedAt {
		return fmt.Errorf("%w: started_at must be set and ended_at not before it (epoch seconds), got %d to %d", ErrInvalid, startedAt, endedAt)
	}
	return nil
}

func checkProbability(probability *float64) error {
	if probability != nil && (*probability < 0 || *probability > 1) {
		return fmt.Errorf("%w: probability must be between 0 and 1, got %v", ErrInvalid, *probability)
	}
	return nil
}

// Entity types, as revisions and the routes name them.
const (
	EntityPlace = "place"
	EntityVisit = "visit"
	EntityTrip  = "trip"
)

// snapshot copies the whole current record into revisions before a change, in
// the change's own transaction.
func snapshot(tx *sql.Tx, entityType, id, reason string, at int64) error {
	var row any
	var err error
	switch entityType {
	case EntityPlace:
		row, err = getPlaceRow(tx, id)
	case EntityVisit:
		row, err = getVisitRow(tx, id)
	case EntityTrip:
		row, err = getTripRow(tx, id)
	default:
		return fmt.Errorf("snapshot of unknown entity type %q", entityType)
	}
	if err != nil {
		return err
	}
	encoded, err := json.Marshal(row)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO revisions (entity_type, entity_id, reason, replaced_at, row_json) VALUES (?,?,?,?,?)`,
		entityType, id, reason, at, string(encoded))
	return err
}

const (
	defaultLimit = 100
	maximumLimit = 1000
)

func clampLimit(limit, fallbackDefault, maximum int) int {
	if limit <= 0 {
		return fallbackDefault
	}
	return min(limit, maximum)
}

// Revisions is a place's, visit's or trip's history, newest first.
func (s *Store) Revisions(entityType, id string, limit, offset int) ([]Revision, error) {
	var err error
	switch entityType {
	case EntityPlace:
		_, err = getPlaceRow(s.db, id)
	case EntityVisit:
		_, err = getVisitRow(s.db, id)
	case EntityTrip:
		_, err = getTripRow(s.db, id)
	}
	if err != nil {
		return nil, err
	}
	rows, err := s.db.Query(`SELECT revision_id, entity_type, entity_id, reason, replaced_at, row_json FROM revisions
		WHERE entity_id = ? ORDER BY revision_id DESC LIMIT ? OFFSET ?`, id, clampLimit(limit, defaultLimit, maximumLimit), max(offset, 0))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	revisions := []Revision{}
	for rows.Next() {
		var revision Revision
		var row string
		if err := rows.Scan(&revision.RevisionID, &revision.EntityType, &revision.EntityID, &revision.Reason, &revision.ReplacedAt, &row); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(row), &revision.Row); err != nil {
			return nil, fmt.Errorf("revision %d of %s: %w", revision.RevisionID, id, err)
		}
		revisions = append(revisions, revision)
	}
	return revisions, rows.Err()
}

// Counts reports what GET /health shows.
func (s *Store) Counts() (Counts, error) {
	var counts Counts
	err := s.db.QueryRow(`SELECT
		(SELECT COUNT(*) FROM places WHERE deleted_at = 0),
		(SELECT COUNT(*) FROM visits WHERE deleted_at = 0),
		(SELECT COUNT(*) FROM trips WHERE deleted_at = 0),
		(SELECT COUNT(*) FROM points),
		(SELECT COUNT(*) FROM imports),
		(SELECT COUNT(*) FROM places WHERE deleted_at != 0) + (SELECT COUNT(*) FROM visits WHERE deleted_at != 0) + (SELECT COUNT(*) FROM trips WHERE deleted_at != 0),
		(SELECT COUNT(*) FROM revisions)`).
		Scan(&counts.Places, &counts.Visits, &counts.Trips, &counts.Points, &counts.Imports, &counts.Deleted, &counts.Revisions)
	return counts, err
}

// likePattern turns a word into a LIKE pattern that matches it anywhere.
func likePattern(word string) string {
	return "%" + strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(word) + "%"
}
