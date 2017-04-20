package main

import (
	"crypto/subtle"
	"encoding/json"
	"flag"
	"github.com/coreos/go-systemd/activation"
	"github.com/gorilla/mux"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
	_ "github.com/niklas88/feel-at-home-server/server/deviceapi/brightness"
	_ "github.com/niklas88/feel-at-home-server/server/deviceapi/brightnessscaling"
	_ "github.com/niklas88/feel-at-home-server/server/deviceapi/clock"
	_ "github.com/niklas88/feel-at-home-server/server/deviceapi/clockcolor"
	_ "github.com/niklas88/feel-at-home-server/server/deviceapi/color"
	_ "github.com/niklas88/feel-at-home-server/server/deviceapi/colorfade"
	_ "github.com/niklas88/feel-at-home-server/server/deviceapi/fire"
	_ "github.com/niklas88/feel-at-home-server/server/deviceapi/heart"
	_ "github.com/niklas88/feel-at-home-server/server/deviceapi/power"
	_ "github.com/niklas88/feel-at-home-server/server/deviceapi/rainbow"
	_ "github.com/niklas88/feel-at-home-server/server/deviceapi/random"
	_ "github.com/niklas88/feel-at-home-server/server/deviceapi/strobe"
	_ "github.com/niklas88/feel-at-home-server/server/deviceapi/sunrise"
	_ "github.com/niklas88/feel-at-home-server/server/deviceapi/wheel"
	_ "github.com/niklas88/feel-at-home-server/server/deviceapi/whitefade"
	"github.com/niklas88/feel-at-home-server/server/devicemaster"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	lampDelay      int
	listenAddr     string
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
	flag.StringVar(&configFileName, "config", "config.json", "Filepath of the configfile")
}

type BasicAuthConfig struct {
	Username string
	Password string
	Realm    string
}

// BasicAuth wraps a handler requiring HTTP basic auth for it using the given
// username and password and the specified realm, which shouldn't contain quotes.
//
// Most web browser display a dialog with something like:
//
//    The website says: "<realm>"
//
func BasicAuth(handler http.HandlerFunc, authconf *BasicAuthConfig) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(authconf.Username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(authconf.Password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+authconf.Realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		handler(w, r)
	}
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
	Username      string
	Password      string
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
	config := deviceapi.DefaultRegistry.Config(put.Name)
	if config == nil {
		log.Println("Did not find", deviceId)
		http.NotFound(w, req)
		return
	}
	err = json.Unmarshal(put.Config, config)
	if err != nil {
		log.Println(err)
		http.Error(w, "darn fuck it config broken", http.StatusInternalServerError)
		return
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
	effectList := deviceapi.DefaultRegistry.CompatibleEffects(device.Device)
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

func deviceFromConfig(config configMap) (name string, id string, device devices.Device) {
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
		udpWordClock := devices.NewUdpWordClock(timeUpdateInterval)
		if err = udpWordClock.Dial(nil, addr, uint8(devicePort)); err != nil {
			log.Fatal("Couldn't create UdpWordClock", err)
		}
		device = udpWordClock
		break
	case "udpmatrixlamp":
		udpMatrix := devices.NewUdpMatrixLamp()
		if err = udpMatrix.Dial(nil, addr, uint8(devicePort)); err != nil {
			log.Fatal("Couldn't create UdpMatrix", err)
		}
		device = udpMatrix
		break
	case "udpstripelamp":
		//lampStripes := config.ensureNumeric("lampStripes")
		//lampLedsPerStripe := config.ensureNumeric("lampLedsPerStripe")

		udpStripeLamp := devices.NewUdpStripeLamp()
		if err = udpStripeLamp.Dial(nil, addr, uint8(devicePort)); err != nil {
			log.Fatal("Couldn't create UdpStripeLamp", err)
		}
		device = udpStripeLamp
		break
	case "udpcolorlamp":
		udpAnalogColorLamp := devices.NewUdpColorLamp()
		if err = udpAnalogColorLamp.Dial(nil, addr, uint8(devicePort)); err != nil {
			log.Fatal("Couldn't create UdpAnalogColorLamp", err)
		}
		device = udpAnalogColorLamp
		break
	case "udpdimlamp":
		udpDimLamp := devices.NewUdpDimLamp()
		if err = udpDimLamp.Dial(nil, addr, uint8(devicePort)); err != nil {
			log.Fatal("Couldn't create UdpDimLamp", err)
		}
		device = udpDimLamp
		break
	case "udppowerdevice":
		udpPowerDevice := devices.NewUdpPowerDevice()
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
	dm = devicemaster.New(deviceapi.DefaultRegistry)
	for _, value := range serverConfig.Devices {
		name, id, lamp := deviceFromConfig(value)
		if name != "" && id != "" && lamp != nil {
			dm.AddDevice(name, id, lamp)
		}
	}

	r := mux.NewRouter()
	auth := &BasicAuthConfig{serverConfig.Username, serverConfig.Password, "Feel@Home"}

	r.HandleFunc("/devices", BasicAuth(DeviceListHandler, auth))
	r.HandleFunc("/devices/{id}", BasicAuth(DeviceHandler, auth)).Methods("GET")
	r.HandleFunc("/devices/{id}/effect", BasicAuth(EffectGetHandler, auth)).Methods("GET")
	r.HandleFunc("/devices/{id}/effect", BasicAuth(EffectPutHandler, auth)).Methods("PUT")
	r.HandleFunc("/devices/{id}/active", BasicAuth(ActivePutHandler, auth)).Methods("PUT")
	r.HandleFunc("/devices/{id}/available", BasicAuth(EffectListHandler, auth))

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
