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
	SetBrightness(brightness uint8) error
	Fade(delay time.Duration, maxBrightness uint8) error
	Stroboscope(delay time.Duration) error	
}

type ColorLamp interface {
	DimLamp
	SetColor(color color.Color) error
	ColorFade(delay time.Duration, col color.Color) error
	Sunrise(delay time.Duration) error
	ColorWheel(delay time.Duration) error
}

type StripeLamp interface {
	ColorLamp
	Stripes() []Stripe
	UpdateAll() error
	Rainbow(delay time.Duration) error
	RandomPixelBrightness(delay time.Duration) error
	RandomPixelWhiteFade(delay time.Duration) error
	RandomPixelColor(delay time.Duration) error
}

type Stripe []color.RGBA
