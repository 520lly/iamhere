package modules

import ()

type Configuration struct {
	AppName   string
	Version   string
	AppConfig struct {
		Host          string
		Port          string
		LoggerPrefix  string
		EnableDebug   bool
		EnableSSL     bool
		EnableDumpLog bool
		LogPath       string
	}

	ApiConfig struct {
		Secret          string
		RandomItemLimit int
		BodySizeLimit   string
	}

	Database struct {
		Name     string
		User     string
		Password string
		Host     string
		Port     string
	}
}
