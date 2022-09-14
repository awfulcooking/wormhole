package wormhole

type ServerConfig struct {
	NameGenerator      NameGenerator
	WebsocketReadLimit int64
}

func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		NameGenerator:      StaticNameGenerator("host"),
		WebsocketReadLimit: 512 * 1024,
	}
}
