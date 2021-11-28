package server

import (
	"github.com/spf13/viper"
)

// TzSchema configuration
type TzSchema struct {
	DatabaseType       string `mapstructure:"database_type"`
	DatabaseName       string `mapstructure:"database_name"`
	Snappy             bool   `mapstructure:"snappy"`
	Encoding           string `mapstructure:"encoding"`
	DownloadTzData     bool   `mapstructure:"download_tz_data"`
	DownloadTzDataURL  string `mapstructure:"download_tz_data_url"`
	DownloadTzFilename string `mapstructure:"download_tz_filename"`
}

// WebSchema configuration
type WebSchema struct {
	ListenAddress      string  `mapstructure:"listen_address,omitempty"`
	AuthTokenValue     string  `mapstructure:"auth_token_value,omitempty"`
	AuthTokenParamName string  `mapstructure:"auth_token_param_name,omitempty"`
	RateLimit          float64 `mapstructure:"rate_limit,omitempty"`
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
	viper.SetDefault("tz.database_type", "boltdb")
	viper.SetDefault("tz.database_name", "timezone")
	viper.SetDefault("tz.snappy", true)
	viper.SetDefault("tz.encoding", "msgpack")
	viper.SetDefault("tz.download_tz_data", true)
	viper.SetDefault("tz.download_tz_data_url", "https://api.github.com/repos/evansiroky/timezone-boundary-builder/releases/latest")
	viper.SetDefault("tz.download_tz_filename", "timezones.geojson.zip")
	// web
	viper.SetDefault("web.listen_address", ":2004")
	viper.SetDefault("web.auth_token_value", "") // GEO2TZ_WEB_AUTH_TOKEN_VALUE="ciao"
	viper.SetDefault("web.auth_token_param_name", "t")
	viper.SetDefault("web.rate_limit", 0)

}

// Validate a configuration
func Validate(schema *ConfigSchema) (err []error) {
	// TODO: implement this one
	return
}

// Settings general settings
var Settings ConfigSchema
