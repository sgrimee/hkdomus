package main

import (
	"fmt"
	"log"

	"github.com/sgrimee/godomus"
	"github.com/spf13/viper"
)

type Config struct {
	Pin         string           // HomeKit pairing pin
	GroupKey    godomus.GroupKey // group containing devices to use
	DomusConfig godomus.Config
	// optional
	Debug      bool
	BridgeName string
}

// getConfig reads in config file.
func getConfig() Config {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".ldclient") // name of config file (without extension)
	viper.AddConfigPath("$HOME")     // adding home directory as first search path

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	viper.SetDefault("debug", "false")
	viper.SetDefault("name", "HKDomus bridge")

	gdCfg := godomus.Config{
		SiteKey:    cfgSiteKey(),
		UserKey:    cfgUserKey(),
		Password:   cfgPassword(),
		Url:        cfgUrl(),
		SocketPort: cfgSocketPort(),
	}

	return Config{
		Pin:         viper.GetString("pin"),
		GroupKey:    cfgGroupKey(),
		DomusConfig: gdCfg,
		Debug:       viper.GetBool("debug"),
		BridgeName:  viper.GetString("name"),
	}
}

func cfgGroupKey() godomus.GroupKey {
	if !viper.IsSet("group") || (viper.GetInt("group") < 1) {
		log.Fatal("You must five a group (int), or set it in config file")
	}
	return godomus.NewGroupKey(viper.GetInt("group"))
}

func cfgSiteKey() godomus.SiteKey {
	if !viper.IsSet("site") || (viper.GetInt("site") < 1) {
		log.Fatal("You must give a site (int), or set it in config file")
	}
	return godomus.NewSiteKey(viper.GetInt("site"))
}

func cfgSocketPort() int {
	if !viper.IsSet("socket_port") || (viper.GetInt("socket_port") < 1) {
		log.Fatal("You must give a socket_port (int), or set it in config file")
	}
	return viper.GetInt("socket_port")
}

func cfgUserKey() godomus.UserKey {
	if !viper.IsSet("user") || (viper.GetInt("user") < 1) {
		log.Fatal("You must give a user (int), or set it in config file")
	}
	return godomus.NewUserKey(viper.GetInt("user"))
}

func cfgUrl() string {
	if !viper.IsSet("url") || (viper.GetString("url") == "") {
		log.Fatal("You must give a url, or set it in config file")
	}
	return viper.GetString("url")
}

func cfgPassword() string {
	if !viper.IsSet("password") || (viper.GetString("password") == "") {
		log.Fatal("You must give a password, or set it in config file")
	}
	return viper.GetString("password")
}
