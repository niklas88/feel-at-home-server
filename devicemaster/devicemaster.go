package devicemaster

import (
	"errors"
	"lamp/effect"
	"lamp/lampbase"
	"log"
	"sync"
)

type DeviceInfoShort struct {
	Name string
	Id   string
}

type DeviceInfo struct {
	Name          string
	Id            string
	CurrentEffect *effect.Info    `json:"-"`
	Device        lampbase.Device `json:"-"`
	controller    *effect.Controller
}

type DeviceMaster struct {
	sync.RWMutex
	devices map[string]*DeviceInfo
	reg     effect.Registry
}

func New(registry effect.Registry) *DeviceMaster {
	return &DeviceMaster{devices: make(map[string]*DeviceInfo),
		reg: registry}
}

func (d *DeviceMaster) AddDevice(name, id string, dev lampbase.Device) {
	d.Lock()
	defer d.Unlock()
	if _, ok := d.devices[id]; ok {
		panic("Readded device " + id)
	}
	device := &DeviceInfo{Name: name,
		Id:            id,
		CurrentEffect: nil,
		Device:        dev,
		controller:    effect.NewController()}
	go device.controller.Run()
	d.devices[id] = device
}

func (d *DeviceMaster) SetEffect(deviceId, effectName string, config effect.Config) error {
	d.Lock()
	defer d.Unlock()
	dev, ok := d.devices[deviceId]
	if !ok {
		return errors.New("Unknown device " + deviceId)
	}

	eff, info := d.reg.CreateEffect(effectName, dev.Device)
	if eff == nil {
		return errors.New("Unknown or incompatible effect")
	}

	if c, ok := eff.(effect.Configurer); ok {
		log.Println("Configuring")
		c.Configure(config)
	}
	info.Config = config
	dev.CurrentEffect = info
	dev.controller.EffectChan <- eff
	return nil
}

func (d *DeviceMaster) DeviceList() []DeviceInfoShort {
	var devList []DeviceInfoShort
	d.RLock()
	defer d.RUnlock()
	for _, v := range d.devices {
		devList = append(devList, DeviceInfoShort{Name: v.Name, Id: v.Id})
	}
	return devList
}

func (d *DeviceMaster) Device(id string) (*DeviceInfo, bool) {
	d.RLock()
	defer d.RUnlock()
	dev, ok := d.devices[id]
	if !ok {
		return &DeviceInfo{}, ok
	}
	return dev, true
}
