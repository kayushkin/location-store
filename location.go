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
	ImportID *string `json:"import_id"`
	// Locality is the town the place is in; null until
	// location-name-localities names it.
	Locality  *Locality `json:"locality"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
	DeletedAt int64     `json:"deleted_at"`
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
	// SupersededByImportID is the newer export that no longer holds this
	// visit, when that is why it is deleted.
	SupersededByImportID *string `json:"superseded_by_import_id"`
	// Locality is the town the visit was in: the visit's own coordinate, or
	// its place's. Null until location-name-localities names it.
	Locality  *Locality `json:"locality"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
	DeletedAt int64     `json:"deleted_at"`
}

// Trip is a movement from one place to another, or, with children, a road
// trip or holiday made of them.
type Trip struct {
	// ID is trip_000001.
	ID string `json:"id"`
	// Kind is one of GET /vocabularies trip_kinds: a movement a source
	// recorded, or a journey made of them.
	Kind string `json:"kind"`
	// ParentTripID is the trip this one is part of; never set by an import.
	ParentTripID *string `json:"parent_trip_id"`
	// ParentSetBy says who set or cleared ParentTripID: "" nobody yet,
	// "hand" a person or agent through the routes, or the command that
	// grouped it (location-group-journeys). The command moves only trips it
	// placed itself or that nobody has placed.
	ParentSetBy    string   `json:"parent_set_by"`
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
	// TitleBy says who wrote the title: "hand" through the routes, or the
	// command that generated it. A command rewrites only its own titles.
	TitleBy     string   `json:"title_by"`
	Description string   `json:"description"`
	DescribedBy string   `json:"described_by"`
	Tags        []string `json:"tags"`
	Source      string   `json:"source"`
	SourceRef   string   `json:"source_ref"`
	ImportID    *string  `json:"import_id"`
	// SupersededByImportID is the newer export that no longer holds this
	// trip, when that is why it is deleted.
	SupersededByImportID *string `json:"superseded_by_import_id"`
	// StartLocality and EndLocality are the towns the trip started and ended
	// in; null until location-name-localities names them. A journey's are
	// those of its first and last trip.
	StartLocality *Locality `json:"start_locality"`
	EndLocality   *Locality `json:"end_locality"`
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
	// Snapshot is the span a full-snapshot export covers, and what says how
	// new it is; null for any other import.
	Snapshot *ImportSnapshot `json:"snapshot"`
	// SnapshotChanges is what reconciling the snapshots did when this import
	// last finished; null when it did not run.
	SnapshotChanges *SnapshotChanges `json:"snapshot_changes"`
}

// ImportSnapshot describes an on-device Timeline export, which holds every
// visit and activity of the span it covers: an export newer than another
// supersedes the older one's records inside its own span.
type ImportSnapshot struct {
	// CoversFrom and CoversTo are the earliest start and the latest end of
	// the file's semantic segments.
	CoversFrom int64 `json:"covers_from"`
	CoversTo   int64 `json:"covers_to"`
	// RawSignalsOldestAt and RawSignalsNewestAt are the oldest and newest
	// raw position the file holds; 0 when it holds none. The phone keeps raw
	// positions in a rolling buffer, so a later export has a later newest
	// position, or the same newest and a later oldest.
	RawSignalsOldestAt int64 `json:"raw_signals_oldest_at"`
	RawSignalsNewestAt int64 `json:"raw_signals_newest_at"`
}

// SnapshotChanges is what reconciling the snapshot imports did: records an
// older export held that the newest export covering them does not, set
// aside, and records set aside before that a newer export holds again,
// brought back.
type SnapshotChanges struct {
	VisitsSuperseded int `json:"visits_superseded"`
	TripsSuperseded  int `json:"trips_superseded"`
	VisitsRestored   int `json:"visits_restored"`
	TripsRestored    int `json:"trips_restored"`
	// KeptEditedIDs are records the newest export no longer holds, kept
	// because someone edited them.
	KeptEditedIDs []string `json:"kept_edited_ids"`
}

// ImportFinish is the body of POST /imports/{id}/finish; it may be empty.
type ImportFinish struct {
	// Snapshot, for an on-device Timeline export, makes the import a full
	// snapshot of its span.
	Snapshot *ImportSnapshot `json:"snapshot"`
}

// ImportStart is the body of POST /imports.
type ImportStart struct {
	Kind          string `json:"kind"`
	Sha256        string `json:"sha256"`
	SourceMachine string `json:"source_machine"`
	SourcePath    string `json:"source_path"`
	// Reapply opens a finished import of the same file again, so the file is
	// sent once more: to record which records it holds, for an import made
	// before the store kept that.
	Reapply bool `json:"reapply"`
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
	// Kind is journey when left out.
	Kind           string   `json:"kind"`
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

// Locality is the town a coordinate is in, named offline from GeoNames'
// towns of 1000 people or more: of the towns whose reach covers the
// coordinate, the most populous, else the nearest.
type Locality struct {
	City string `json:"city"`
	// Region is the first-level division: a state, province or county.
	Region      string `json:"region"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	// TimeZone is the town's IANA time zone (Africa/Ndjamena).
	TimeZone   string `json:"time_zone"`
	GeonamesID int64  `json:"geonames_id"`
	// DistanceMeters is from the coordinate to the town's centre.
	DistanceMeters float64 `json:"distance_meters"`
	// DescribedBy says who named it: the command and the data it read.
	DescribedBy string `json:"described_by"`
}

// LocalityFillRequest is the body of POST /localities/fill.
type LocalityFillRequest struct {
	DryRun bool `json:"dry_run"`
}

// LocalityFill answers POST /localities/fill: how many empty localities it
// filled (or, in a dry run, would fill), and how many it could not name.
type LocalityFill struct {
	DryRun          bool `json:"dry_run"`
	Places          int  `json:"places"`
	Visits          int  `json:"visits"`
	TripStarts      int  `json:"trip_starts"`
	TripEnds        int  `json:"trip_ends"`
	AlreadyNamed    int  `json:"already_named"`
	NoTownNearby    int  `json:"no_town_nearby"`
	NoCoordinate    int  `json:"no_coordinate"`
	DistinctTowns   int  `json:"distinct_towns"`
	GeonamesTownsIn int  `json:"geonames_towns_in"`
}

// HomePeriod is a stretch of time with one home base: the place the user
// slept most nights. EndedAt 0 means it runs on past the newest data.
type HomePeriod struct {
	// ID is home_period_000001.
	ID string `json:"id"`
	// PlaceID is the home place; null when the nights there named no place.
	PlaceID   *string   `json:"place_id"`
	PlaceName string    `json:"place_name"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Locality  *Locality `json:"locality"`
	StartedAt int64     `json:"started_at"`
	EndedAt   int64     `json:"ended_at"`
	// NightsKnown is how many nights of the period the data places anywhere;
	// NightsAtHome how many of those were at the home.
	NightsKnown  int    `json:"nights_known"`
	NightsAtHome int    `json:"nights_at_home"`
	DescribedBy  string `json:"described_by"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
	DeletedAt    int64  `json:"deleted_at"`
}

// JourneyGroupRequest is the body of POST /journeys/group. A zero threshold
// takes its default.
type JourneyGroupRequest struct {
	DryRun bool `json:"dry_run"`
	// AwayDistanceMeters is how far from home counts as away (80 km).
	AwayDistanceMeters float64 `json:"away_distance_meters"`
	// MinimumAwaySeconds is how long a span away without a night away must
	// last to be a journey (6 h).
	MinimumAwaySeconds int64 `json:"minimum_away_seconds"`
	// HomeChangeDays is how many days a new night place must lead before it
	// becomes the home (21).
	HomeChangeDays int `json:"home_change_days"`
	// ReturnBreakSeconds is how long back near home, without a night at
	// home, ends a journey (12 h).
	ReturnBreakSeconds int64 `json:"return_break_seconds"`
}

// PlannedJourney is one journey a grouping found, and what it did or would do.
type PlannedJourney struct {
	// Action is create, update, unchanged or skip.
	Action string `json:"action"`
	// TripID is the journey's trip; "" for one a dry run would create.
	TripID    string `json:"trip_id"`
	Title     string `json:"title"`
	StartedAt int64  `json:"started_at"`
	EndedAt   int64  `json:"ended_at"`
	Mode      string `json:"mode"`
	// Nights is how many nights (03:00 local) fall in the journey.
	Nights int `json:"nights"`
	// FarthestMeters is the farthest the journey got from home.
	FarthestMeters float64  `json:"farthest_meters"`
	ChildTripIDs   []string `json:"child_trip_ids"`
	// ChildrenAttached counts trips this run made (or would make) part of it;
	// ChildrenDetached trips it took back out; ChildrenKeptElsewhere trips in
	// its span whose parent someone set by hand, left where they are.
	ChildrenAttached      int `json:"children_attached"`
	ChildrenDetached      int `json:"children_detached"`
	ChildrenKeptElsewhere int `json:"children_kept_elsewhere"`
	// Notes says why it was skipped, or what looks doubtful.
	Notes []string `json:"notes"`
}

// JourneyGrouping answers POST /journeys/group.
type JourneyGrouping struct {
	DryRun      bool                `json:"dry_run"`
	Request     JourneyGroupRequest `json:"request"`
	HomePeriods []HomePeriod        `json:"home_periods"`
	Journeys    []PlannedJourney    `json:"journeys"`
	// DeletedJourneyIDs are journeys an earlier run made that the data no
	// longer supports; soft-deleted, so they can be restored.
	DeletedJourneyIDs []string `json:"deleted_journey_ids"`
	// HomePlacesMarked are places whose kind this run set to home.
	HomePlacesMarked []string `json:"home_places_marked"`
	Notes            []string `json:"notes"`
}

// JourneyLocality is one town a journey stayed in, in order of arrival.
type JourneyLocality struct {
	Locality Locality `json:"locality"`
	// Seconds is the time the journey's visits spent there.
	Seconds   int64 `json:"seconds"`
	ArrivedAt int64 `json:"arrived_at"`
}

// JourneySummary answers GET /journeys/{id}: the journey's own trip, the
// trips that make it up, the visits in its span and the towns it stayed in.
type JourneySummary struct {
	Trip       Trip              `json:"trip"`
	Children   []Trip            `json:"children"`
	Visits     []Visit           `json:"visits"`
	Localities []JourneyLocality `json:"localities"`
}
