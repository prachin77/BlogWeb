package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	clienthandlers "github.com/prachin77/client/ClientHandlers"
)

func main() {
	r := mux.NewRouter()

	fmt.Println("listening to port :1234")

	r.HandleFunc("/", clienthandlers.DefaultRoute).Methods("GET")
	r.HandleFunc("/home", clienthandlers.RenderHomePage).Methods("GET")
	// r.HandleFunc("/home?userid={userid}",clienthandlers.RenderHomeWebPage).Methods("GET")

	// AUTH URLs 
	r.HandleFunc("/register",clienthandlers.RenderRegisterPage).Methods("GET")
	r.HandleFunc("/login",clienthandlers.RenderLoginPage).Methods("GET")
	r.HandleFunc("/login", clienthandlers.Login).Methods("POST")
	r.HandleFunc("/register",clienthandlers.Register).Methods("POST")

	log.Fatal(http.ListenAndServe(":1234", r))
}
