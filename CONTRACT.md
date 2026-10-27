# location-store — API contract

`127.0.0.1:8333`, loopback only, no auth: a front door such as dash adds it.
Where the user has been: named places, visits to them, trips between them, the
raw location points underneath, and the imports that brought them in. SQLite
at `<data dir>/location-store.db` (default `~/.config/location-store`,
directory 0700, file 0600). Timestamps are epoch seconds.

## Settings

| Variable | Default | What it is |
|---|---|---|
| `LOCATION_STORE_ADDR` | `127.0.0.1:8333` | listen address |
| `LOCATION_STORE_DATA_DIR` | `~/.config/location-store` | holds `location-store.db` |
| `LOCATION_STORE_URL` | `http://127.0.0.1:8333` | the import commands: where to post |
| `PHOTO_STORE_URL` | `http://127.0.0.1:8329` | `location-import-photos`: where photo-store answers |

## Records

Places, visits and trips have ids from this store (`place_000001`,
`visit_000001`, `trip_000001`), registered in kanban-store's entity-type
registry as types `place`, `visit` and `trip`, so a timeline item or a note
cites one as a chip. Each update, import change, delete and restore first
copies the whole record to `revisions`; deletes are soft.

**An import never touches** a name, title, description, tags or
`parent_trip_id`, never changes a place link already set, never changes a
trip's `mode` once it no longer matches what the source said (someone
corrected it), and leaves a deleted record deleted.

### `places`

| Field | Meaning |
|---|---|
| `id` | `place_000001` |
| `name`, `address` | an import fills them only while blank |
| `latitude`, `longitude`, `radius_meters` | `radius_meters` 0 = not given |
| `kind` | `GET /vocabularies` `place_kinds`: `home`, `work`, `venue`, `saved`, `other`. An import turns `other` into what the source says, never anything else |
| `google_place_id` | Google's id, when a source gave one; unique; an import finds the place again by it |
| `description`, `described_by` | markdown; `described_by` (who wrote it) is required with a description |
| `tags` | JSON array of strings |
| `source`, `import_id` | `google_timeline`, `google_semantic` or `manual`; the import that created it |
| `deleted_at` | 0 = live |

### `visits`

| Field | Meaning |
|---|---|
| `id` | `visit_000001` |
| `place_id`, `place_name` | the place, nullable; `place_name` is read with the visit, for display |
| `started_at`, `ended_at` | |
| `latitude`, `longitude` | where the source put the visit, nullable |
| `probability` | 0..1, how sure the source is of the place, nullable |
| `source_semantic_type` | the source's own word (`HOME`, `TYPE_WORK`, `UNKNOWN`) |
| `description`, `described_by`, `tags` | as on places |
| `source`, `source_ref`, `import_id` | an import finds the visit again by `source` and `source_ref` (`<started_at>-<ended_at>`) |

### `trips`

| Field | Meaning |
|---|---|
| `id` | `trip_000001` |
| `parent_trip_id` | the road trip or holiday this trip is part of; set by hand, never by an import |
| `started_at`, `ended_at` | |
| `start_latitude` … `end_longitude` | nullable, in pairs |
| `start_place_id`, `end_place_id` | nullable; `start_place_name`, `end_place_name` for display |
| `mode` | `GET /vocabularies` `trip_modes`: `driving`, `motorcycling`, `walking`, `running`, `hiking`, `cycling`, `bus`, `train`, `tram`, `subway`, `ferry`, `boating`, `flying`, `skiing`, `mixed`, `unknown` |
| `source_mode` | the source's own word (`IN_PASSENGER_VEHICLE`) |
| `distance_meters`, `probability` | nullable; `probability` is 0..1, how sure the source is of the mode |
| `title`, `description`, `described_by`, `tags` | |
| `source`, `source_ref`, `import_id` | as on visits |
| `path_points` | how many points `GET /trips/{id}/path` answers |
| `child_trips` | how many live trips name this one as their parent |

### `points`

Raw fixes, possibly millions; no public ids. `timestamp`, `latitude`,
`longitude`, `accuracy_meters` (nullable), `source` (`point_sources`:
`google_records`, `google_timeline_raw`, `google_timeline_path`,
`google_semantic_path`, `photo`), `source_ref` (a photo id for `photo`, else
`""`), `import_id`. A second copy of a fix from the same source is a no-op; a
point with a `source_ref` replaces an earlier one with the same source and
`source_ref`.

### `imports`

`id` (`import_000001`), `kind` (`import_kinds`: `google_timeline`,
`google_records`, `google_semantic`, `photo_store`), `sha256` (of the file;
unique), `source_machine`, `source_path`, `imported_at`, `finished_at` (0
while sending), `counts` (`points_added`, `points_unchanged`,
`places_created`, `places_updated`, `visits_created`, `visits_updated`,
`trips_created`, `trips_updated`, `records_unchanged`).

## Routes

Lists take `?limit=` (default 100, max 1000) and `?offset=`. A time range
`?from=&to=` selects what overlaps it: ends at or after `from`, starts before
`to`; either may be left out. `?include_deleted=true` shows deleted records.
A `PATCH` with an unknown key, or `null` on a field that is not a link, is a
400.

| Method | Path | |
|---|---|---|
| `GET` | `/health` | `{status, counts:{places, visits, trips, points, imports, deleted, revisions}}` |
| `GET` | `/settings` | the declared settings |
| `GET` | `/vocabularies` | `{place_kinds, trip_modes, record_sources, point_sources, import_kinds}` |
| `GET` | `/places` | oldest first. `?q=` (every word in name, address, description or tags), `?kind=`, `?tag=`, `?google_place_id=` |
| `POST` | `/places` | `{name, latitude, longitude, radius_meters, address, kind, description, described_by, tags}`; name, latitude, longitude, kind required. **201** |
| `GET` | `/places/{id}` | deleted is 404 unless `?include_deleted=true` |
| `PATCH` | `/places/{id}` | any of `name`, `latitude`, `longitude`, `radius_meters`, `address`, `kind`, `description`, `described_by`, `tags` |
| `GET` | `/visits` | earliest first. `?from=&to=`, `?place_id=`, `?source=`, `?tag=` |
| `POST` | `/visits` | `{place_id, started_at, ended_at, latitude, longitude, description, described_by, tags}`; a place or a coordinate required. **201** |
| `GET` | `/visits/{id}` | |
| `PATCH` | `/visits/{id}` | any of `place_id` (null clears), `description`, `described_by`, `tags` |
| `GET` | `/trips` | earliest first. `?from=&to=`, `?mode=`, `?parent_trip_id=` (an id lists its children; `none` lists trips that are part of no other), `?min_distance_meters=`, `?source=`, `?tag=`, `?exclude_parents=true` |
| `POST` | `/trips` | `{title, started_at, ended_at, mode, parent_trip_id, start_/end_ latitude/longitude, start_place_id, end_place_id, distance_meters, description, described_by, tags}`; title, times and mode required. Usually a road trip or holiday to attach imported drives to. **201** |
| `GET` | `/trips/{id}` | |
| `PATCH` | `/trips/{id}` | any of `parent_trip_id`, `start_place_id`, `end_place_id` (null clears each), `mode`, `title`, `description`, `described_by`, `tags`. A parent that would make a loop is a 400 |
| `GET` | `/trips/{id}/path` | `[{latitude, longitude, timestamp}]` in order; `timestamp` null when the source gave none |
| `POST` | `/trips/{id}/children` | `{"trip_ids":[…]}`: makes each part of this trip, all or none. Answers the parent |
| `DELETE` | `/places/{id}`, `/visits/{id}`, `/trips/{id}` | soft, **204**; again is a no-op |
| `POST` | `/places/{id}/restore`, `/visits/{id}/restore`, `/trips/{id}/restore` | answers the record |
| `GET` | `/places/{id}/revisions`, `/visits/{id}/revisions`, `/trips/{id}/revisions` | newest first: `{revision_id, entity_type, entity_id, reason, replaced_at, row}`, `row` the whole record as it was |
| `GET` | `/timeline` | visits and trips overlapping `?from=&to=`, in start order: `[{kind:"visit"\|"trip", visit, trip}]`. Trips that have children are left out unless `?include_parent_trips=true`. `?source=`; `?limit=` default 500, max 5000 |
| `GET` | `/points` | earliest first. `?from=` (inclusive) `&to=` (exclusive) on `timestamp`, `?source=`, `?source_ref=`; `?limit=` default 1000, max 10000 |
| `GET` | `/imports` | newest first |
| `POST` | `/imports` | `{kind, sha256, source_machine, source_path}`: starts an import, or answers the one with that hash. `{import, already_finished}` |
| `GET` | `/imports/{id}` | |
| `POST` | `/imports/{id}/batches` | `{places, visits, trips, points}` (`ImportBatch` in `location.go`), one transaction: one bad record is a 400 naming it and writes nothing. Places are written first, so a visit may name a place of the same batch by `google_place_id`. A finished import is a **409**. Answers what it did |
| `POST` | `/imports/{id}/finish` | stamps `finished_at`; the same file is then a no-op |

## Importing

`~/bin/location-import-google -file <path> [-source-machine <name>] [-dry-run]`
reads a Google export and posts it in batches of 2000. `<path>` is a file or a
directory; a directory imports every `.json` file under it, each as its own
import, and skips files that are not location history. It tells the formats
apart by their keys:

- **On-device Timeline export** (`Timeline.json`, 2024 on). `semanticSegments`:
  a `visit` at `hierarchyLevel` 0 becomes a visit, and its `topCandidate`
  (`placeId`, `semanticType`, `placeLocation`) a place; an `activity` becomes
  a trip (`topCandidate.type` → `mode`, `distanceMeters`), whose path is the
  `timelinePath` points inside its time; every `timelinePath` point is also a
  point. `rawSignals` `position`s become points. `userLocationProfile.frequentPlaces`
  become places (`label` → kind). Visits inside other visits
  (`hierarchyLevel` above 0) and `timelineMemory` are left out and counted.
- **Legacy Takeout `Records.json`**: `locations[]` become points
  (`latitudeE7`/`longitudeE7`, `timestamp` or `timestampMs`, `accuracy`).
- **Legacy Takeout `Semantic Location History/YYYY/YYYY_MONTH.json`**:
  `placeVisit` becomes a visit at its centre and, with a `placeId`, a place
  (`name`, `address`, `semanticType`); `activitySegment` becomes a trip whose
  path is `waypointPath`, or `simplifiedRawPath` when there are no waypoints.
  Every `simplifiedRawPath` point is also a point.

The iPhone export (a top-level JSON array) is refused. Records with no
readable time or place are left out and counted in the command's output.

`~/bin/location-import-photos [-date-sources exif,filename] [-dry-run]` reads
photo-store's photos with GPS (`GET /photos?has_gps=true`, paged) and makes
each into a point, source `photo`, `source_ref` the photo id. Only photos
dated by a trusted source are read: a photo dated when it was sent or posted
was taken somewhere earlier. The import's hash is that of the points it
would write, so running it again with nothing new changes nothing.
