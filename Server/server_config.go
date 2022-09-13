package wormhole

type ServerConfig struct {
	WebsocketReadLimit int64
}

func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		WebsocketReadLimit: 512 * 1024,
	}
}
