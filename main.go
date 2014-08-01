package main

import (
	"database/sql"
	"flag"
	"fmt"
	log "github.com/golang/glog"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"

	"runtime"
)

const (
	INFO  = 2
	WARN  = 1
	ERROR = 0
	FATAL = -1
)

var configFile string

var indexTmpl *template.Template

var clientListTmpl *template.Template

var db *sql.DB

func init() {
	flag.StringVar(&configFile, "config", "", "a configuration file is needed")
}

func main() {

	flag.Parse()

	if log.V(INFO) {
		log.V(INFO).Infoln("Flushing the log on every request")
	}

	if configFile == "" {
		log.Errorln("An example config file can look like this:")
		log.Errorf("\n%v\n", ExampleConfig())
		log.Fatalln("You have to supply a config file '~/proj :> planner -config=/path/to/config.json'")
	}

	config, err := LoadConfig(configFile)

	if err != nil {
		log.Fatalf("No config available bailing out")
	}

	runtime.GOMAXPROCS(config.getConcurrency())

	indexTmpl = template.Must(template.ParseFiles(
		config.getWebRoot()+"/templates/index.html",
		config.getWebRoot()+"/templates/navigation.html",
		config.getWebRoot()+"/templates/head.html",
		config.getWebRoot()+"/templates/leftmenu.html",
		config.getWebRoot()+"/templates/dashboard.html",
		config.getWebRoot()+"/templates/scripts.html"))

	clientListTmpl = template.Must(template.ParseFiles(
		config.getWebRoot()+"/templates/clients.html",
		config.getWebRoot()+"/templates/navigation.html",
		config.getWebRoot()+"/templates/head.html",
		config.getWebRoot()+"/templates/leftmenu.html",
		config.getWebRoot()+"/templates/clientlist.html",
		config.getWebRoot()+"/templates/scripts.html"))

	db, err = sql.Open("postgres", fmt.Sprintf("user=%v sslmode=disable", config.DBUser))

	if err != nil {
		log.Fatalf("Unable to connect to the database, reason -> %v\n", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(config.getMaxOpenConns())
	log.Infof("Setting max open db connections to %v\n", config.getMaxOpenConns())
	db.SetMaxIdleConns(config.getMaxIdleConns())
	log.Infof("Setting max idle db connections to %v\n", config.getMaxIdleConns())

	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(config.getWebRoot()+"/js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(config.getWebRoot()+"/css"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir(config.getWebRoot()+"/fonts"))))

	r := httprouter.New()

	r.POST("/data/employee", RunInTransaction(EmployeeCreateHandler))
	r.POST("/data/employee/:id", RunInTransaction(EmployeeSaveHandler))
	r.DELETE("/data/employee/:id", RunInTransaction(EmployeeDeleteHandler))
	r.GET("/data/employee/:id", RunInTransaction(EmployeeHandler))
	r.GET("/data/employees", RunInTransaction(EmployeeListHandler))

	r.POST("/data/customer", RunInTransaction(CustomerCreateHandler))
	r.POST("/data/customer/:id", RunInTransaction(CustomerSaveHandler))
	r.DELETE("/data/customer/:id", RunInTransaction(CustomerDeleteHandler))
	r.GET("/data/customer/:id", RunInTransaction(CustomerHandler))
	r.GET("/data/customers", RunInTransaction(CustomerListHandler))

	r.GET("/", RunInTransaction(index))
	r.GET("/clients", RunInTransaction(clients))

	http.Handle("/", r)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Can't start web server, err = %v\n", err.Error())
	}

	log.Flush()
}

func index(w http.ResponseWriter, r *http.Request, tx Transaction, vars map[string]string) error {

	e := &EmployeeList{}

	if err := e.Load(tx, 0, 40); err != nil {
		log.Warningf("[MAIN]: Unable to load data from database, error: '%v'", err)
		return err
	}

	err := indexTmpl.Execute(w, e.Employees)

	if err != nil {
		log.Warningf("[MAIN]: Unable to execute template for page 'index', error: '%v'", err)
		return err
	}

	return nil
}

func clients(w http.ResponseWriter, r *http.Request, tx Transaction, vars map[string]string) error {

	e := &CustomerList{}

	if err := e.Load(tx, 0, 40); err != nil {
		log.Warningf("[MAIN]: Unable to load data from database, error: '%v'", err)
		return err
	}

	err := clientListTmpl.Execute(w, e.Customers)

	if err != nil {
		log.Warningf("[MAIN]: Unable to execute template for page 'clients', error: '%v'", err)
		return err
	}

	return nil
}
