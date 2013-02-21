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
	reg               map[string]*effects.Info
)

func init() {
	reg = make(map[string]*effects.Info, 10)
	reg["fire"] = &effects.Info{
		Name:          "fire",
		ConfigFactory: func() effects.Config { return &effects.FireConfig{} },
		Factory:       effects.StripeLampEffectFactory(effects.NewFireEffect)}
	reg["wheel"] = &effects.Info{
		Name:          "wheel",
		ConfigFactory: func() effects.Config { return nil },
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

	var eff effects.Effect = reg[effectName].CreateEffect(lamp)
	if fire, ok := eff.(*effects.FireEffect); ok {
		fire.Configure(&effects.FireConfig{color.RGBA{255, 0, 0, 0}, color.RGBA{0, 0, 255, 0}, color.RGBA{0, 0, 0, 0}})
	}
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
	fmt.Println(lamp.Power(false))
}
