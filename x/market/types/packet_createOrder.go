package types

// ValidateBasic is used for validating the packet
func (p CreateOrderPacketData) ValidateBasic() error {

	// TODO: Validate the packet data

	return nil
}

// GetBytes is a helper for serialising
func (p CreateOrderPacketData) GetBytes() ([]byte, error) {
	var modulePacket MarketPacketData

	modulePacket.Packet = &MarketPacketData_CreateOrderPacket{&p}

	return modulePacket.Marshal()
}
