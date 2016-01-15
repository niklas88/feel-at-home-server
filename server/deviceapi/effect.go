// The effect package defines data structures and functions to export the
// functionality of devices, giving functions a textual name, description and
// providing an interface to configure their parameters
package deviceapi

import (
	"github.com/niklas88/feel-at-home-server/devices"
)

// Config is a type used to configure the parameters of a function supported by
// a lamp this could be a color to be set or define the speed of an animation
type Config interface{}

type ConfigFactory func() Config

// EmptyConfig is used to configure effects that don't have any parameters
type EmptyConfig struct{}

// DelayConfig sets the Delay between steps of an animation
type DelayConfig struct {
	Delay string
}

// Effect are all types that provide a configureable Apply function
type Effect interface {
	Name() string
	Description() string
	DefaultConfig() Config
	Apply(d devices.Device, config Config) error
	Compatible(d devices.Device) bool
}
