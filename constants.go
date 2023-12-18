package visca

import "errors"

var IncorrectDeviceIndexErr error = errors.New("incorrect index, should be in [1; 7], and be free")

type requestType int

const (
	commandRequest requestType = iota
)

const HeaderLen = 8

var CommandHeaderPrefix = []byte{0x01, 0x00}
