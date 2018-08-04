package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"lheinrich.de/secpass/conf"
	"lheinrich.de/secpass/handler"
	"lheinrich.de/secpass/shorts"
)

// main function
func main() {
	// setups
	setupConfig()
	setupLogging()
	setupDB()
	setupWebserver()
}

// setup config and languages
func setupConfig() {
	// load config
	conf.Config = conf.ReadConfig("resources/config.json")

	// define language directory and list files
	langDir := conf.Config["app"]["languageDirectory"]
	files, err := ioutil.ReadDir(langDir)
	shorts.Check(err, true)

	// loop files and load language
	for _, file := range files {
		lang := strings.Split(file.Name(), ".")
		conf.Lang[lang[0]] = conf.ReadLanguage(langDir + "/" + file.Name())
	}
}

// setup logging
func setupLogging() {
	// define logging configurations
	logFile := conf.Config["app"]["logFile"]
	logInfo := conf.Config["app"]["logInfo"] == "true"

	// initialize logging to file
	shorts.InitLoggingFile(time.Now().Format(logFile), logInfo)
}

// setup database
func setupDB() {
	// define config values
	host := conf.Config["postgresql"]["host"]
	port := conf.Config["postgresql"]["port"]
	database := conf.Config["postgresql"]["database"]
	username := conf.Config["postgresql"]["username"]
	password := conf.Config["postgresql"]["password"]
	ssl := conf.Config["postgresql"]["ssl"]

	// connect to postgresql database
	conf.DB = shorts.ConnectPostgreSQL(host, port, database, username, password, ssl)
}

// setup webserver
func setupWebserver() {
	// setup handlers
	setupHandlers()

	// define config values
	address := conf.Config["webserver"]["address"]
	cert := conf.Config["webserver"]["certificateFile"]
	key := conf.Config["webserver"]["keyFile"]

	// start http server
	http.ListenAndServeTLS(address, cert, key, nil)
}

// setup webserver handlers
func setupHandlers() {
	// load templates
	handler.LoadTemplates()

	// register handlers
	http.HandleFunc("/", handler.Index)
	http.HandleFunc("/css/", handler.CSS)
	http.HandleFunc("/register", handler.Register)
	http.HandleFunc("/login", handler.Login)
}
