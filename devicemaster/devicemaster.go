package devicemaster

import (
	"errors"
	"fmt"
	"lamp/effect"
	"lamp/lampbase"
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

	powerInfo := d.reg.Info("Power")

	newDeviceInfo := &DeviceInfo{Name: name,
		Id:            id,
		CurrentEffect: powerInfo,
		Active:        false,
		Device:        dev}
	d.devices = append(d.devices, newDeviceInfo)
	d.deviceMap[id] = newDeviceInfo
}

func (d *DeviceMaster) SetEffect(deviceId, effectName string, config effect.Config) error {
	d.Lock()
	defer d.Unlock()
	dev, ok := d.deviceMap[deviceId]
	if !ok {
		return errors.New("Unknown device " + deviceId)
	}

	eff := d.reg.Effect(effectName, dev.Device)
	if eff == nil {
		return fmt.Errorf("Incompatible effect %v for lamp type %v", effectName, dev.Device)
	}
	err := eff.Apply(config)
	if err != nil {
		return err
	}

	info := d.reg.Info(effectName)
	info.Config = config
	dev.CurrentEffect = info
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
	var err error

	if active {
		eff := d.reg.Effect(dev.CurrentEffect.Name, dev.Device)
		err = eff.Apply(dev.CurrentEffect.Config)
	} else {
		err = dev.Device.Power(active)
	}
	dev.Active = active
	return err
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
