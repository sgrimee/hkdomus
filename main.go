package main

import (
	"fmt"
	"log"

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
	initConfig()
	initDomus()
	domusLogin()

	devices, err := domus.DevicesInRoom(
		// just lamps in a given room for now
		godomus.NewRoomKey(58), godomus.CategoryClassId("CLSID-DEVC-A-EC"))
	if err != nil {
		log.Fatalf("Could not get devices from LD: %s\n", err)
	}

	for _, dev := range devices {
		go AddDevice(dev)
	}

	select {}
}

func AddDevice(dev godomus.Device) {
	switchInfo := accessory.Info{
		Name: dev.Label,
	}
	acc := accessory.NewSwitch(switchInfo)

	config := hc.Config{Pin: pin, StoragePath: fmt.Sprintf("$HOME/.hkdomus/db/%s", dev.Key)}
	t, err := hc.NewIPTransport(config, acc.Accessory)

	if err != nil {
		log.Fatal(err)
	}

	// Log to console when client (e.g. iOS app) changes the value of the on characteristic
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
		err := domus.ExecuteAction(action, property, dev.Key)
		if err != nil {
			log.Print(err)
		}
	})

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
