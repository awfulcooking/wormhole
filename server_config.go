package wormhole

import (
	"embed"
	"io/fs"
)

type ServerConfig struct {
	NameGenerator      NameGenerator
	StaticFS           fs.FS
	WebsocketReadLimit int64
}

func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		NameGenerator:      StaticNameGenerator("host"),
		WebsocketReadLimit: 512 * 1024,
		StaticFS:           DefaultStaticFS(),
	}
}

//go:embed html/*
var defaultStaticFS embed.FS

func DefaultStaticFS() fs.FS {
	fs, _ := fs.Sub(defaultStaticFS, "html")
	return fs
}
