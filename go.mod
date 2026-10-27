module github.com/kayushkin/location-store

go 1.26.0

require (
	github.com/kayushkin/llm-bridge v0.0.0
	github.com/kayushkin/photo-store v0.0.0
	github.com/mattn/go-sqlite3 v1.14.37
)

replace (
	github.com/kayushkin/llm-bridge => ../llm-bridge
	github.com/kayushkin/photo-store => ../photo-store
)
