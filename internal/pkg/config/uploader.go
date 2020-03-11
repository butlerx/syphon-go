package config

type uploaderConfig struct {
	File []fileUploadConfig `toml:"file"`
	Grpc []grpcUploadConfig `toml:"grpc"`
	TCP  []tcpUploadConfig  `toml:"tcp"`
	UDP  []udpUploadConfig  `toml:"udp"`
}

type tcpUploadConfig struct {
	Enabled bool   `toml:"enabled"`
	Host    string `toml:"host"`
	Port    int64  `toml:"port"`
	Pattern string `toml:"pattern"`
}

type grpcUploadConfig struct {
	Enabled bool   `toml:"enabled"`
	Host    string `toml:"host"`
	Port    int64  `toml:"port"`
	Pattern string `toml:"pattern"`
}

type udpUploadConfig struct {
	Enabled bool   `toml:"enabled"`
	Host    string `toml:"host"`
	Port    int64  `toml:"port"`
	Pattern string `toml:"pattern"`
}

type fileUploadConfig struct {
	Enabled bool   `toml:"enabled"`
	Path    string `toml:"path"`
	Pattern string `toml:"pattern"`
}
