package main

import (
	"encoding/json"
	"flag"
	"github.com/coreos/go-systemd/activation"
	"github.com/gorilla/mux"
	"lamp/devicemaster"
	"lamp/effect"
	_ "lamp/effect/brightness"
	_ "lamp/effect/fire"
	_ "lamp/effect/power"
	_ "lamp/effect/static"
	_ "lamp/effect/sunrise"
	_ "lamp/effect/wheel"
	_ "lamp/effect/wheel2"
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
	/*flag.StringVar(&listenAddr, "listen", ":8080", "Address the lampserver listens on")
	flag.StringVar(&staticServeDir, "serve", "./static", "Directory to serve static content from")
	flag.StringVar(&lampAddress, "lamp", "192.168.178.178:8888", "Address of the lamp")
	flag.IntVar(&lampStripes, "stripes", 4, "Number of stripes the lamp has")
	flag.IntVar(&lampLedsPerStripe, "leds", 26, "Number of LEDs per stripe")*/
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
		log.Println("Did not find", device)
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
		log.Println("Did not find", device)
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
		http.Error(w, "darn fuck it", 400)
		return
	}
	config, ok := effect.DefaultRegistry.Config(put.Name)
	if !ok {
		log.Println("Did not find", device)
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
		http.Error(w, "darn fuck it", 400)
		return
	}
	w.Write(out)
}

func deviceFromConfig(configDevice map[string]interface{}) (string, string, lampbase.Device) {
	lampAddress, okLampAddress := configDevice["lampAddress"]
	id, okID := configDevice["id"]
	name, okName := configDevice["name"]
	t, okTyp := configDevice["typ"]
	if okLampAddress && okID && okName && okTyp {
		lampAddress, okLampAddress := lampAddress.(string)
		id, okID := id.(string)
		name, okName := name.(string)
		t, okTyp := t.(string)

		if okLampAddress && okID && okName && okTyp && lampAddress != "" && id != "" && name != "" && t != "" {
			addr, err := net.ResolveUDPAddr("udp4", lampAddress)
			if err != nil {
				log.Println("bam5")
				log.Fatal("Couldn't resolve", err)
				return "", "", nil
			}
			var device lampbase.Device
			switch t {
			case "udpstripelamp":
				lampStripesInterface, okLampStripes := configDevice["lampStripes"]
				lampLedsPerStripeInterface, okLedsPerStripe := configDevice["lampLedsPerStripe"]
				if okLampStripes && okLedsPerStripe {
					lampStripes, okLampStripes := lampStripesInterface.(float64)
					lampLedsPerStripe, okLedsPerStripe := lampLedsPerStripeInterface.(float64)
					if okLampStripes && okLedsPerStripe {
						udpStripeLamp := lampbase.NewUdpStripeLamp(int(lampStripes), int(lampLedsPerStripe))
						if err = udpStripeLamp.Dial(nil, addr); err != nil {
							log.Println(err)
							return "", "", nil
						}
						device = udpStripeLamp
					} else {
						return "", "", nil
					}
				} else {
					return "", "", nil
				}

				break
			case "udpanalogcolorlamp":
				udpAnalogColorLamp := lampbase.NewUdpAnalogColorLamp()
				if err = udpAnalogColorLamp.Dial(nil, addr); err != nil {
					log.Println(err)
					return "", "", nil
				}
				device = udpAnalogColorLamp
				break
			case "udpdimlamp":
				udpDimLamp := lampbase.NewUdpDimLamp()
				if err = udpDimLamp.Dial(nil, addr); err != nil {
					log.Println(err)
					return "", "", nil
				}
				device = udpDimLamp
				break
			case "udppowerdevice":
				udpPowerDevice := lampbase.NewUdpPowerDevice()
				if err = udpPowerDevice.Dial(nil, addr); err != nil {
					log.Println(err)
					return "", "", nil
				}
				device = udpPowerDevice
				break
			}
			return name, id, device
		}
	}
	return "", "", nil

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
	if err = http.Serve(l, r); err != nil {
		log.Fatal("Serve: ", err)
	}
}
