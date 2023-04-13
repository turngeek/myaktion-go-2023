package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/turngeek/myaktion-go-2023/src/myaktion/handler"
)

func main() {
	log.Println("Starting My-Aktion API server")
	router := mux.NewRouter()
	router.HandleFunc("/health", handler.Health).Methods("GET")
	router.HandleFunc("/campaign", handler.CreateCampaign).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
