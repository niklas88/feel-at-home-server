package main

import (
	"fmt"
	"lamp/effects"
	"lamp/lampbase"
	"net"
	"time"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp4", "192.168.178.178:8888")
	if err != nil {
		fmt.Println("Couldn't resolve", err)
	}
	lamp := lampbase.NewLamp(4, 26, addr)

	lamp.Update()

	var eff effects.Effect = effects.NewFireEffect()

	for true {
		eff.ColorizeLamp(lamp)
		lamp.Update()
		time.Sleep(25 * time.Millisecond)
	}
}
