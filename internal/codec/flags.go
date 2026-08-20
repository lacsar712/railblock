package codec

// FlagSet provides helpers for constructing and inspecting frame flag bytes.
type FlagSet uint8

func NewFlagSet() FlagSet { return FlagSet(FlagNone) }

func (fs FlagSet) WithForceClear() FlagSet {
	return FlagSet(uint8(fs) | FlagForceClear)
}

func (fs FlagSet) WithPriority() FlagSet {
	return FlagSet(uint8(fs) | FlagPriority)
}

func (fs FlagSet) WithEmergency() FlagSet {
	return FlagSet(uint8(fs) | FlagEmergency)
}

func (fs FlagSet) Byte() uint8 { return uint8(fs) }

func (fs FlagSet) HasForceClear() bool { return uint8(fs)&FlagForceClear != 0 }

func (fs FlagSet) HasPriority() bool { return uint8(fs)&FlagPriority != 0 }

func (fs FlagSet) HasEmergency() bool { return uint8(fs)&FlagEmergency != 0 }
