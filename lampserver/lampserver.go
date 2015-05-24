package main

import (
	"encoding/json"
	"flag"
	"github.com/coreos/go-systemd/activation"
	"github.com/gorilla/mux"
	"lamp/devicemaster"
	"lamp/effect"
	_ "lamp/effect/brightness"
	_ "lamp/effect/colorfade"
	_ "lamp/effect/power"
	_ "lamp/effect/rainbow"
	_ "lamp/effect/static"
	_ "lamp/effect/strobe"
	_ "lamp/effect/sunrise"
	_ "lamp/effect/wheel"
	_ "lamp/effect/whitefade"
	"lamp/lampbase"
	"log"
	"net"
	"net/http"
	"os"
)

var (
	lampDelay      int
	listenAddr     string
	staticServeDir string
	configFileName string
	serverConfig   ServerConfig
	dm             *devicemaster.DeviceMaster
)

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
		http.Error(w, "darn fuck it", 400)
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
		http.Error(w, "darn fuck it", 400)
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
		http.Error(w, "darn fuck it", 400)
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
			http.Error(w, "darn fuck it config broken", 400)
			return
		}
	}
	dm.SetEffect(deviceId, put.Name, config)
	w.Write([]byte{})
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
		http.Error(w, "darn fuck it", 400)
		return
	}
	dm.SetActive(deviceId, put.Active)
	w.Write([]byte{})
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
		http.Error(w, "darn fuck it", 400)
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
	case "udpstripelamp":
		lampStripes := config.ensureNumeric("lampStripes")
		lampLedsPerStripe := config.ensureNumeric("lampLedsPerStripe")

		udpStripeLamp := lampbase.NewUdpStripeLamp(int(lampStripes), int(lampLedsPerStripe))
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
