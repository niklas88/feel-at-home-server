package strobe

import (
	"lamp/effect"
	"lamp/lampbase"
	"log"
	"time"
)

type StrobeConfig struct {
	Delay string
}

type Strobe struct {
	on    bool
	lamp  lampbase.DimLamp
	delay time.Duration
}

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Stroboscope",
			Description: "Stroboscope"},
		ConfigFactory: func() effect.Config {
			return &StrobeConfig{"30ms"}
		},
		Factory: effect.DimLampEffectFactory(NewStrobeEffect)})
}

func NewStrobeEffect(l lampbase.DimLamp) effect.Effect {
	return &Strobe{false, l, 30 * time.Millisecond}
}

func (s *Strobe) Configure(conf effect.Config) {
	strobeConf := conf.(*StrobeConfig)
	var err error

	s.delay, err = time.ParseDuration(strobeConf.Delay)
	if err != nil {
		log.Println(err)
		s.delay = 30 * time.Millisecond
	}
}

func (s *Strobe) Apply() (time.Duration, error) {
	var brightness uint8
	s.on = !s.on
	if s.on {
		brightness = 255
	} else {
		brightness = 0
	}
	err := s.lamp.SetBrightness(brightness)

	return s.delay, err
}
