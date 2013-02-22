package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"lamp/effect"
	_ "lamp/effects/fire"
	_ "lamp/effects/wheel"
	"lamp/lampbase"
	"log"
	"net"
	"net/http"
	"strconv"
)

type DeviceInfo struct {
	Name   string
	Device lampbase.Device `json:"-"`
}

var (
	lampAddress       string
	lampStripes       int
	lampLedsPerStripe int
	lampDelay         int
	deviceList        []DeviceInfo
)

func init() {
	flag.StringVar(&lampAddress, "lamp", "192.168.178.178:8888", "Address of the lamp")
	flag.IntVar(&lampStripes, "stripes", 4, "Number of stripes the lamp has")
	flag.IntVar(&lampLedsPerStripe, "leds", 26, "Number of LEDs per stripe")
	flag.IntVar(&lampDelay, "delay", 25, "Milliseconds between updates")
}

func DeviceListHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Requesting device list")
	w.Header().Set("Content-Type", "application/json")
	out, err := json.Marshal(deviceList)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(out)
}

type EffectDescription struct {
	Name   string
	Config interface{}
}

func EffectListHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	id := vars["id"]
	deviceId, _ := strconv.Atoi(id)
	if deviceId < 0 || deviceId >= len(deviceList) {
		http.NotFound(w, req)
		return
	}
	device := deviceList[deviceId]
	effectList := effect.DefaultRegistry.CompatibleEffects(device.Device)
	out, err := json.Marshal(effectList)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(out)
}

func main() {
	flag.Parse()
	addr, err := net.ResolveUDPAddr("udp4", lampAddress)
	if err != nil {
		log.Println("Couldn't resolve", err)
	}
	lamp := lampbase.NewUdpStripeLamp(lampStripes, lampLedsPerStripe)
	if err = lamp.Dial(nil, addr); err != nil {
		log.Println(err)
		return
	}
	defer lamp.Close()

	deviceList = append(deviceList, DeviceInfo{"Big Lamp", lamp})

	r := mux.NewRouter()
	r.HandleFunc("/devices", DeviceListHandler)
	r.HandleFunc("/devices/{id:[0-9]+}/effects", EffectListHandler)

	if err = http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
