package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"lamp/effects"
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
	reg               map[string]*effects.EffectInfo
)

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
		w.WriteHeader(404)
		return
	}
	device := deviceList[deviceId]
	var effectList []EffectDescription
	for effectName, effectInfo := range reg {
		if effectInfo.Compatible(device.Device) {
			effectList = append(effectList, EffectDescription{effectName, effectInfo.ConfigFactory()})
		}
	}

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
