package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	http.FileServer(http.Dir("static/"))

	r := mux.NewRouter()
	r.HandleFunc("/data", DataHandler)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Data struct {
	Id    int
	Value string
}

func DataHandler(w http.ResponseWriter, r *http.Request) {

	d := Data{
		Id:    10,
		Value: "Some data",
	}

	b, _ := json.Marshal(d)

	w.Write(b)
}
