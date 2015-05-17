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
	return &Strobe{l, 30 * time.Millisecond}
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
	err := s.lamp.Stroboscope(s.delay)
	return -1, err
}
