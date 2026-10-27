package locationstore

// Wire types: what the routes take and answer. ./generate-ts.sh renders this
// file and vocabulary.go to ts/model.ts. Timestamps are epoch seconds; 0
// means unset. A nullable field is null when no source said.

// Place is a named spot.
type Place struct {
	// ID is place_000001: what anything else joins on.
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	// RadiusMeters is how far the place reaches; 0 when not given.
	RadiusMeters float64 `json:"radius_meters"`
	Address      string  `json:"address"`
	// Kind is one of GET /vocabularies place_kinds.
	Kind string `json:"kind"`
	// GooglePlaceID is Google's id for the place, when a source gave one; an
	// import finds the place again by it.
	GooglePlaceID *string `json:"google_place_id"`
	// Description is markdown; [text](ref:<id>) links name people and events.
	Description string   `json:"description"`
	DescribedBy string   `json:"described_by"`
	Tags        []string `json:"tags"`
	Source      string   `json:"source"`
	// ImportID is the import that created the place; null when made by hand.
	ImportID  *string `json:"import_id"`
	CreatedAt int64   `json:"created_at"`
	UpdatedAt int64   `json:"updated_at"`
	DeletedAt int64   `json:"deleted_at"`
}

// Visit is being somewhere from StartedAt to EndedAt.
type Visit struct {
	// ID is visit_000001.
	ID      string  `json:"id"`
	PlaceID *string `json:"place_id"`
	// PlaceName is the place's name when the visit was read, for display only.
	PlaceName string   `json:"place_name"`
	StartedAt int64    `json:"started_at"`
	EndedAt   int64    `json:"ended_at"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	// Probability is 0..1: how sure the source is of the place.
	Probability *float64 `json:"probability"`
	// SourceSemanticType is the source's own word for the place (HOME).
	SourceSemanticType string   `json:"source_semantic_type"`
	Description        string   `json:"description"`
	DescribedBy        string   `json:"described_by"`
	Tags               []string `json:"tags"`
	Source             string   `json:"source"`
	// SourceRef is what an import finds the visit again by; "" when made by
	// hand.
	SourceRef string  `json:"source_ref"`
	ImportID  *string `json:"import_id"`
	CreatedAt int64   `json:"created_at"`
	UpdatedAt int64   `json:"updated_at"`
	DeletedAt int64   `json:"deleted_at"`
}

// Trip is a movement from one place to another, or, with children, a road
// trip or holiday made of them.
type Trip struct {
	// ID is trip_000001.
	ID string `json:"id"`
	// ParentTripID is the trip this one is part of; set by hand, never by an
	// import.
	ParentTripID   *string  `json:"parent_trip_id"`
	StartedAt      int64    `json:"started_at"`
	EndedAt        int64    `json:"ended_at"`
	StartLatitude  *float64 `json:"start_latitude"`
	StartLongitude *float64 `json:"start_longitude"`
	EndLatitude    *float64 `json:"end_latitude"`
	EndLongitude   *float64 `json:"end_longitude"`
	StartPlaceID   *string  `json:"start_place_id"`
	// StartPlaceName and EndPlaceName are the places' names when the trip was
	// read, for display only.
	StartPlaceName string  `json:"start_place_name"`
	EndPlaceID     *string `json:"end_place_id"`
	EndPlaceName   string  `json:"end_place_name"`
	// Mode is one of GET /vocabularies trip_modes.
	Mode string `json:"mode"`
	// SourceMode is the source's own word for the mode (IN_PASSENGER_VEHICLE).
	SourceMode     string   `json:"source_mode"`
	DistanceMeters *float64 `json:"distance_meters"`
	// Probability is 0..1: how sure the source is of the mode.
	Probability *float64 `json:"probability"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	DescribedBy string   `json:"described_by"`
	Tags        []string `json:"tags"`
	Source      string   `json:"source"`
	SourceRef   string   `json:"source_ref"`
	ImportID    *string  `json:"import_id"`
	// PathPoints is how many points GET /trips/{id}/path answers.
	PathPoints int `json:"path_points"`
	// ChildTrips is how many live trips name this one as their parent.
	ChildTrips int   `json:"child_trips"`
	CreatedAt  int64 `json:"created_at"`
	UpdatedAt  int64 `json:"updated_at"`
	DeletedAt  int64 `json:"deleted_at"`
}

// PathPoint is one point of a trip's path, in order.
type PathPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	// Timestamp is null when the source gave the point no time.
	Timestamp *int64 `json:"timestamp"`
}

// Point is one raw location fix.
type Point struct {
	Timestamp      int64    `json:"timestamp"`
	Latitude       float64  `json:"latitude"`
	Longitude      float64  `json:"longitude"`
	AccuracyMeters *float64 `json:"accuracy_meters"`
	// Source is one of GET /vocabularies point_sources.
	Source string `json:"source"`
	// SourceRef names the record the point came from (a photo id), or "".
	SourceRef string `json:"source_ref"`
	ImportID  string `json:"import_id"`
}

// TimelineItem is one visit or one trip, in GET /timeline.
type TimelineItem struct {
	// Kind is "visit" or "trip"; the other field is null.
	Kind  string `json:"kind"`
	Visit *Visit `json:"visit"`
	Trip  *Trip  `json:"trip"`
}

// ImportCounts says what an import, or one batch of it, did.
type ImportCounts struct {
	PointsAdded     int `json:"points_added"`
	PointsUnchanged int `json:"points_unchanged"`
	PlacesCreated   int `json:"places_created"`
	PlacesUpdated   int `json:"places_updated"`
	VisitsCreated   int `json:"visits_created"`
	VisitsUpdated   int `json:"visits_updated"`
	TripsCreated    int `json:"trips_created"`
	TripsUpdated    int `json:"trips_updated"`
	// RecordsUnchanged counts places, visits and trips already as the import
	// says, or deleted by someone, so left alone.
	RecordsUnchanged int `json:"records_unchanged"`
}

// Import is one imported file.
type Import struct {
	// ID is import_000001.
	ID string `json:"id"`
	// Kind is one of GET /vocabularies import_kinds.
	Kind          string `json:"kind"`
	Sha256        string `json:"sha256"`
	SourceMachine string `json:"source_machine"`
	SourcePath    string `json:"source_path"`
	ImportedAt    int64  `json:"imported_at"`
	// FinishedAt is 0 while the import is still sending.
	FinishedAt int64        `json:"finished_at"`
	Counts     ImportCounts `json:"counts"`
}

// ImportStart is the body of POST /imports.
type ImportStart struct {
	Kind          string `json:"kind"`
	Sha256        string `json:"sha256"`
	SourceMachine string `json:"source_machine"`
	SourcePath    string `json:"source_path"`
}

// ImportStarted answers POST /imports. AlreadyFinished means the same file was
// imported before and there is nothing to send.
type ImportStarted struct {
	Import          Import `json:"import"`
	AlreadyFinished bool   `json:"already_finished"`
}

// ImportPlace is a place in an import batch, found again by GooglePlaceID.
type ImportPlace struct {
	GooglePlaceID string  `json:"google_place_id"`
	Name          string  `json:"name"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Address       string  `json:"address"`
	Kind          string  `json:"kind"`
	Source        string  `json:"source"`
}

// ImportVisit is a visit in an import batch, found again by Source and
// SourceRef. GooglePlaceID names a place in this batch or an earlier one.
type ImportVisit struct {
	Source             string   `json:"source"`
	SourceRef          string   `json:"source_ref"`
	StartedAt          int64    `json:"started_at"`
	EndedAt            int64    `json:"ended_at"`
	Latitude           *float64 `json:"latitude"`
	Longitude          *float64 `json:"longitude"`
	GooglePlaceID      string   `json:"google_place_id"`
	Probability        *float64 `json:"probability"`
	SourceSemanticType string   `json:"source_semantic_type"`
}

// ImportTrip is a trip in an import batch, found again by Source and
// SourceRef.
type ImportTrip struct {
	Source         string   `json:"source"`
	SourceRef      string   `json:"source_ref"`
	StartedAt      int64    `json:"started_at"`
	EndedAt        int64    `json:"ended_at"`
	StartLatitude  *float64 `json:"start_latitude"`
	StartLongitude *float64 `json:"start_longitude"`
	EndLatitude    *float64 `json:"end_latitude"`
	EndLongitude   *float64 `json:"end_longitude"`
	// StartGooglePlaceID and EndGooglePlaceID link the trip's ends to places
	// already known by those ids; an unknown id links nothing.
	StartGooglePlaceID string      `json:"start_google_place_id"`
	EndGooglePlaceID   string      `json:"end_google_place_id"`
	Mode               string      `json:"mode"`
	SourceMode         string      `json:"source_mode"`
	DistanceMeters     *float64    `json:"distance_meters"`
	Probability        *float64    `json:"probability"`
	Path               []PathPoint `json:"path"`
}

// ImportPoint is a raw point in an import batch.
type ImportPoint struct {
	Timestamp      int64    `json:"timestamp"`
	Latitude       float64  `json:"latitude"`
	Longitude      float64  `json:"longitude"`
	AccuracyMeters *float64 `json:"accuracy_meters"`
	Source         string   `json:"source"`
	SourceRef      string   `json:"source_ref"`
}

// ImportBatch is the body of POST /imports/{id}/batches. Places are written
// first, so a visit may name a place in the same batch.
type ImportBatch struct {
	Places []ImportPlace `json:"places"`
	Visits []ImportVisit `json:"visits"`
	Trips  []ImportTrip  `json:"trips"`
	Points []ImportPoint `json:"points"`
}

// PlaceCreate is the body of POST /places.
type PlaceCreate struct {
	Name         string   `json:"name"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	RadiusMeters float64  `json:"radius_meters"`
	Address      string   `json:"address"`
	Kind         string   `json:"kind"`
	Description  string   `json:"description"`
	DescribedBy  string   `json:"described_by"`
	Tags         []string `json:"tags"`
}

// VisitCreate is the body of POST /visits.
type VisitCreate struct {
	PlaceID     *string  `json:"place_id"`
	StartedAt   int64    `json:"started_at"`
	EndedAt     int64    `json:"ended_at"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	Description string   `json:"description"`
	DescribedBy string   `json:"described_by"`
	Tags        []string `json:"tags"`
}

// TripCreate is the body of POST /trips: usually a road trip or holiday that
// imported drives are then attached to.
type TripCreate struct {
	ParentTripID   *string  `json:"parent_trip_id"`
	StartedAt      int64    `json:"started_at"`
	EndedAt        int64    `json:"ended_at"`
	StartLatitude  *float64 `json:"start_latitude"`
	StartLongitude *float64 `json:"start_longitude"`
	EndLatitude    *float64 `json:"end_latitude"`
	EndLongitude   *float64 `json:"end_longitude"`
	StartPlaceID   *string  `json:"start_place_id"`
	EndPlaceID     *string  `json:"end_place_id"`
	Mode           string   `json:"mode"`
	DistanceMeters *float64 `json:"distance_meters"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	DescribedBy    string   `json:"described_by"`
	Tags           []string `json:"tags"`
}

// AttachChildren is the body of POST /trips/{id}/children.
type AttachChildren struct {
	TripIDs []string `json:"trip_ids"`
}

// Revision is the whole state of a place, visit or trip just before one
// change: Row is that record as GET answered it then.
type Revision struct {
	RevisionID int64  `json:"revision_id"`
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	Reason     string `json:"reason"`
	ReplacedAt int64  `json:"replaced_at"`
	Row        any    `json:"row"`
}

// Counts is what GET /health counts.
type Counts struct {
	Places    int `json:"places"`
	Visits    int `json:"visits"`
	Trips     int `json:"trips"`
	Points    int `json:"points"`
	Imports   int `json:"imports"`
	Deleted   int `json:"deleted"`
	Revisions int `json:"revisions"`
}

// Health is the answer to GET /health.
type Health struct {
	Status string `json:"status"`
	Counts Counts `json:"counts"`
}
