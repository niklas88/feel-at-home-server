package lampbase

import (
	"image/color"
)

type Device interface {
	Power(on bool) error
}

type DimLamp interface {
	Device
	SetBrightness(brightness uint8) error
}

type ColorLamp interface {
	DimLamp
	SetColor(color color.Color) error
}

type StripeLamp interface {
	ColorLamp
	Stripes() []Stripe
	UpdateAll() error
}

type Stripe []color.RGBA
