package web

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	TZDBFile      = "tzdata/timezone.db"
	TZVersionFile = "tzdata/version.json"

	GeoDataURLTemplate        = "https://github.com/evansiroky/timezone-boundary-builder/releases/download/%s/timezones-with-oceans.geojson.zip"
	GeoDataReleaseURLTemplate = "https://github.com/evansiroky/timezone-boundary-builder/releases/tag/%s"
)

type TzRelease struct {
	Version    string `json:"version"`
	URL        string `json:"url"`
	GeoDataURL string `json:"geo_data_url"`
}

func NewTzRelease(version string) TzRelease {
	return TzRelease{
		Version:    version,
		URL:        fmt.Sprintf(GeoDataReleaseURLTemplate, version),
		GeoDataURL: fmt.Sprintf(GeoDataURLTemplate, version),
	}
}

// TzSchema configuration
type TzSchema struct {
	DatabaseName string `mapstructure:"database_name"`
	VersionFile  string `mapstructure:"version_file"`
}

// WebSchema configuration
type WebSchema struct {
	ListenAddress      string `mapstructure:"listen_address,omitempty"`
	AuthTokenValue     string `mapstructure:"auth_token_value,omitempty"`
	AuthTokenParamName string `mapstructure:"auth_token_param_name,omitempty"`
}

// ConfigSchema main configuration for the news room
type ConfigSchema struct {
	Tz             TzSchema  `mapstructure:"tz"`
	Web            WebSchema `mapstructure:"web"`
	RuntimeVersion string    `mapstructure:"-"`
}

// Defaults configure defaults
func Defaults() {
	// tz defaults
	viper.SetDefault("tz.database_name", TZDBFile)
	viper.SetDefault("tz.version_file", TZVersionFile)
	// web
	viper.SetDefault("web.listen_address", ":2004")
	viper.SetDefault("web.auth_token_value", "") // GEO2TZ_WEB_AUTH_TOKEN_VALUE="ciao"
	viper.SetDefault("web.auth_token_param_name", "t")
}

// Validate a configuration
func Validate(_ *ConfigSchema) (err []error) {
	// TODO: implement this one
	return
}

// Settings general settings
var Settings ConfigSchema
