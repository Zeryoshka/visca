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

// CamDirectGainCommand set direct gain, val in 0..255 if RGain/BGain, else val in 0..14
func CamDirectGainCommand(c GainMode, val byte) Command {
	if c == GainChannel && val > 0x0E {
		return nil
	}
	p := (0xF0 & val) >> 4
	q := 0x0F & val
	return Command{0x01, 0x04, 0x40 | byte(c), 0x00, 0x00, p, q}
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

// CamDirectIrisCommand set direct iris position, val in 0x00..0x11
func CamDirectIrisCommand(val byte) Command {
	if val > 0x11 {
		return nil
	}
	p := (0xF0 & val) >> 4
	q := 0x0F & val
	return Command{0x01, 0x04, 0x4B, 0x00, 0x00, p, q}
}

func CamShutterCommand(s Selector) Command {
	return Command{0x01, 0x04, 0x0A, byte(s)}
}

// CamDirectShutterCommand set direct shutter position, val in 0x00..0x15
func CamDirectShutterCommand(val byte) Command {
	if val > 0x15 {
		return nil
	}
	p := (0xF0 & val) >> 4
	q := 0x0F & val
	return Command{0x01, 0x04, 0x4A, 0x00, 0x00, p, q}
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

func CamAFModeIntervalCommand(operationTime, stayingTime byte) Command {
	p := (0xF0 & operationTime) >> 4
	q := 0x0F & operationTime
	r := (0xF0 & stayingTime) >> 4
	s := 0x0F & stayingTime
	return Command{0x01, 0x04, 0x27, p, q, r, s}
}

// CamFocusDirectCommand set direct focus position. focusPosition in 0..0xFFFF
func CamFocusDirectCommand(focusPosition uint16) Command {
	p := (0xF000 & focusPosition) >> 12
	q := (0x0F00 & focusPosition) >> 8
	r := (0x00F0 & focusPosition) >> 4
	s := 0x000F & focusPosition
	return Command{0x01, 0x04, 0x48, byte(p), byte(q), byte(r), byte(s)}
}

func CamAFSensitivityNormalCommand() Command {
	return Command{0x01, 0x04, 0x58, 0x02}
}

func CamAFSensitivityLowCommand() Command {
	return Command{0x01, 0x04, 0x58, 0x03}
}

type HueColorSpecification byte

const (
	HueMaster HueColorSpecification = 0
	Magenta   HueColorSpecification = 1
	HueRed    HueColorSpecification = 2
	HueYellow HueColorSpecification = 3
	HueGreen  HueColorSpecification = 4
	HueCyan   HueColorSpecification = 5
	HueBlue   HueColorSpecification = 6
)

// CamColorHue set direct hue. level in range 0..0x0E. The initial value of level is 4
func CamColorHue(color HueColorSpecification, level byte) Command {
	if level > 0x0E {
		return nil
	}
	return Command{0x01, 0x04, 0x4F, 0x00, 0x00, byte(color), level}
}

func CamSaturation(level byte) Command {
	p := (0xF0 & level) >> 4
	q := 0x0F & level
	return Command{0x01, 0x04, 0xA1, 0x00, 0x00, p, q}
}
