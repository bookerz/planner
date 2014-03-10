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
	r.GET("/data/employee/:id", RunInTransaction(EmployeeHandler))
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func EmployeeHandler(w http.ResponseWriter, r *http.Request, tx *sql.Tx, vars map[string]string) error {

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		log.Printf("[EMPLOYEE]: Unable to format input parameter. error: '%v'", err)
		http.Error(w, "The id have to be a number", http.StatusBadRequest)
		return err
	}

	e := &Employee{
		Id: id,
	}

	if err = e.Load(tx); err != nil {
		switch {
		case err == sql.ErrNoRows:
			log.Printf("[EMPLOYEE]: Unable to to find employee with id '%v', error: '%v'", id, err)
			http.Error(w, "Employee not found", http.StatusNotFound)
		default:
			log.Printf("[EMPLOYEE]: Unable to load data from database, error: '%v'", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return err
	}

	b, err := json.Marshal(e)

	if err != nil {
		log.Printf("[EMPLOYEE]: Unable to marshal json data, error: '%v'", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)

	return nil
}

func DoRollback(tx *sql.Tx, prefix string) {
	if err := tx.Rollback(); err != nil {
		log.Printf(prefix+": unable to rollback transaction, error: '%v'", err)
	}
}
