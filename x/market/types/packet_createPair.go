package types

// ValidateBasic is used for validating the packet
func (p CreatePairPacketData) ValidateBasic() error {

	// TODO: Validate the packet data

	return nil
}

// GetBytes is a helper for serialising
func (p CreatePairPacketData) GetBytes() ([]byte, error) {
	var modulePacket MarketPacketData

	modulePacket.Packet = &MarketPacketData_CreatePairPacket{&p}

	return modulePacket.Marshal()
}
