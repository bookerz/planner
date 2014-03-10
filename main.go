package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

var configFile string

var db *sql.DB

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

	rand.Seed(42)

	db, err = sql.Open("postgres", fmt.Sprintf("user=%v sslmode=disable", config.DBUser))

	if err != nil {
		log.Fatalf("Unable to connect to the database, reason -> %v\n", err)
	}

	defer db.Close()

	http.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("./web/app/"))))

	r := httprouter.New()
	r.GET("/data/employee/:id", EmployeeHandler)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func EmployeeHandler(w http.ResponseWriter, r *http.Request, vars map[string]string) {

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		log.Printf("[EMPLOYEE]: Unable to format input parameter. error: '%v'", err)
		http.Error(w, "The id have to be a number", http.StatusBadRequest)
		return
	}

	e := &Employee{
		Id: id,
	}

	tx, err := db.Begin()

	if err != nil {
		log.Printf("[EMPLOYEE]: Unable to create database transaction. error: '%v'", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = e.Load(tx)

	if err != nil {
		DoRollback(tx, "[EMPLOYEE]")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(e)

	if err != nil {
		DoRollback(tx, "[EMPLOYEE]")
		log.Printf("[EMPLOYEE]: Unable to marshal json data, error: '%v'", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[EMPLOYEE]: Unable to commit transaction. error: '%v'", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func DoRollback(tx *sql.Tx, prefix string) {
	if err := tx.Rollback(); err != nil {
		log.Printf(prefix+": unable to rollback transaction, error: '%v'", err)
	}
}
