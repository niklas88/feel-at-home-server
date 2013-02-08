package main

import (
	"flag"
	"fmt"
	"lamp/effects"
	"lamp/lampbase"
	"net"
	"strings"
	"time"
)

var (
	effectName        string
	lampAddress       string
	lampStripes       int
	lampLedsPerStripe int
	reg               map[string]effects.Effect
)

func init() {
	reg = make(map[string]effects.Effect, 10)
	reg["fire"] = effects.NewFireEffect()
	reg["wheel"] = &effects.Wheel{}

	effectList := make([]string, 0, 10)
	for k, _ := range reg {
		effectList = append(effectList, k)
	}
	flag.StringVar(&effectName, "effect", "fire", "Effect (available: "+strings.Join(effectList, ", ")+")")
	flag.StringVar(&lampAddress, "lamp", "192.168.178.178:8888", "Address of the lamp")
	flag.IntVar(&lampStripes, "stripes", 4, "Number of stripes the lamp has")
	flag.IntVar(&lampLedsPerStripe, "leds", 26, "Number of LEDs per stripe")
}

func main() {
	flag.Parse()
	addr, err := net.ResolveUDPAddr("udp4", lampAddress)
	if err != nil {
		fmt.Println("Couldn't resolve", err)
	}
	lamp := lampbase.NewLamp(lampStripes, lampLedsPerStripe, addr)

	lamp.Update()

	var eff effects.Effect = reg[effectName]

	for true {
		eff.ColorizeLamp(lamp)
		lamp.Update()
		time.Sleep(25 * time.Millisecond)
	}
}
