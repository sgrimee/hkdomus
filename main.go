package main

import (
	"fmt"
	"os"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/log"

	"github.com/sgrimee/godomus"
)

var (
	cfgFile string
	config  Config
	domus   *godomus.Domus
)

func main() {
	config = getConfig()
	domus, err := godomus.New(config.DomusConfig)
	if err != nil {
		log.Fatal(err)
	}

	// disable logs of hc package
	log.Verbose = false

	// Setup bridge
	info := accessory.Info{
		Name: config.BridgeName,
	}
	bridge := accessory.New(info, accessory.TypeBridge)

	var accessories []*accessory.Accessory
	accForDev := make(map[godomus.DeviceKey]*accessory.Accessory)

	// Get devices from LD exported group
	group, err := domus.GetGroup(config.GroupKey)
	if err != nil {
		log.Fatal("Could not get group %d\n", config.GroupKey.Num())
	}

	for _, d := range group.Devices {
		dev, err := domus.GetDeviceState(d.Key)
		if err != nil {
			log.Fatal("Could not get state for device %s\n", d.Key)
		}

		switchInfo := accessory.Info{
			// adding the room to the name to make it possible to recognise the device
			// in the HomeKit apps at only the name is shown
			Name: fmt.Sprintf("%s (%s)", dev.Label, dev.RoomLabel),
		}
		acc := accessory.NewSwitch(switchInfo)

		acc.Switch.On.OnValueRemoteUpdate(func(on bool) {
			log.Printf("dev: %+v\n", dev)
			var err error
			log.Printf("[INFO] Client changing switch %s to %t\n", switchInfo.Name, on)
			if on {
				err = dev.On()
			} else {
				err = dev.Off()
			}
			if err != nil {
				log.Println(err)
			}
		})

		fmt.Printf("Adding device: %s\n", switchInfo.Name)
		accessories = append(accessories, acc.Accessory)
		accForDev[dev.Key] = acc.Accessory
	}

	hcConfig := hc.Config{Pin: config.Pin, StoragePath: fmt.Sprintf("%s/.hkdomus/db", os.Getenv("HOME"))}
	t, err := hc.NewIPTransport(hcConfig, bridge, accessories...)
	if err != nil {
		log.Fatal(err)
	}

	hc.OnTermination(func() {
		t.Stop()
	})

	go t.Start()

	// listen for device events from the LD server and update accessory status as needed
	devices := make(chan godomus.Device, 1)
	errs := make(chan error, 1)
	done := make(chan struct{})
	go domus.ListenForDeviceUpdates(devices, errs, done)

	for {
		select {
		case d := <-devices:
			fmt.Printf("got LD device update %s (%s): %s\n", d.Label, d.RoomLabel, d.States[0])
			if acc, ok := accForDev[d.Key]; ok {
				fmt.Printf("Found accessory: %s\n", acc.Info.Name)
			}
		case e := <-errs:
			fmt.Printf("got LD error: %s\n", e)
		case <-done:
			break
		}
	}
	t.Stop()
}
