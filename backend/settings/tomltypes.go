package settings

import "time"

type TomlApp struct {
	Debug           bool
	Swag            bool
	FilrepApi       string
	FilecoinApi     string
	JwtSecret       string
	PasswordSalt    string
	PublishDate     time.Time
	OfficialWebsite string
	OfficialEmail   string
}

type TomlServer struct {
	HttpPort     int
	ReadTimeout  int
	WriteTimeout int
}

type TomlDatabase struct {
	Type        string
	Dsn         string
	TablePrefix string
}

type TomlSmtp struct {
	Basic TomlBasicSmtp
}

type TomlBasicSmtp struct {
	Host     string
	User     string
	Password string
}

type AppConfig struct {
	App      TomlApp
	Server   TomlServer
	Database TomlDatabase
	Smtp     TomlSmtp
}
