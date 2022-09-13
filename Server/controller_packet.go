package wormhole

type PipeID uint64

type ControllerMeta struct {
	Slug string `json:"slug"`
}

type ControllerPacket struct {
	PipeID      PipeID          `json:"pipeID,omitempty"`
	State       PipeState       `json:"state,omitempty"`
	Data        *string         `json:"data,omitempty"`
	DataType    PipeDataType    `json:"dataType,omitempty"`
	Subprotocol string          `json:"subprotocol,omitempty"`
	Meta        *ControllerMeta `json:"meta,omitempty"`
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
