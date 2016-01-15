// Package devicemaster implements a central hub for handling different
// devices, maintaining metadata of devices and keeping track of currently
// running deviceapis
package devicemaster

import (
	"errors"
	"fmt"
	"github.com/niklas88/feel-at-home-server/devices"
	"github.com/niklas88/feel-at-home-server/server/deviceapi"
	"github.com/niklas88/feel-at-home-server/server/deviceapi/power"
	"sync"
)

// DeviceInfoShort functions as POD type for storing the most important
// information on a device such as its name, id and whether it's currently
// active i.e. running a static or dynamic effect other than being Power()'ed
// off
type DeviceInfoShort struct {
	Name   string
	Id     string
	Active bool
}

// DeviceInfo holds all information maintained for a device under control
type DeviceInfo struct {
	Name          string
	Id            string
	Active        bool
	Config        deviceapi.Config `json:"-"`
	CurrentEffect deviceapi.Effect `json:"-"`
	Device        devices.Device   `json:"-"`
}

// DeviceMaster is the main type of this package through its methods allows
// controlling devices under its control as well as putting devices under its
// control
type DeviceMaster struct {
	sync.RWMutex
	deviceMap map[string]*DeviceInfo
	devices   []*DeviceInfo
	reg       *deviceapi.Registry
}

// New creates a new DeviceMaster instance using the provided deviceapi.Registry
// which maintains available effects and their metadata
func New(registry *deviceapi.Registry) *DeviceMaster {
	return &DeviceMaster{deviceMap: make(map[string]*DeviceInfo),
		devices: make([]*DeviceInfo, 0),
		reg:     registry}
}

// AddDevice puts a device under the control of this DeviceMaster instance
// registering it with a name and id. Readding an already added device results
// in a panic to prevent misuse
func (d *DeviceMaster) AddDevice(name, id string, dev devices.Device) {
	d.Lock()
	defer d.Unlock()
	if _, ok := d.deviceMap[id]; ok {
		panic("Readded device " + id)
	}

	powerEffect := d.reg.Effect("Power")

	newDeviceInfo := &DeviceInfo{Name: name,
		Id:            id,
		CurrentEffect: powerEffect,
		Config:        &power.PowerConfig{false},
		Active:        false,
		Device:        dev}
	d.devices = append(d.devices, newDeviceInfo)
	d.deviceMap[id] = newDeviceInfo
}

// SetEffect makes the effect given by effectName active for the device given
// by deviceId using the provided deviceapi.Config
func (d *DeviceMaster) SetEffect(deviceId, effectName string, config deviceapi.Config) error {
	d.Lock()
	defer d.Unlock()
	dev, ok := d.deviceMap[deviceId]
	if !ok {
		return errors.New("Unknown device " + deviceId)
	}

	eff := d.reg.Effect(effectName)
	if eff == nil {
		return fmt.Errorf("Incompatible effect %v for lamp type %T", effectName, dev.Device)
	}
	err := eff.Apply(dev.Device, config)
	if err != nil {
		return err
	}

	dev.CurrentEffect = eff
	dev.Config = config
	dev.Active = true
	return nil
}

// SetActive activates (active == true) or suspends (active == false) the
// current effect running on the device given by deviceId. If a device is
// already in the state being requested this is a no-op
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
		eff := d.reg.Effect(dev.CurrentEffect.Name())
		err = eff.Apply(dev.Device, dev.Config)
	} else {
		err = dev.Device.Power(active)
	}
	dev.Active = active
	return err
}

// DeviceList returns a list of all devices under the control of this
// DeviceMaster containing the most important metadata for each device
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

// Device returns detailed information such as the currently active effect for
// the device given by id
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
