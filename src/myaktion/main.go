package main

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/turngeek/myaktion-go-2023/src/myaktion/handler"
)

func init() {
	// init logger
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
	level, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.Info("Log level not specified, using default log level: INFO")
		log.SetLevel(log.InfoLevel)
		return
	}
	log.SetLevel(level)
}

func main() {
	log.Println("Starting My-Aktion API server")
	router := mux.NewRouter()
	router.HandleFunc("/health", handler.Health).Methods("GET")
	router.HandleFunc("/campaign", handler.CreateCampaign).Methods("POST")
	router.HandleFunc("/campaigns", handler.GetCampaigns).Methods("GET")
	router.HandleFunc("/campaigns/{id}", handler.GetCampaign).Methods("GET")
	router.HandleFunc("/campaigns/{id}", handler.UpdateCampaign).Methods("PUT")
	router.HandleFunc("/campaigns/{id}", handler.DeleteCampaign).Methods("DELETE")
	router.HandleFunc("/campaigns/{id}/donation", handler.AddDonation).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
