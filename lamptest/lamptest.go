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
	reg               map[string]effects.StripeLampEffectFactory
)

func init() {
	reg = make(map[string]effects.StripeLampEffectFactory, 10)
	reg["fire"] = effects.StripeLampEffectFactory(effects.NewFireEffect)
	//reg["wheel"] = &effects.Wheel{}

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

	var eff effects.Effect = reg[effectName](lamp)
	config := eff.Config().(*effects.FireConfig)
	config.BottomColor = color.RGBA{255, 0, 0, 0}
	config.MidColor = color.RGBA{0, 0, 255, 0}
	config.TopColor = color.RGBA{0, 0, 0, 0}
	t := eff.Tomb()
	c := eff.ConfigChan()
	go eff.Apply()
	time.Sleep(3 * time.Second)
	c <- config
	go func(t *tomb.Tomb) {
		time.Sleep(30 * time.Second)
		fmt.Println("Killing")
		t.Kill(nil)
	}(t)

	t.Wait()
}
