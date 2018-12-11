package main

import (
	"context"
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

const wait = time.Second * 15

func main() {
	appConfig := new(conf.Config)
	_, err := toml.DecodeFile("app.toml", appConfig)
	if err != nil {
		log.Fatalln(err)
	}

	api := &handlers.APIHandler{}
	f2e := &handlers.FrontEndHandler{}
	ser := &handlers.ServiceHandler{}

	api.Init(appConfig.Secret, appConfig.Cookie)
	f2e.Init(appConfig.Secret, appConfig.Cookie)
	ser.Init(appConfig.Secret, appConfig.Cookie)

	r := mux.NewRouter()

	r.Methods(http.MethodGet).Path("/").HandlerFunc(f2e.Form)
	r.Methods(http.MethodPost).Path("/").HandlerFunc(f2e.Login)
	r.Methods(http.MethodGet).Path("/kancolle").HandlerFunc(f2e.Normal)
	r.Methods(http.MethodGet).Path("/kcv").HandlerFunc(f2e.KCV)
	r.Methods(http.MethodGet).Path("/poi").HandlerFunc(f2e.Poi)
	r.Methods(http.MethodGet).Path("/connector").HandlerFunc(f2e.Connector)
	r.Methods(http.MethodGet).Path("/logout").HandlerFunc(f2e.Logout)
	r.Methods(http.MethodGet, http.MethodPost).Path("/kcsapi/{action:.+}").HandlerFunc(api.API)
	r.Methods(http.MethodGet).Path("/kcs/resources/image/world/{server:.+}_{size:[lst]}.png").HandlerFunc(api.WorldImage)
	r.Methods(http.MethodPost).Path("/service/osapi").HandlerFunc(ser.GetOSAPI)
	r.Methods(http.MethodPost).Path("/service/flash").HandlerFunc(ser.GetFlash)
	r.Methods(http.MethodGet).PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", appConfig.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		log.Printf("OOI serving on http://localhost:%d\n", appConfig.Port)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srv.Shutdown(ctx)

	log.Println("OOI is shutting down")
	os.Exit(0)
}
