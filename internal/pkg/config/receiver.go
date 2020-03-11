package config

type fileConfig struct {
	Enabled bool   `toml:"enabled"`
	Path    string `toml:"path"`
}

type udpConfig struct {
	Listen  string `toml:"listen"`
	Mode    string `toml:"mode"`
	Enabled bool   `toml:"enabled"`
}

type tcpConfig struct {
	Listen  string `toml:"listen"`
	Enabled bool   `toml:"enabled"`
}

type promConfig struct {
	Listen  string `toml:"listen"`
	Enabled bool   `toml:"enabled"`
}
