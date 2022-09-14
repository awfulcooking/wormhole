package wormhole

type PipeID uint64

type ControllerPacket struct {
	PipeID      PipeID          `json:"pipeID,omitempty"`
	State       PipeState       `json:"state,omitempty"`
	Data        *string         `json:"data,omitempty"`
	DataType    DataType        `json:"dataType,omitempty"`
	Subprotocol string          `json:"subprotocol,omitempty"`
	Meta        *ControllerMeta `json:"meta,omitempty"`
}

type ControllerMeta struct {
	Name string `json:"name"`
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
