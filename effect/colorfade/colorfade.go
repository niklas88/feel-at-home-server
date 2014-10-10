package colorefade

import (
	"github.com/pwaller/go-hexcolor"
	"image/color"
	"lamp/effect"
	"lamp/lampbase"
	"log"
	"math"
	"time"
)

type ColorfadeConfig struct {
	Color hexcolor.Hex
	Delay string
}

type Colorfade struct {
	color       color.RGBA
	currentStep uint8
	upward      bool
	lamp        lampbase.ColorLamp
	delay       time.Duration
}

const maxSteps = 255

func fadeColor(currentStep uint8, c color.RGBA) color.RGBA {
	var newColor color.RGBA
	multiplicator := (math.Pow(float64(currentStep)/maxSteps, 2.5) + float64(currentStep)/maxSteps) / 2 * maxSteps

	newColor.R = uint8(float64(c.R) / maxSteps * multiplicator)
	newColor.G = uint8(float64(c.G) / maxSteps * multiplicator)
	newColor.B = uint8(float64(c.B) / maxSteps * multiplicator)
	return newColor

}

func init() {
	effect.DefaultRegistry.Register(&effect.Registration{
		Info: effect.Info{
			Name:        "Colorfade",
			Description: "Fades with Color"},
		ConfigFactory: func() effect.Config { return &ColorfadeConfig{"#ffffff", "15ms"} },
		Factory:       effect.ColorLampEffectFactory(NewColorfadeEffect)})
}

func NewColorfadeEffect(l lampbase.ColorLamp) effect.Effect {
	return &Colorfade{currentStep: 0, upward: true, lamp: l, delay: 15 * time.Millisecond}
}

func (cf *Colorfade) Configure(conf effect.Config) {
	colorfadeConf := conf.(*ColorfadeConfig)
	var err error

	cf.delay, err = time.ParseDuration(colorfadeConf.Delay)
	if err != nil {
		log.Println(err)
		cf.delay = 30 * time.Millisecond
	}

	m := color.RGBAModel
	cf.color = m.Convert(colorfadeConf.Color).(color.RGBA)
}

func (cf *Colorfade) Apply() (time.Duration, error) {
	newColor := fadeColor(cf.currentStep, cf.color)

	if cf.upward {
		cf.currentStep++
	} else {
		cf.currentStep--
	}

	if cf.currentStep == maxSteps || cf.currentStep == 0 {
		cf.upward = !cf.upward
	}

	err := cf.lamp.SetColor(&newColor)
	return cf.delay, err

}
