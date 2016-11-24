package main

import (
	"log"

	"github.com/sgrimee/godomus"
	"github.com/spf13/viper"
)

// initDomus sets-up the domus object
func initDomus() {
	var err error
	domus, err = godomus.New(
		viper.GetString("url"),
		viper.GetInt("socket_port"),
	)
	if err != nil {
		log.Fatal(err)
	}
	if debug {
		domus.Debug = true
	}
}

// validateSiteSet ensures a site number was provided
func validateSiteSet() {
	if !viper.IsSet("site") || (viper.GetInt("site") < 1) {
		log.Fatal("You must give a site (int), or set it in config file")
	}
}

// validateUserSet ensures a userid and password were provided
func validateUserSet() {
	if !viper.IsSet("user") || (viper.GetInt("user") < 1) {
		log.Fatal("You must give a user (int), or set it in config file")
	}
	if !viper.IsSet("password") || (viper.GetString("password") == "") {
		log.Fatal("You must give a password, or set it in config file")
	}
}

// domusInfos logs in an returns infos
func domusInfos() godomus.LoginInfos {
	validateSiteSet()
	validateUserSet()
	sk := godomus.NewSiteKey(viper.GetInt("site"))
	uk := godomus.NewUserKey(viper.GetInt("user"))
	pass := viper.GetString("password")
	infos, err := domus.LoginInfos(sk, uk, pass)
	if err != nil {
		log.Fatal(err)
	}
	return infos
}

// domusLogin performs a login and stores the session key
func domusLogin() {
	validateSiteSet()
	validateUserSet()
	sk := godomus.NewSiteKey(viper.GetInt("site"))
	uk := godomus.NewUserKey(viper.GetInt("user"))
	pass := viper.GetString("password")
	_, err := domus.Login(sk, uk, pass)
	if err != nil {
		log.Fatal(err)
	}
}
