package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	serverhandlers "github.com/prachin77/server/ServerHandlers"
)

func main() {
	r := mux.NewRouter()

	fmt.Println("listening to port :8080")

	r.HandleFunc("/register",serverhandlers.Register).Methods("POST")
	r.HandleFunc("/login",serverhandlers.Login).Methods("POST")
	r.HandleFunc("/login",serverhandlers.Logout).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080" ,	 r))
}
