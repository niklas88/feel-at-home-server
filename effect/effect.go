package effect

import (
	"lamp/lampbase"
)

type Config interface{}

type Info struct {
	Name        string
	Description string
	Config      Config
}

type Effect interface{}

type DeviceEffect func(p lampbase.Device, config Config) error
type DimLampEffect func(d lampbase.DimLamp, config Config) error
type ColorLampEffect func(c lampbase.ColorLamp, config Config) error
type StripeLampEffect func(s lampbase.StripeLamp, config Config) error
