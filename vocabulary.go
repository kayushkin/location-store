package locationstore

import "strings"

// The words a record's fields may take. GET /vocabularies serves them all.

// What a place is to the user.
const (
	PlaceKindHome  = "home"
	PlaceKindWork  = "work"
	PlaceKindVenue = "venue" // a named business or landmark
	PlaceKindSaved = "saved" // an address the user saved, searched for or gave a name
	PlaceKindOther = "other" // not known; an import may later say more
)

// PlaceKinds is every kind a place can have.
var PlaceKinds = []string{PlaceKindHome, PlaceKindWork, PlaceKindVenue, PlaceKindSaved, PlaceKindOther}

// How a trip moved.
const (
	ModeDriving      = "driving"
	ModeMotorcycling = "motorcycling"
	ModeWalking      = "walking"
	ModeRunning      = "running"
	ModeHiking       = "hiking"
	ModeCycling      = "cycling"
	ModeBus          = "bus"
	ModeTrain        = "train"
	ModeTram         = "tram"
	ModeSubway       = "subway"
	ModeFerry        = "ferry"
	ModeBoating      = "boating"
	ModeFlying       = "flying"
	ModeSkiing       = "skiing"
	// ModeMixed is a trip made of others that moved in more than one way: a
	// holiday with flights and drives.
	ModeMixed   = "mixed"
	ModeUnknown = "unknown"
)

// TripModes is every mode a trip can have.
var TripModes = []string{ModeDriving, ModeMotorcycling, ModeWalking, ModeRunning, ModeHiking, ModeCycling, ModeBus,
	ModeTrain, ModeTram, ModeSubway, ModeFerry, ModeBoating, ModeFlying, ModeSkiing, ModeMixed, ModeUnknown}

// Where a place, visit or trip came from.
const (
	// SourceGoogleTimeline is the on-device Timeline export (Timeline.json,
	// 2024 on).
	SourceGoogleTimeline = "google_timeline"
	// SourceGoogleSemantic is the legacy Takeout Semantic Location History.
	SourceGoogleSemantic = "google_semantic"
	// SourceManual is made by hand, by a person or an agent.
	SourceManual = "manual"
)

// SourceGenerated is made by one of this repo's commands from other records:
// a journey grouped by location-group-journeys.
const SourceGenerated = "generated"

// RecordSources is every source a place, visit or trip can have.
var RecordSources = []string{SourceGoogleTimeline, SourceGoogleSemantic, SourceManual, SourceGenerated}

// What a trip is.
const (
	// TripKindMovement is one movement a source recorded: a drive, a flight.
	TripKindMovement = "movement"
	// TripKindJourney is a span away from home made of movements: a road trip
	// or a holiday. Made by hand, or by location-group-journeys.
	TripKindJourney = "journey"
)

// TripKinds is every kind a trip can have.
var TripKinds = []string{TripKindMovement, TripKindJourney}

// Who set a field the routes and the commands both write.
const (
	// SetByHand is a person or an agent, through the routes.
	SetByHand = "hand"
	// SetByJourneyGrouping is location-group-journeys.
	SetByJourneyGrouping = "location-group-journeys"
	// SetByLocalityNaming is location-name-localities.
	SetByLocalityNaming = "location-name-localities"
)

// Where a raw point came from.
const (
	PointSourceGoogleRecords      = "google_records"       // legacy Takeout Records.json
	PointSourceGoogleTimelineRaw  = "google_timeline_raw"  // Timeline.json rawSignals positions
	PointSourceGoogleTimelinePath = "google_timeline_path" // Timeline.json timelinePath points
	PointSourceGoogleSemanticPath = "google_semantic_path" // legacy Semantic Location History simplifiedRawPath points
	PointSourcePhoto              = "photo"                // a photo's GPS, from photo-store; source_ref is the photo id
)

// PointSources is every source a point can have.
var PointSources = []string{PointSourceGoogleRecords, PointSourceGoogleTimelineRaw, PointSourceGoogleTimelinePath,
	PointSourceGoogleSemanticPath, PointSourcePhoto}

// What kind of file an import read.
const (
	ImportKindGoogleTimeline = "google_timeline" // on-device Timeline.json
	ImportKindGoogleRecords  = "google_records"  // legacy Takeout Records.json
	ImportKindGoogleSemantic = "google_semantic" // one legacy Semantic Location History month file
	ImportKindPhotoStore     = "photo_store"     // photos with GPS, read from photo-store
)

// ImportKinds is every kind an import can have.
var ImportKinds = []string{ImportKindGoogleTimeline, ImportKindGoogleRecords, ImportKindGoogleSemantic, ImportKindPhotoStore}

// Vocabularies is the answer to GET /vocabularies.
type Vocabularies struct {
	PlaceKinds    []string `json:"place_kinds"`
	TripKinds     []string `json:"trip_kinds"`
	TripModes     []string `json:"trip_modes"`
	RecordSources []string `json:"record_sources"`
	PointSources  []string `json:"point_sources"`
	ImportKinds   []string `json:"import_kinds"`
}

func knownWord(words []string, word string) bool {
	for _, known := range words {
		if word == known {
			return true
		}
	}
	return false
}

// googleActivityModes maps Google's activity types, in both the on-device
// export and legacy Takeout, to a trip mode. A type not here is unknown, and
// the trip keeps Google's word in source_mode.
var googleActivityModes = map[string]string{
	"IN_PASSENGER_VEHICLE": ModeDriving,
	"IN_VEHICLE":           ModeDriving,
	"IN_CAR":               ModeDriving,
	"IN_TAXI":              ModeDriving,
	"DRIVING":              ModeDriving,
	"MOTORCYCLING":         ModeMotorcycling,
	"WALKING":              ModeWalking,
	"ON_FOOT":              ModeWalking,
	"RUNNING":              ModeRunning,
	"HIKING":               ModeHiking,
	"CYCLING":              ModeCycling,
	"ON_BICYCLE":           ModeCycling,
	"IN_BUS":               ModeBus,
	"IN_TRAIN":             ModeTrain,
	"IN_TRAM":              ModeTram,
	"IN_SUBWAY":            ModeSubway,
	"IN_FERRY":             ModeFerry,
	"SAILING":              ModeBoating,
	"BOATING":              ModeBoating,
	"KAYAKING":             ModeBoating,
	"ROWING":               ModeBoating,
	"FLYING":               ModeFlying,
	"SKIING":               ModeSkiing,
}

// ModeForGoogleActivityType is the trip mode for one of Google's activity
// types (IN_PASSENGER_VEHICLE), or unknown.
func ModeForGoogleActivityType(activityType string) string {
	if mode, found := googleActivityModes[strings.ToUpper(strings.TrimSpace(activityType))]; found {
		return mode
	}
	return ModeUnknown
}

// PlaceKindForGoogleSemanticType is the place kind for one of Google's
// semantic types: HOME and INFERRED_HOME in the on-device export, TYPE_HOME in
// legacy Takeout, and the like. named says the source gave the place a name,
// which makes a place of no stated type a venue.
func PlaceKindForGoogleSemanticType(semanticType string, named bool) string {
	switch strings.TrimPrefix(strings.ToUpper(strings.TrimSpace(semanticType)), "TYPE_") {
	case "HOME", "INFERRED_HOME":
		return PlaceKindHome
	case "WORK", "INFERRED_WORK":
		return PlaceKindWork
	case "SEARCHED_ADDRESS", "ALIASED_LOCATION":
		return PlaceKindSaved
	}
	if named {
		return PlaceKindVenue
	}
	return PlaceKindOther
}
