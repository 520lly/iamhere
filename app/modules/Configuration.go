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
		Prefix          string
		Version         string
		Secret          string
		RandomItemLimit int
		BodySizeLimit   string
		Accounts        struct {
			Group string
		}

		Messages struct {
			Group string
		}

		Areas struct {
			Group string
		}

		Trails struct {
			Group string
		}
	}

	Database struct {
		Name     string
		User     string
		Password string
		Host     string
		Port     string
	}
}
