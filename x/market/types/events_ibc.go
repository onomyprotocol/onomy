package types

// IBC events
const (
	EventTypeTimeout = "timeout"
	// this line is used by starport scaffolding # ibc/packet/event
	EventTypeCreateOrderPacket = "createOrder_packet"

	EventTypeCreatePairPacket = "createPair_packet"

	AttributeKeyAckSuccess = "success"
	AttributeKeyAck        = "acknowledgement"
	AttributeKeyAckError   = "error"
)
