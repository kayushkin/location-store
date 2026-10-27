# About location-store

## What it owns

`:8333`, loopback only. Where the user has been: named places (`place_000001`), visits to them (`visit_000001`), trips between them (`trip_000001`), the raw location points underneath (millions possible, no public ids, read by time) and the imports that brought them in (`import_000001`). Place, visit and trip ids are registered in kanban-store's entity-type registry as types `place`, `visit` and `trip`, so a journal-store timeline item or a note cites a particular drive or trip as a chip. **Location history is private**: a front door such as dash adds the auth this service has none of, and nothing here is copied to noteboard. `CONTRACT.md` is the route table.

# How it works

## Imports and what they never touch

`location-import-google -file <file or directory>` reads Google's on-device Timeline export (`Timeline.json`), legacy Takeout `Records.json`, and legacy Semantic Location History month files, told apart by their top-level keys (`internal/googletimeline`). `location-import-photos` makes each photo-store photo with GPS and a trusted date (`exif`, `filename` by default) into a point with source `photo` and `source_ref` the photo id. Both go through `internal/importclient`: `POST /imports` with the file's sha256 (a finished hash is a no-op, an unfinished one resumes), `POST /imports/{id}/batches` of 2000 records in one transaction each, `POST /imports/{id}/finish`. Records are found again by key: a place by `google_place_id`, a visit or trip by `source` + `source_ref` (`<started_at>-<ended_at>`), a point by time, coordinates and source. **An import never touches** a name, title, description, tags or `parent_trip_id`, never changes a place link already set, never changes a `mode` someone corrected (one that no longer matches `ModeForGoogleActivityType(source_mode)`), and leaves a deleted record deleted. Every update, import change, delete and restore of a place, visit or trip first copies the whole record to `revisions`. Grouping drives into a road trip is by hand: `POST /trips` for the parent, then `POST /trips/{id}/children`. Visits inside other visits (`hierarchyLevel` above 0), `timelineMemory` and the iPhone export (a top-level array) are not read. Which files and machines this host's history came from is not in this repo: see the noteboard note tagged `location-store`.

# Working in this repo

## Build, test and deploy

Module `github.com/kayushkin/location-store`, root package `locationstore`, server `cmd/location-store`, importers `cmd/location-import-google` and `cmd/location-import-photos`. No build tags. `go.mod` replaces `../llm-bridge` (settings) and `../photo-store` (the photo importer reads `photostore.Photo`). Every variable it reads is declared once in `settings.go` with llm-bridge `servicesettings`, served at `GET /settings`; `TestEveryEnvironmentVariableTheServiceReadsIsDeclared` fails on an undeclared `os.Getenv`. `./deploy.sh` runs the shared deploy gate, tests, checks the running unit's environment, installs `location-store` and both importers to `~/bin`, restarts the user unit and fails unless `/health` is ok and the port is bound to loopback only. `./generate-ts.sh` renders `location.go` and `vocabulary.go` to `ts/model.ts`, published source-only as `@kayushkin/location-store-types` (`file:../location-store/ts`); bridge-ui's place, visit and trip chips read it.

## Public repository: nothing personal

This repo is public. Nothing about the user's own data goes in it — not in code, tests, comments, commit messages or this prompt: no real coordinates or addresses, no names of places the user goes, no machine names, file paths of their exports, account names, home directories or tokens. Tests use invented data on a grid near 10°N 20°E with made-up place ids (`example-home`). A fact an agent here needs about this host's data goes in a noteboard note tagged `location-store`.
