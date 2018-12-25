package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/MZIchenjl/ooi4/conf"
	"github.com/MZIchenjl/ooi4/handlers"
	"github.com/gorilla/mux"
)

const waitTime = time.Second * 15

func main() {
	// Read the config
	confile := flag.String("config", "app.toml", "Set the config file(toml)")
	flag.Parse()

	appConfig := new(conf.Config)
	_, err := toml.DecodeFile(*confile, appConfig)
	if err != nil {
		log.Fatalln(err)
	}
	// Set the http.DefaultClient's proxy
	http.DefaultClient.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment}

	// Init all handlers
	handlers.Init(appConfig.Secret, appConfig.Cookie)

	// Register routers
	r := mux.NewRouter()

	r.Methods(http.MethodGet).Path("/").HandlerFunc(handlers.Form)
	r.Methods(http.MethodPost).Path("/").HandlerFunc(handlers.Login)
	r.Methods(http.MethodGet).Path("/kancolle").HandlerFunc(handlers.Normal)
	r.Methods(http.MethodGet).Path("/kcv").HandlerFunc(handlers.KCV)
	r.Methods(http.MethodGet).Path("/flash").HandlerFunc(handlers.Flash)
	r.Methods(http.MethodGet).Path("/poi").HandlerFunc(handlers.Poi)
	r.Methods(http.MethodGet).Path("/connector").HandlerFunc(handlers.Connector)
	r.Methods(http.MethodGet).Path("/logout").HandlerFunc(handlers.Logout)

	r.Methods(http.MethodGet, http.MethodPost).Path("/kcsapi/{action:.+}").HandlerFunc(handlers.API)
	r.Methods(http.MethodGet).Path("/kcs2/resources/world/{server:.+}_{size:[lst]}.png").HandlerFunc(handlers.WorldImage)

	r.Methods(http.MethodPost).Path("/service/osapi").HandlerFunc(handlers.GetOSAPI)
	r.Methods(http.MethodPost).Path("/service/flash").HandlerFunc(handlers.GetFlash)

	r.Methods(http.MethodGet).PathPrefix("/kcs").HandlerFunc(handlers.Proxy)
	r.Methods(http.MethodGet).PathPrefix("/_kcs").HandlerFunc(handlers.Proxy)
	r.Methods(http.MethodGet).PathPrefix("/kcs2").HandlerFunc(handlers.Proxy)
	r.Methods(http.MethodGet).PathPrefix("/_kcs2").HandlerFunc(handlers.Proxy)

	r.Methods(http.MethodGet).PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Create the server
	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", appConfig.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// Start the server by using a go routine
	go func() {
		log.Printf("ooi serving on http://localhost:%d\n", appConfig.Port)
		if err := srv.ListenAndServe(); err != nil {
			// If the server can't start up, exit with code 1
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)

	// Wait for the system's Interrupt/Kill signal
	signal.Notify(c, os.Interrupt, os.Kill)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), waitTime)
	defer cancel()
	// Close the server and the context
	srv.Shutdown(ctx)

	log.Println("ooi is shutting down")

	// Exit the program with code 0
	os.Exit(0)
}
