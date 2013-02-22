package main

import (
	"flag"
	"fmt"
	"image/color"
	"lamp/effect"
	"lamp/effects/fire"
	_ "lamp/effects/wheel"
	"lamp/lampbase"
	"launchpad.net/tomb"
	"net"
	"time"
)

var (
	effectName        string
	lampAddress       string
	lampStripes       int
	lampLedsPerStripe int
	lampDelay         int
)

func init() {
	flag.StringVar(&effectName, "effect", "fire", "Effect")
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

	effectInfo := effect.Info{Name: effectName}
	eff := effect.DefaultRegistry.CreateEffect(&effectInfo, lamp)
	if fireEffect, ok := eff.(*fire.FireEffect); ok {
		fireEffect.Configure(&fire.FireConfig{color.RGBA{255, 0, 0, 0}, color.RGBA{0, 0, 255, 0}, color.RGBA{0, 0, 0, 0}})
	}
	controller := effect.NewController()
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
