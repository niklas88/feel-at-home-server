package lampbase

import (
	"image/color"
	"time"
)

type Device interface {
	Power(on bool) error
}

type DimLamp interface {
	Device
	Brightness(brightness uint8) error
	Fade(delay time.Duration, maxBrightness uint8) error
	Stroboscope(delay time.Duration) error
	BrightnessScaling(brightness uint8) error
}

type ColorLamp interface {
	DimLamp
	Color(color color.Color) error
	ColorFade(delay time.Duration, col color.Color) error
	Sunrise(delay time.Duration) error
	ColorWheel(delay time.Duration) error
}

type StripeLamp interface {
	ColorLamp
	Rainbow(delay time.Duration) error
	RandomPixelBrightness(delay time.Duration) error
	RandomPixelWhiteFade(delay time.Duration) error
	RandomPixelColor(delay time.Duration) error
}

type MatrixLamp interface {
	StripeLamp
	Heart() error
}

type Clock interface {
	TimeUpdate(time time.Time) error
}

type WordClock interface {
	MatrixLamp
	Clock
	Clock() error
	ClockColor(color color.Color) error
}
