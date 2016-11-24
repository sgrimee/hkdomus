package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/log"

	"github.com/sgrimee/godomus"
)

//"flag"
//"github.com/spf13/pflag"

var (
	cfgFile  string
	debug    bool
	domus    *godomus.Domus
	pin      string
	groupKey godomus.GroupKey
)

// validateGroupSet ensures a group with exported LD devices is set
func getConfigGroupKey() godomus.GroupKey {
	if !viper.IsSet("group") || (viper.GetInt("group") < 1) {
		log.Fatal("You must five a group (int), or set it in config file")
	}
	return godomus.NewGroupKey(viper.GetInt("group"))
}

func main() {
	//pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	//pflag.Parse()
	debug = true
	initConfig()
	initDomus()
	domusLogin()
	log.Verbose = false

	// Setup bridge
	info := accessory.Info{
		Name: viper.GetString("name"),
	}
	bridge := accessory.New(info, accessory.TypeBridge)

	var accessories []*accessory.Accessory

	// Get devices from LD exported group
	group, err := domus.GetGroup(getConfigGroupKey())
	if err != nil {
		log.Fatal("Could not get group %d\n", groupKey.Num())
	}

	for _, d := range group.Devices {
		dev, err := domus.GetDeviceState(d.Key)
		if err != nil {
			log.Fatal("Could not get state for device %s\n", d.Key)
		}
		fmt.Printf("Adding device: %+v\n", dev)

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
