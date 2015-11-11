package main

import (
	"encoding/json"
	"flag"
	"github.com/coreos/go-systemd/activation"
	"github.com/gorilla/mux"
	"github.com/niklas88/feel-at-home-server/devicemaster"
	"github.com/niklas88/feel-at-home-server/effect"
	_ "github.com/niklas88/feel-at-home-server/effect/brightness"
	_ "github.com/niklas88/feel-at-home-server/effect/brightnessscaling"
	_ "github.com/niklas88/feel-at-home-server/effect/clock"
	_ "github.com/niklas88/feel-at-home-server/effect/clockcolor"
	_ "github.com/niklas88/feel-at-home-server/effect/color"
	_ "github.com/niklas88/feel-at-home-server/effect/colorfade"
	_ "github.com/niklas88/feel-at-home-server/effect/heart"
	_ "github.com/niklas88/feel-at-home-server/effect/power"
	_ "github.com/niklas88/feel-at-home-server/effect/rainbow"
	_ "github.com/niklas88/feel-at-home-server/effect/random"
	_ "github.com/niklas88/feel-at-home-server/effect/strobe"
	_ "github.com/niklas88/feel-at-home-server/effect/sunrise"
	_ "github.com/niklas88/feel-at-home-server/effect/whitefade"
	"github.com/niklas88/feel-at-home-server/lampbase"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	lampDelay      int
	listenAddr     string
	staticServeDir string
	configFileName string
	serverConfig   ServerConfig
	dm             *devicemaster.DeviceMaster
)

type StatusResult struct {
	Status string
	Error  string
}

func writeStatusResult(w http.ResponseWriter, err error) {
	var resp []byte
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		resp, _ = json.Marshal(&StatusResult{Status: "error", Error: err.Error()})
	} else {
		resp, _ = json.Marshal(&StatusResult{Status: "success", Error: ""})
	}
	w.Write(resp)
}

func init() {
	flag.IntVar(&lampDelay, "delay", 25, "Milliseconds between updates")
	flag.StringVar(&configFileName, "configfilename", "config.json", "Filepath of the configfile")
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
		log.Println("Did not find", deviceId)
		http.NotFound(w, req)
		return
	}
	out, err := json.Marshal(device)
	if err != nil {
		log.Println(err)
		http.Error(w, "darn fuck it", http.StatusInternalServerError)
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
		log.Println("Did not find", deviceId)
		http.NotFound(w, req)
		return
	}
	out, err := json.Marshal(device.CurrentEffect)
	if err != nil {
		log.Println(err)
		http.Error(w, "darn fuck it", http.StatusInternalServerError)
		return
	}
	w.Write(out)
}

type EffectPut struct {
	Name   string
	Config json.RawMessage
}

type ServerConfig struct {
	ListenAddress string
	Devices       []map[string]interface{}
}

func EffectPutHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	deviceId := vars["id"]
	_, ok := dm.Device(deviceId)
	if !ok {
		log.Println("Did not find", deviceId)
		http.NotFound(w, req)
		return
	}
	put := &EffectPut{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(put)
	if err != nil {
		log.Println(err)
		http.Error(w, "darn fuck it", http.StatusInternalServerError)
		return
	}
	config := effect.DefaultRegistry.Config(put.Name)
	if config == nil {
		log.Println("Did not find", deviceId)
		http.NotFound(w, req)
		return
	}
	if config != nil {
		err = json.Unmarshal(put.Config, config)
		if err != nil {
			log.Println(err)
			http.Error(w, "darn fuck it config broken", http.StatusInternalServerError)
			return
		}
	}

	err = dm.SetEffect(deviceId, put.Name, config)
	writeStatusResult(w, err)
}

type ActivePut struct {
	Active bool
}

func ActivePutHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	deviceId := vars["id"]
	_, ok := dm.Device(deviceId)
	if !ok {
		log.Println("Did not find", deviceId)
		http.NotFound(w, req)
		return
	}
	put := &ActivePut{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(put)
	if err != nil {
		log.Println(err)
		http.Error(w, "darn fuck it", http.StatusInternalServerError)
		return
	}
	err = dm.SetActive(deviceId, put.Active)
	writeStatusResult(w, err)
}

func EffectListHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	deviceId := vars["id"]
	device, ok := dm.Device(deviceId)
	if !ok {
		log.Println("Did not find", deviceId)
		http.NotFound(w, req)
		return
	}
	effectList := effect.DefaultRegistry.CompatibleEffects(device.Device)
	out, err := json.Marshal(effectList)
	if err != nil {
		log.Println(err)
		http.Error(w, "darn fuck it", http.StatusInternalServerError)
		return
	}
	w.Write(out)
}

type configMap map[string]interface{}

func (config configMap) ensureDefined(name string) interface{} {
	val, ok := config[name]
	if !ok {
		log.Fatalf("Missing %s field in config", name)
	}
	return val
}

func (config configMap) ensureNonEmptyString(name string) (result string) {
	val := config.ensureDefined(name)
	result, ok := val.(string)
	if !ok || result == "" {
		log.Fatalf("Value \"%v\" of field %s is not a string or empty", val, name)
	}

	return result
}

func (config configMap) ensureNumeric(name string) (result float64) {
	val := config.ensureDefined(name)
	result, ok := val.(float64)
	if !ok {
		log.Fatalf("Value \"%v\" of field %s is not numeric", val, name)
	}

	return result
}

func deviceFromConfig(config configMap) (name string, id string, device lampbase.Device) {
	// TODO Refactor
	lampAddress := config.ensureNonEmptyString("lampAddress")
	id = config.ensureNonEmptyString("id")
	name = config.ensureNonEmptyString("name")
	typeName := config.ensureNonEmptyString("type")
	devicePort := config.ensureNumeric("devicePort")

	addr, err := net.ResolveUDPAddr("udp4", lampAddress)
	if err != nil {
		log.Fatal("Couldn't resolve", err)
	}

	switch typeName {
	case "udpwordclock":
		timeUpdateInterval, err := time.ParseDuration(config.ensureNonEmptyString("timeUpdateInterval"))
		if err != nil {
			timeUpdateInterval = 1 * time.Minute
		}
		udpWordClock := lampbase.NewUdpWordClock(timeUpdateInterval)
		if err = udpWordClock.Dial(nil, addr, uint8(devicePort)); err != nil {
			log.Fatal("Couldn't create UdpWordClock", err)
		}
		device = udpWordClock
		break
	case "udpmatrixlamp":
		udpMatrix := lampbase.NewUdpMatrixLamp()
		if err = udpMatrix.Dial(nil, addr, uint8(devicePort)); err != nil {
			log.Fatal("Couldn't create UdpMatrix", err)
		}
		device = udpMatrix
		break
	case "udpstripelamp":
		//lampStripes := config.ensureNumeric("lampStripes")
		//lampLedsPerStripe := config.ensureNumeric("lampLedsPerStripe")

		udpStripeLamp := lampbase.NewUdpStripeLamp()
		if err = udpStripeLamp.Dial(nil, addr, uint8(devicePort)); err != nil {
			log.Fatal("Couldn't create UdpStripeLamp", err)
		}
		device = udpStripeLamp
		break
	case "udpcolorlamp":
		udpAnalogColorLamp := lampbase.NewUdpColorLamp()
		if err = udpAnalogColorLamp.Dial(nil, addr, uint8(devicePort)); err != nil {
			log.Fatal("Couldn't create UdpAnalogColorLamp", err)
		}
		device = udpAnalogColorLamp
		break
	case "udpdimlamp":
		udpDimLamp := lampbase.NewUdpDimLamp()
		if err = udpDimLamp.Dial(nil, addr, uint8(devicePort)); err != nil {
			log.Fatal("Couldn't create UdpDimLamp", err)
		}
		device = udpDimLamp
		break
	case "udppowerdevice":
		udpPowerDevice := lampbase.NewUdpPowerDevice()
		if err = udpPowerDevice.Dial(nil, addr, uint8(devicePort)); err != nil {
			log.Fatal("Couldn't create UdpPowerDevice", err)
		}
		device = udpPowerDevice
		break
	}
	return name, id, device
}

func main() {
	flag.Parse()

	file, err := os.Open(configFileName)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&serverConfig)
	if err != nil {
		log.Fatal(err)
	}
	dm = devicemaster.New(effect.DefaultRegistry)
	for _, value := range serverConfig.Devices {
		name, id, lamp := deviceFromConfig(value)
		if name != "" && id != "" && lamp != nil {
			dm.AddDevice(name, id, lamp)
		}
	}

	r := mux.NewRouter()
	r.HandleFunc("/devices", DeviceListHandler)
	r.HandleFunc("/devices/{id}", DeviceHandler).Methods("GET")
	r.HandleFunc("/devices/{id}/effect", EffectGetHandler).Methods("GET")
	r.HandleFunc("/devices/{id}/effect", EffectPutHandler).Methods("PUT")
	r.HandleFunc("/devices/{id}/active", ActivePutHandler).Methods("PUT")
	r.HandleFunc("/devices/{id}/available", EffectListHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir(staticServeDir))))
	// Redirect toplevel requests to the static folder so browsers find index.html
	r.Path("/").Handler(http.RedirectHandler("/static/", 302))

	files := activation.Files(false)
	var l net.Listener
	if len(files) != 1 {
		l, err = net.Listen("tcp", serverConfig.ListenAddress)
	} else {
		l, err = net.FileListener(files[0])
	}

	if err != nil {
		log.Fatal("Could not create Listener: ", err)
	}

	if err = http.Serve(l, r); err != nil {
		log.Fatal("Serve: ", err)
	}
}
