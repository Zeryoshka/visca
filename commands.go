package visca

type Command []byte

type WhiteBalance byte

const (
	WbAuto1   WhiteBalance = 0x00
	WbIndoor  WhiteBalance = 0x01
	WbOutdoor WhiteBalance = 0x02
	WbOnePush WhiteBalance = 0x03
	WbAuto2   WhiteBalance = 0x04
	WbManual  WhiteBalance = 0x05
)

// CamWbCommand build ViscaCommand to manage WhiteBalance
func CamWbCommand(mode WhiteBalance) Command {
	return Command{0x01, 0x04, 0x35, byte(mode)}
}

type Selector byte

const (
	Reset Selector = 0x00
	Up    Selector = 0x02
	Down  Selector = 0x03
)

type GainMode byte

const (
	RGainChannel GainMode = 0x03
	BGainChannel GainMode = 0x04
)

func CamGainCommand(c GainMode, s Selector) Command {
	return Command{0x01, 0x04, byte(c), byte(s)}
}

func CamDirectGainCommand(c GainMode, val int8) Command {
	return Command{0x01, 0x04, 0x40 | byte(c), byte(val)}
}
