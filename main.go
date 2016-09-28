package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"

	"github.com/sgrimee/godomus"
)

//"flag"
//"github.com/spf13/pflag"

var (
	cfgFile string
	debug   bool
	domus   *godomus.Domus
	devices godomus.Devices
	pin     string
)

func main() {
	//pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	//pflag.Parse()
	debug = true
	//log.Verbose = false
	initConfig()
	initDomus()
	domusLogin()

	// Get devices from LD
	for _, devNum := range []int{207, 208} {
		dev, err := domus.GetDeviceState(godomus.NewDeviceKey(devNum))
		if err != nil {
			log.Fatalf("Could not get device %d: %s\n", devNum, err)
		}
		devices = append(devices, *dev)
	}

	// Setup bridge
	info := accessory.Info{
		Name: "LifeDomus",
	}
	bridge := accessory.New(info, accessory.TypeBridge)

	var accessories []*accessory.Accessory

	for _, dev := range devices {
		switchInfo := accessory.Info{
			Name: dev.Label,
		}
		fmt.Printf("Adding device: %+v\n", dev)
		acc := accessory.NewSwitch(switchInfo)

		dk := dev.Key
		acc.Switch.On.OnValueRemoteUpdate(func(on bool) {
			property := godomus.PropClassId("CLSID-DEVC-PROP-TOR-SW")
			var action godomus.ActionClassId

			if on == true {
				log.Printf("[INFO] Client changed switch %s to on\n", acc.Info.Name)
				action = godomus.ActionClassId("CLSID-ACTION-ON")
			} else {
				log.Printf("[INFO] Client changed switch %s to off\n", acc.Info.Name)
				action = godomus.ActionClassId("CLSID-ACTION-OFF")
			}
			err := domus.ExecuteAction(action, property, dk)
			if err != nil {
				log.Println(err)
			}
		})

		accessories = append(accessories, acc.Accessory)
	}

	config := hc.Config{Pin: pin, StoragePath: fmt.Sprintf("%s/.hkdomus/db", os.Getenv("HOME"))}
	t, err := hc.NewIPTransport(config, bridge, accessories...)

	if err != nil {
		log.Fatal(err)
	}

	hc.OnTermination(func() {
		t.Stop()
	})

	t.Start()
}

// initConfig reads in config file.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".ldclient") // name of config file (without extension)
	viper.AddConfigPath("$HOME")     // adding home directory as first search path

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && debug {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	pin = viper.GetString("pin")

}
