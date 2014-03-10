package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"

	"runtime"
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

	runtime.GOMAXPROCS(config.getConcurrency())

	db, err = sql.Open("postgres", fmt.Sprintf("user=%v sslmode=disable", config.DBUser))

	if err != nil {
		log.Fatalf("Unable to connect to the database, reason -> %v\n", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(config.getMaxOpenConns())
	db.SetMaxIdleConns(config.getMaxIdleConns())

	http.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("./web/app/"))))

	r := httprouter.New()

	r.POST("/data/employee", RunInTransaction(EmployeeCreateHandler))
	r.POST("/data/employee/:id", RunInTransaction(EmployeeSaveHandler))
	r.DELETE("/data/employee/:id", RunInTransaction(EmployeeDeleteHandler))
	r.GET("/data/employee/:id", RunInTransaction(EmployeeHandler))
	r.GET("/data/employees", RunInTransaction(EmployeeListHandler))

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
