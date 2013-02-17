package main

import (
	"flag"
	"fmt"
	"image/color"
	"lamp/effects"
	"lamp/lampbase"
	"launchpad.net/tomb"
	"net"
	"strings"
	"time"
)

var (
	effectName        string
	lampAddress       string
	lampStripes       int
	lampLedsPerStripe int
	lampDelay         int
	reg               map[string]*effects.EffectInfo
)

func configureEffect(info *effects.EffectInfo, lamp lampbase.Powerable) effects.Effect {
	if info == nil {
		return nil
	}

	var config interface{}
	switch info.Name {
	case "fire":
		config = &effects.FireConfig{color.RGBA{255, 0, 0, 0}, color.RGBA{0, 0, 255, 0}, color.RGBA{0, 0, 0, 0}}
	default:
		config = nil
	}

	switch fac := info.Factory.(type) {
	case effects.PowerableEffectFactory:
		if l, ok := lamp.(lampbase.Powerable); ok {
			return fac(l, config)
		}
	case effects.DimLampEffectFactory:
		if l, ok := lamp.(lampbase.DimLamp); ok {
			return fac(l, config)
		}
	case effects.ColorLampEffectFactory:
		if l, ok := lamp.(lampbase.ColorLamp); ok {
			return fac(l, config)
		}
	case effects.StripeLampEffectFactory:
		if l, ok := lamp.(lampbase.StripeLamp); ok {
			return fac(l, config)
		}
	default:
		panic("Unknow lamp factory type")
	}
	return nil
}

func init() {
	reg = make(map[string]*effects.EffectInfo, 10)
	reg["fire"] = &effects.EffectInfo{
		Name:          "fire",
		ConfigFactory: func() interface{} { return &effects.FireConfig{} },
		Factory:       effects.StripeLampEffectFactory(effects.NewFireEffect)}
	reg["wheel"] = &effects.EffectInfo{
		Name:          "wheel",
		ConfigFactory: func() interface{} { return nil },
		Factory:       effects.ColorLampEffectFactory(effects.NewWheelAllEffect)}

	effectList := make([]string, 0, 10)
	for k, _ := range reg {
		effectList = append(effectList, k)
	}
	flag.StringVar(&effectName, "effect", "fire", "Effect (available: "+strings.Join(effectList, ", ")+")")
	flag.StringVar(&lampAddress, "lamp", "192.168.178.178:8888", "Address of the lamp")
	flag.IntVar(&lampStripes, "stripes", 4, "Number of stripes the lamp has")
	flag.IntVar(&lampLedsPerStripe, "leds", 26, "Number of LEDs per stripe")
	flag.IntVar(&lampDelay, "delay", 25, "Milliseconds between updates")
}

func main() {
	flag.Parse()
	addr, err := net.ResolveUDPAddr("udp4", lampAddress)
	if err != nil {
		fmt.Println("Couldn't resolve", err)
	}
	lamp := lampbase.NewUdpStripeLamp(lampStripes, lampLedsPerStripe)
	if err = lamp.Dial(nil, addr); err != nil {
		fmt.Println(err)
		return
	}
	defer lamp.Close()

	lamp.UpdateAll()

	var eff effects.Effect = configureEffect(reg[effectName], lamp)
	controller := effects.NewController()
	go controller.Run()
	time.Sleep(3 * time.Second)
	controller.EffectChan <- eff
	go func(t *tomb.Tomb) {
		time.Sleep(30 * time.Second)
		fmt.Println("Killing")
		t.Kill(nil)
	}(&controller.Tomb)

	controller.Tomb.Wait()
	lamp.Power(false)
}
