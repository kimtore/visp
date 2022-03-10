package spotify_devices

import (
	"fmt"

	"github.com/ambientsound/visp/list"
	"github.com/ambientsound/visp/log"
	"github.com/ambientsound/visp/utils"
	"github.com/zmb3/spotify"
)

type List struct {
	list.Base
	devices map[string]spotify.PlayerDevice
}

var _ list.List = &List{}

func New(client spotify.Client) (*List, error) {
	var err error

	devices, err := client.PlayerDevices()

	if err != nil {
		return nil, err
	}

	this := &List{
		devices: make(map[string]spotify.PlayerDevice),
	}
	this.Clear()

	for _, device := range devices {
		if len(device.ID) == 0 {
			log.Debugf("ignoring encountered device with empty ID: %+v", device)
			continue
		}
		this.devices[device.ID.String()] = device
		this.Add(Row(device))
	}

	this.SetID("devices")
	this.SetName("Player devices")
	this.SetVisibleColumns([]string{
		"deviceName",
		"deviceType",
		"active",
		"restricted",
		"volume",
	})

	return this, nil
}

func Row(device spotify.PlayerDevice) list.Row {
	return list.NewRow(
		device.ID.String(),
		list.DataTypeDevice,
		map[string]string{
			"deviceName": device.Name,
			"deviceType": device.Type,
			"active":     utils.HumanFormatBool(device.Active),
			"restricted": utils.HumanFormatBool(device.Restricted),
			"volume":     fmt.Sprintf("%d%%", device.Volume),
		},
	)
}

// CursorDevice returns the device currently selected by the cursor.
func (s *List) CursorDevice() *spotify.PlayerDevice {
	return s.Device(s.Cursor())
}

// Device returns the device at a specific index.
func (s *List) Device(index int) *spotify.PlayerDevice {
	row := s.Row(index)
	if row == nil {
		return nil
	}
	device := s.devices[row.ID()]
	return &device
}

// Device returns the device with a specific ID.
func (s *List) DeviceByID(id string) *spotify.PlayerDevice {
	device, ok := s.devices[id]
	if !ok {
		return nil
	}
	return &device
}
