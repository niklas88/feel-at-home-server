// The effect package defines data structures and functions to export the
// functionality of devices, giving functions a textual name, description and
// providing an interface to configure their parameters
package effect

import (
	"github.com/niklas88/feel-at-home-server/lampbase"
)

// Config is a type used to configure the parameters of a function supported by
// a lamp this could be a color to be set or define the speed of an animation
type Config interface{}

// EmptyConfig is used to configure effects that don't have any parameters
type EmptyConfig struct{}

// DelayConfig sets the Delay between steps of an animation
type DelayConfig struct {
	Delay string
}

// DelayConfigFactory is used to create a preconfigured DelayConfig that
// configures a delay of 30ms appropriate for smooth but slow animations
func DelayConfigFactory() Config {
	return &DelayConfig{"30ms"}
}

// EmptyConfigFactory creates an EmptyEffectConfig
func EmptyConfigFactory() Config {
	return &EmptyConfig{}
}

// Info holds metadata on an effect including its current configuration
type Info struct {
	Name        string
	Description string
	Config      Config
}

// Effect are all types that provide a configureable Apply function
type Effect interface {
	Apply(config Config) error
}

// EffectFunc turns any function taking a config into an Effect
type EffectFunc func(config Config) error

// Apply let's EffectFunc implement the Effect interface
func (f EffectFunc) Apply(config Config) error {
	return f(config)
}

type DeviceEffectFactory func(p lampbase.Device) Effect
type DimLampEffectFactory func(d lampbase.DimLamp) Effect
type ColorLampEffectFactory func(c lampbase.ColorLamp) Effect
type StripeLampEffectFactory func(s lampbase.StripeLamp) Effect
type MatrixLampEffectFactory func(s lampbase.MatrixLamp) Effect
type WordClockEffectFactory func(s lampbase.WordClock) Effect
