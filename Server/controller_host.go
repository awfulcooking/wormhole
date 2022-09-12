package wormhole

import "io"

type ControllerHost interface {
	ControllerPacketReadWriter
	io.Closer
}
