package wormhole

type PipeID uint64

type ControllerPacket struct {
	PipeID  *PipeID
	State   *string
	Message []byte
}

type ControllerPacketReader interface {
	ReadControllerPacket() (ControllerPacket, error)
}

type ControllerPacketWriter interface {
	WriteControllerPacket(ControllerPacket) error
}

type ControllerPacketReadWriter interface {
	ControllerPacketReader
	ControllerPacketWriter
}
