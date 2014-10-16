package devicemaster

import (
	"errors"
	"lamp/effect"
	"lamp/lampbase"
	"log"
	"sync"
)

type DeviceInfoShort struct {
	Name   string
	Id     string
	Active bool
}

type DeviceInfo struct {
	Name          string
	Id            string
	Active        bool
	CurrentEffect *effect.Info    `json:"-"`
	Device        lampbase.Device `json:"-"`
	controller    *effect.Controller
}

type DeviceMaster struct {
	sync.RWMutex
	deviceMap map[string]*DeviceInfo
	devices   []*DeviceInfo
	reg       effect.Registry
}

func New(registry effect.Registry) *DeviceMaster {
	return &DeviceMaster{deviceMap: make(map[string]*DeviceInfo),
		devices: make([]*DeviceInfo, 0),
		reg:     registry}
}

func (d *DeviceMaster) AddDevice(name, id string, dev lampbase.Device) {
	d.Lock()
	defer d.Unlock()
	if _, ok := d.deviceMap[id]; ok {
		panic("Readded device " + id)
	}
	newDeviceInfo := &DeviceInfo{Name: name,
		Id:            id,
		CurrentEffect: nil,
		Device:        dev,
		controller:    effect.NewController(dev)}
	d.devices = append(d.devices, newDeviceInfo)

	go newDeviceInfo.controller.Run()
	d.deviceMap[id] = newDeviceInfo
}

func (d *DeviceMaster) SetEffect(deviceId, effectName string, config effect.Config) error {
	d.Lock()
	defer d.Unlock()
	dev, ok := d.deviceMap[deviceId]
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
	dev.Active = true
	return nil
}

func (d *DeviceMaster) SetActive(deviceId string, active bool) error {
	d.Lock()
	defer d.Unlock()

	dev, ok := d.deviceMap[deviceId]
	if !ok {
		return errors.New("Unknown device " + deviceId)
	}
	if dev.Active == active {
		return nil
	}

	if active {
		dev.controller.StateChange <- effect.Activate
	} else {
		dev.controller.StateChange <- effect.Deactivate
	}
	dev.Active = active
	return nil
}

func (d *DeviceMaster) DeviceList() []DeviceInfoShort {
	// Copy for concurrency safety
	var devList []DeviceInfoShort
	d.RLock()
	defer d.RUnlock()
	for _, v := range d.devices {
		devList = append(devList, DeviceInfoShort{Name: v.Name, Id: v.Id, Active: v.Active})
	}
	return devList
}

func (d *DeviceMaster) Device(id string) (DeviceInfo, bool) {
	d.RLock()
	defer d.RUnlock()
	dev, ok := d.deviceMap[id]
	if !ok {
		return DeviceInfo{}, ok
	}
	// Make copy so that dev.CurrentEffect can't change concurrently
	return *dev, true
}
