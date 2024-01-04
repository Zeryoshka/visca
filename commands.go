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
	GainChannel  GainMode = 0x0C
)

func CamGainCommand(c GainMode, s Selector) Command {
	return Command{0x01, 0x04, byte(c), byte(s)}
}

// CamDirectGainCommand set direct gain, val in 0..255
func CamDirectGainCommand(c GainMode, val byte) Command {
	return Command{0x01, 0x04, 0x40 | byte(c), val}
}

// CamDirectSpeedCommand set color speed, val in 1(slow)..5(fast)
func CamDirectSpeedCommand(speed byte) Command {
	if 0 < speed && speed <= 5 {
		return Command{0x01, 0x04, 0x56, speed}
	}
	return nil
}

// CamChromaSuppressCommand set color chroma suppress, val in 0(off), 1(weak)..3(strong)
func CamChromaSuppressCommand(level byte) Command {
	if level <= 3 {
		return Command{0x01, 0x04, 0x5F, level}
	}
	return nil
}

func CamIrisCommand(s Selector) Command {
	return Command{0x01, 0x04, 0x0B, byte(s)}
}

func CamDirectIrisCommand(val byte) Command {
	return Command{0x01, 0x04, 0x4B, val}
}

func CamShutterCommand(s Selector) Command {
	return Command{0x01, 0x04, 0x0A, byte(s)}
}

func CamDirectShutterCommand(val byte) Command {
	return Command{0x01, 0x04, 0x4A, val}
}

type ExposureMode byte

const (
	AeFullAuto        ExposureMode = 0x00
	AeManual          ExposureMode = 0x03
	AeShutterPriority ExposureMode = 0x0A
	AeIrisPriority    ExposureMode = 0x0B
	AeBright          ExposureMode = 0x0D
)

func CamAECommand(mode ExposureMode) Command {
	return Command{0x01, 0x04, 0x39, byte(mode)}
}

func CamBackLightOnCommand() Command {
	return Command{0x01, 0x04, 0x33, 0x02}
}

func CamBackLightOffCommand() Command {
	return Command{0x01, 0x04, 0x33, 0x03}
}

// Cam2dNR set 2D noise reduction, val in 0(off), 1(weak)..5(strong)
func Cam2dNR(level byte) Command {
	if level <= 5 {
		return Command{0x01, 0x04, 0x53, level}
	}
	return nil
}

// Cam3dNR set 3D noise reduction, val in 0(off), 1(weak)..5(strong)
func Cam3dNR(level byte) Command {
	if level <= 5 {
		return Command{0x01, 0x04, 0x54, level}
	}
	return nil
}

func CamFocusAutoModeCommand() Command {
	return Command{0x01, 0x04, 0x38, 0x02}
}

func CamFocusManualModeCommand() Command {
	return Command{0x01, 0x04, 0x38, 0x03}
}

func CamFocusToggleModeCommand() Command {
	return Command{0x01, 0x04, 0x38, 0x10}
}

func CamFocusOnePushTriggerCommand() Command {
	return Command{0x01, 0x04, 0x18, 0x01}
}

type AFMode byte

const (
	AFNormal      AFMode = 0x00
	AFInterval    AFMode = 0x01
	AFZoomTrigger AFMode = 0x02
)

func CamAFModeCommand(mode AFMode) Command {
	return Command{0x01, 0x04, 0x57, byte(mode)}
}

func CamAFModeIntervalCommand(val byte) Command {
	return Command{0x01, 0x04, 0x27, val}
}

func CamAFModeDirectCommand(focusPosition []byte) Command {
	if len(focusPosition) > 4 {
		return nil
	}
	return append([]byte{0x01, 0x04, 0x48}, focusPosition...)
}

func CamAFSensitivityNormalCommand() Command {
	return Command{0x01, 0x04, 0x58, 0x02}
}

func CamAFSensitivityLowCommand() Command {
	return Command{0x01, 0x04, 0x58, 0x03}
}
