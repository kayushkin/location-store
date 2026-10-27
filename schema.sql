-- location-store schema. Create-only and idempotent: Open() runs this whole
-- file on every boot. A new column goes in an additive migration in store.go,
-- never in an edit to a CREATE TABLE below.
--
-- Nothing a person can cite is hard-deleted: DELETE stamps deleted_at, and
-- every update, import change, delete and restore of a place, visit or trip
-- first copies the whole row into revisions.

PRAGMA foreign_keys = ON;

-- One imported file, found again by its hash: importing the same file twice
-- changes nothing. finished_at is 0 while an import is still sending; running
-- the importer again on the same file resumes it.
CREATE TABLE IF NOT EXISTS imports (
    id              TEXT PRIMARY KEY,              -- import_000001
    seq             INTEGER NOT NULL UNIQUE,
    kind            TEXT NOT NULL,                 -- vocabulary.go ImportKinds
    sha256          TEXT NOT NULL UNIQUE,          -- of the file (or of the records, for a photo-store import)
    source_machine  TEXT NOT NULL DEFAULT '',      -- where the file was read
    source_path     TEXT NOT NULL DEFAULT '',      -- its path there
    imported_at     INTEGER NOT NULL,
    finished_at     INTEGER NOT NULL DEFAULT 0,
    counts          TEXT NOT NULL DEFAULT '{}'     -- ImportCounts as JSON, summed over every batch
);

-- A named spot.
CREATE TABLE IF NOT EXISTS places (
    id              TEXT PRIMARY KEY,              -- place_000001
    seq             INTEGER NOT NULL UNIQUE,
    name            TEXT NOT NULL DEFAULT '',
    latitude        REAL NOT NULL,
    longitude       REAL NOT NULL,
    radius_meters   REAL NOT NULL DEFAULT 0,       -- 0 = not given
    address         TEXT NOT NULL DEFAULT '',
    kind            TEXT NOT NULL,                 -- vocabulary.go PlaceKinds
    google_place_id TEXT,                          -- NULL when no source gave one; an import finds the place again by it
    description     TEXT NOT NULL DEFAULT '',
    described_by    TEXT NOT NULL DEFAULT '',
    tags            TEXT NOT NULL DEFAULT '[]',
    source          TEXT NOT NULL,                 -- vocabulary.go RecordSources
    import_id       TEXT REFERENCES imports(id),   -- the import that created it; NULL when made by hand
    created_at      INTEGER NOT NULL,
    updated_at      INTEGER NOT NULL,
    deleted_at      INTEGER NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_places_google ON places(google_place_id) WHERE google_place_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_places_kind ON places(deleted_at, kind);

-- Being somewhere from started_at to ended_at.
CREATE TABLE IF NOT EXISTS visits (
    id                   TEXT PRIMARY KEY,         -- visit_000001
    seq                  INTEGER NOT NULL UNIQUE,
    place_id             TEXT REFERENCES places(id),
    started_at           INTEGER NOT NULL,
    ended_at             INTEGER NOT NULL,
    latitude             REAL,
    longitude            REAL,
    probability          REAL,                     -- 0..1, how sure the source is of the place; NULL when it gave none
    source_semantic_type TEXT NOT NULL DEFAULT '', -- the source's own word for the place (HOME, TYPE_WORK…)
    description          TEXT NOT NULL DEFAULT '',
    described_by         TEXT NOT NULL DEFAULT '',
    tags                 TEXT NOT NULL DEFAULT '[]',
    source               TEXT NOT NULL,
    source_ref           TEXT NOT NULL DEFAULT '', -- what an import finds the visit again by; '' when made by hand
    import_id            TEXT REFERENCES imports(id),
    created_at           INTEGER NOT NULL,
    updated_at           INTEGER NOT NULL,
    deleted_at           INTEGER NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_visits_source_ref ON visits(source, source_ref) WHERE source_ref != '';
CREATE INDEX IF NOT EXISTS idx_visits_time ON visits(deleted_at, started_at, seq);
CREATE INDEX IF NOT EXISTS idx_visits_place ON visits(place_id, started_at);

-- A movement. parent_trip_id groups the drives of a road trip or a holiday
-- under one trip made by hand; an import never sets it.
CREATE TABLE IF NOT EXISTS trips (
    id               TEXT PRIMARY KEY,             -- trip_000001
    seq              INTEGER NOT NULL UNIQUE,
    parent_trip_id   TEXT REFERENCES trips(id),
    started_at       INTEGER NOT NULL,
    ended_at         INTEGER NOT NULL,
    start_latitude   REAL,
    start_longitude  REAL,
    end_latitude     REAL,
    end_longitude    REAL,
    start_place_id   TEXT REFERENCES places(id),
    end_place_id     TEXT REFERENCES places(id),
    mode             TEXT NOT NULL,                -- vocabulary.go TripModes
    source_mode      TEXT NOT NULL DEFAULT '',     -- the source's own word (IN_PASSENGER_VEHICLE)
    distance_meters  REAL,
    probability      REAL,                         -- 0..1, how sure the source is of the mode
    title            TEXT NOT NULL DEFAULT '',
    description      TEXT NOT NULL DEFAULT '',
    described_by     TEXT NOT NULL DEFAULT '',
    tags             TEXT NOT NULL DEFAULT '[]',
    source           TEXT NOT NULL,
    source_ref       TEXT NOT NULL DEFAULT '',
    import_id        TEXT REFERENCES imports(id),
    created_at       INTEGER NOT NULL,
    updated_at       INTEGER NOT NULL,
    deleted_at       INTEGER NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_trips_source_ref ON trips(source, source_ref) WHERE source_ref != '';
CREATE INDEX IF NOT EXISTS idx_trips_time ON trips(deleted_at, started_at, seq);
CREATE INDEX IF NOT EXISTS idx_trips_parent ON trips(parent_trip_id, started_at);

-- A trip's path, in order. An import replaces it whole when the trip changes.
CREATE TABLE IF NOT EXISTS trip_path_points (
    trip_id    TEXT NOT NULL REFERENCES trips(id),
    position   INTEGER NOT NULL,
    latitude   REAL NOT NULL,
    longitude  REAL NOT NULL,
    timestamp  INTEGER,                            -- NULL when the source gives the point no time
    PRIMARY KEY (trip_id, position)
) WITHOUT ROWID;

-- Raw location fixes; there may be millions. No public ids: a point is read by
-- time. The unique index both serves time ranges and makes a second copy of a
-- fix from the same source a no-op. A point with a source_ref (a photo id)
-- replaces any earlier point with the same source and source_ref.
CREATE TABLE IF NOT EXISTS points (
    timestamp        INTEGER NOT NULL,
    latitude         REAL NOT NULL,
    longitude        REAL NOT NULL,
    accuracy_meters  REAL,
    source           TEXT NOT NULL,                -- vocabulary.go PointSources
    source_ref       TEXT NOT NULL DEFAULT '',
    import_id        TEXT NOT NULL REFERENCES imports(id)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_points_identity ON points(timestamp, latitude, longitude, source);
CREATE INDEX IF NOT EXISTS idx_points_source_ref ON points(source, source_ref) WHERE source_ref != '';

CREATE TABLE IF NOT EXISTS revisions (
    revision_id  INTEGER PRIMARY KEY AUTOINCREMENT,
    entity_type  TEXT NOT NULL,                    -- place | visit | trip
    entity_id    TEXT NOT NULL,
    reason       TEXT NOT NULL,                    -- update | import | delete | restore | attach
    replaced_at  INTEGER NOT NULL,
    row_json     TEXT NOT NULL                     -- the whole row as it was, as JSON
);
CREATE INDEX IF NOT EXISTS idx_revisions_entity ON revisions(entity_id, revision_id DESC);

-- A stretch of time with one home base: the place the user slept most
-- nights. location-group-journeys writes these from the visits, and a
-- re-run updates the ones it wrote. ended_at 0 runs on past the newest data.
-- A level above journeys, when there is one, is likely to start from these.
CREATE TABLE IF NOT EXISTS home_periods (
    id              TEXT PRIMARY KEY,              -- home_period_000001
    seq             INTEGER NOT NULL UNIQUE,
    place_id        TEXT REFERENCES places(id),    -- NULL when the nights there named no place
    latitude        REAL NOT NULL,
    longitude       REAL NOT NULL,
    locality        TEXT,                          -- Locality as JSON; NULL until named
    started_at      INTEGER NOT NULL,
    ended_at        INTEGER NOT NULL DEFAULT 0,
    nights_known    INTEGER NOT NULL DEFAULT 0,
    nights_at_home  INTEGER NOT NULL DEFAULT 0,
    described_by    TEXT NOT NULL,                 -- who found it: location-group-journeys
    created_at      INTEGER NOT NULL,
    updated_at      INTEGER NOT NULL,
    deleted_at      INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_home_periods_time ON home_periods(deleted_at, started_at);

-- Which visits and trips each import's file holds, whether the import made
-- them, changed them or found them as they were. A full-snapshot export
-- supersedes what an older one held in its span and it does not.
CREATE TABLE IF NOT EXISTS import_records (
    import_id    TEXT NOT NULL REFERENCES imports(id),
    entity_type  TEXT NOT NULL,                    -- visit | trip
    entity_id    TEXT NOT NULL,
    PRIMARY KEY (import_id, entity_type, entity_id)
) WITHOUT ROWID;
CREATE INDEX IF NOT EXISTS idx_import_records_entity ON import_records(entity_id);
