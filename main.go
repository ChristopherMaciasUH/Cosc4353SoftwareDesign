package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"quickstart/controller"
)


func main(){
	fmt.Println("Starting...")
	r := mux.NewRouter()
	r.HandleFunc("/register", controller.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", controller.LoginHandler).Methods("POST")
	r.HandleFunc("/profileInfo",controller.GetProfileInfo).Methods("GET")
	r.HandleFunc("/profileSetter",controller.InsertProfileInfo).Methods("POST")
	r.HandleFunc("/fuelQuoteForm", controller.DeliveryRequestHandler).Methods("POST")
	r.HandleFunc("/fuelQuoteHistory", controller.GetDeliveryRequests).Methods("GET")
	fmt.Println("nice!")
	http.ListenAndServe(":8000",r)
}