package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"lamp/devicemaster"
	"lamp/effect"
	_ "lamp/effect/fire"
	_ "lamp/effect/wheel"
	_ "lamp/effect/wheel2"
	"lamp/lampbase"
	"log"
	"net"
	"net/http"
)

var (
	lampAddress       string
	lampStripes       int
	lampLedsPerStripe int
	lampDelay         int
	listenAddr        string
	staticServeDir    string
	dm                *devicemaster.DeviceMaster
)

func init() {
	flag.StringVar(&listenAddr, "listen", ":8080", "Address the lampserver listens on")
	flag.StringVar(&staticServeDir, "serve", "./static", "Directory to serve static content from")
	flag.StringVar(&lampAddress, "lamp", "192.168.178.178:8888", "Address of the lamp")
	flag.IntVar(&lampStripes, "stripes", 4, "Number of stripes the lamp has")
	flag.IntVar(&lampLedsPerStripe, "leds", 26, "Number of LEDs per stripe")
	flag.IntVar(&lampDelay, "delay", 25, "Milliseconds between updates")
}

func DeviceListHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Requesting device list")
	w.Header().Set("Content-Type", "application/json")
	out, err := json.Marshal(dm.DeviceList())
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(out)
}

func DeviceHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	deviceId := vars["id"]
	device, ok := dm.Device(deviceId)
	if !ok {
		log.Println("Did not find", device)
		http.NotFound(w, req)
		return
	}
	out, err := json.Marshal(device)
	if err != nil {
		log.Println(err)
		http.Error(w, "darn fuck it", 302)
		return
	}
	w.Write(out)
}

func EffectGetHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	deviceId := vars["id"]
	device, ok := dm.Device(deviceId)
	if !ok {
		log.Println("Did not find", device)
		http.NotFound(w, req)
		return
	}
	out, err := json.Marshal(device.CurrentEffect)
	if err != nil {
		log.Println(err)
		http.Error(w, "darn fuck it", 302)
		return
	}
	w.Write(out)
}

type EffectPut struct {
	Name   string
	Config json.RawMessage
}

func EffectPutHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	deviceId := vars["id"]
	device, ok := dm.Device(deviceId)
	if !ok {
		log.Println("Did not find", device)
		http.NotFound(w, req)
		return
	}
	put := &EffectPut{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(put)
	if err != nil {
		log.Println(err)
		// TODO correct error code for malformed input
		http.Error(w, "darn fuck it", 302)
		return
	}
	config, ok := effect.DefaultRegistry.Config(put.Name)
	if !ok {
		log.Println("Did not find", device)
		http.NotFound(w, req)
		return
	}
	err = json.Unmarshal(put.Config, config)
	if err != nil {
		log.Println(err)
		// TODO correct error code for malformed input
		http.Error(w, "darn fuck it", 302)
		return
	}
	dm.SetEffect(deviceId, put.Name, config)
	w.Write([]byte{})
}

func EffectListHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	deviceId := vars["id"]
	device, ok := dm.Device(deviceId)
	if !ok {
		log.Println("Did not find", device)
		http.NotFound(w, req)
		return
	}
	effectList := effect.DefaultRegistry.CompatibleEffects(device.Device)
	out, err := json.Marshal(effectList)
	if err != nil {
		log.Println(err)
		http.Error(w, "darn fuck it", 302)
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

	dm = devicemaster.New(effect.DefaultRegistry)
	dm.AddDevice("Big Lamp", "big", lamp)
	r := mux.NewRouter()
	r.HandleFunc("/devices", DeviceListHandler)
	r.HandleFunc("/devices/{id}", DeviceHandler).Methods("GET")
	r.HandleFunc("/devices/{id}/effect", EffectGetHandler).Methods("GET")
	r.HandleFunc("/devices/{id}/effect", EffectPutHandler).Methods("PUT")
	r.HandleFunc("/devices/{id}/available", EffectListHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir(staticServeDir))))
	// Redirect toplevel requests to the static folder so browsers find index.html
	r.Path("/").Handler(http.RedirectHandler("/static/", 302))

	if err = http.ListenAndServe(listenAddr, r); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
