package locationstore

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/kayushkin/llm-bridge/msg"
	"github.com/kayushkin/llm-bridge/servicesettings"
)

// ServiceName is this service's name in its own settings description, as
// healthcheck and the repo know it.
const ServiceName = "location-store"

// OwnedEnvironmentVariablePrefix is the prefix of the variables that are this
// service's alone. A set variable carrying it that SettingDefinitions does not
// declare stops the service from starting.
const OwnedEnvironmentVariablePrefix = "LOCATION_STORE_"

// Keys of the settings, as GET /settings names them.
const (
	SettingListenAddress = "listen_address"
	SettingDataDirectory = "data_directory"
)

// DefaultListenAddress is loopback, deliberately: this is where the user has
// been, and the service has no auth of its own. dash is the front door.
const DefaultListenAddress = "127.0.0.1:8333"

// DefaultDataDir is where the database lives when the environment says
// nothing.
func DefaultDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".config", "location-store")
}

// SettingDefinitions declares every environment variable the service reads,
// once. Nothing is Editable: GET /settings is as open as every other route.
func SettingDefinitions() []servicesettings.Definition {
	return []servicesettings.Definition{
		{Key: SettingListenAddress, EnvironmentVariable: "LOCATION_STORE_ADDR", Kind: msg.ServiceSettingKindWiring, ValueType: msg.ServiceSettingValueTypeString, Default: DefaultListenAddress,
			Description: "The address the HTTP server listens on. Keep it on loopback: the service has no auth."},
		{Key: SettingDataDirectory, EnvironmentVariable: "LOCATION_STORE_DATA_DIR", Kind: msg.ServiceSettingKindPath, ValueType: msg.ServiceSettingValueTypeString, Default: DefaultDataDir(),
			Description: "The directory that holds location-store.db. Changing it starts the service on whatever is there, or nothing."},
	}
}

// NewSettingsRegistry reads this service's settings from environment.
func NewSettingsRegistry(environment servicesettings.Environment) (*servicesettings.Registry, error) {
	return servicesettings.New(ServiceName, []string{OwnedEnvironmentVariablePrefix}, SettingDefinitions(), environment)
}

// RegisterSettingsHandler serves the registry at GET /settings.
func RegisterSettingsHandler(mux *http.ServeMux, registry *servicesettings.Registry) {
	mux.Handle("GET /settings", servicesettings.Handler(registry, "/settings"))
}

// Keys of the import commands' settings. The commands read them from the
// environment; the service itself reads none.
const (
	SettingLocationStoreURL = "location_store_url"
	SettingPhotoStoreURL    = "photo_store_url"
)

func locationStoreURLDefinition() servicesettings.Definition {
	return servicesettings.Definition{Key: SettingLocationStoreURL, EnvironmentVariable: "LOCATION_STORE_URL", Kind: msg.ServiceSettingKindWiring, ValueType: msg.ServiceSettingValueTypeString, Default: "http://127.0.0.1:8333",
		Description: "Where location-store answers: the import is posted there."}
}

// ImportGoogleSettingDefinitions declares every environment variable
// location-import-google reads.
func ImportGoogleSettingDefinitions() []servicesettings.Definition {
	return []servicesettings.Definition{locationStoreURLDefinition()}
}

// NewImportGoogleSettingsRegistry reads location-import-google's settings. It
// owns no prefix, so only CheckRequired can refuse.
func NewImportGoogleSettingsRegistry(environment servicesettings.Environment) (*servicesettings.Registry, error) {
	return servicesettings.New("location-import-google", nil, ImportGoogleSettingDefinitions(), environment)
}

// ImportPhotosSettingDefinitions declares every environment variable
// location-import-photos reads.
func ImportPhotosSettingDefinitions() []servicesettings.Definition {
	return []servicesettings.Definition{
		locationStoreURLDefinition(),
		{Key: SettingPhotoStoreURL, EnvironmentVariable: "PHOTO_STORE_URL", Kind: msg.ServiceSettingKindWiring, ValueType: msg.ServiceSettingValueTypeString, Default: "http://127.0.0.1:8329",
			Description: "Where photo-store answers: the photos with GPS are read from there."},
	}
}

// NewImportPhotosSettingsRegistry reads location-import-photos' settings.
func NewImportPhotosSettingsRegistry(environment servicesettings.Environment) (*servicesettings.Registry, error) {
	return servicesettings.New("location-import-photos", nil, ImportPhotosSettingDefinitions(), environment)
}
