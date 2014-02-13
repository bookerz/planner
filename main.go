package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "", "a configuration file is needed")
}

func main() {

	flag.Parse()

	if configFile == "" {
		log.Println("An example config file can look like this:")
		log.Printf("\n%v\n", ExampleConfig())
		log.Fatalln("You have to supply a config file '~/proj :> planner -config=/path/to/config.json'")
	}

	config, err := LoadConfig(configFile)

	if err != nil {
		log.Fatalf("No config available bailing out")
	}

	db, err := sql.Open("postgres", fmt.Sprintf("user=%v sslmode=disable", config.DBUser))

	if err != nil {
		log.Fatalf("Unable to connect to the database, reason -> %v\n", err)
	}

	defer db.Close()

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
