// Package lampbase provides interface definitions for different
// types of devices which can be controlled by a lampserver
package devices

import (
	"image/color"
	"time"
)

// Device describes any type of device that can be powered on or off
type Device interface {
	// Power powers a device on (on == true) or off (on == false)
	Power(on bool) error
}

// DimLamp describes a lamp device which can be dimmed. A DimLamp
// supports several types of effects that can be achieved using dimming
type DimLamp interface {
	// A DimLamp can also be powered on or off and therefore also supports the
	// Device interface
	Device
	// Brightness sets the static Brightness of the device
	Brightness(brightness uint8) error
	// Fade starts a fading effect with the specified delay and maximum
	// brightness
	Fade(delay time.Duration, maxBrightness uint8) error
	// Stroboscope starts a stroboscope like effect aka fast switching between
	// full and zero brightness
	Stroboscope(delay time.Duration) error
	// BrightnessScaling sets a scaling factor for all other effects so that they
	// may be used dimmed
	BrightnessScaling(brightness uint8) error
}

// ColorLamp describes a lamp device which supports single RGB color lighting.
// A ColorLamp supports several types of effects that can be achieved using
// colored light
type ColorLamp interface {
	// A ColorLamp can be dimmed by setting a darker color with the same hue so
	// it must also fullfill the DimLamp interface
	DimLamp
	// Color sets the ColorLamp's static color
	Color(color color.Color) error
	// ColorFade starts a fading effect that fades between dark and the specified
	// color
	ColorFade(delay time.Duration, col color.Color) error
	// Sunrise starts an effect imitating the colors of a sunrise, ideal for
	// waking up
	Sunrise(delay time.Duration) error
	// ColorWheel starts an effect looping through all available hues
	ColorWheel(delay time.Duration) error
}

// StripeLamp describes a device built from LED strips with individually
// controllable LEDs
type StripeLamp interface {
	// A StripeLamp may set the same color for all its LEDs and may thus act as a ColorLamp
	ColorLamp
	// Rainbow starts a rainbow like effect which spreads colors accross its LEDs
	Rainbow(delay time.Duration) error
	// RandomPixelBrightness starts an effect utilizing randomly bright but white
	// LEDs
	RandomPixelBrightness(delay time.Duration) error
	// RandomPixelWhiteFade starts an effect that fades a random selection of
	// LEDs using white light
	RandomPixelWhiteFade(delay time.Duration) error
	// RandomPixelColor starts an effect that sets a random color for each LED
	// every delay
	RandomPixelColor(delay time.Duration) error
}

// MatrixLamp describes a lamp device with pixels/LEDs ordered as a square
type MatrixLamp interface {
	// MatrixLamps may act as StripeLamps
	StripeLamp
	// Heart shows a static or animated heart on the device
	Heart() error
}

// Clock describes a device which receives time updates, usually to display the
// current time
type Clock interface {
	TimeUpdate(time time.Time) error
}

// WordClock describes a special clock device built from a matrix of
// pixels/LEDs
type WordClock interface {
	// WordClocks are built from a matrix so they support all matrix functions
	MatrixLamp
	// WordClocks are also clocks receiving TimeUpdate's
	Clock
	// Clock puts the WordClock into time display mode
	Clock() error
	// ClockColor sets the color used to display the time
	ClockColor(color color.Color) error
}
