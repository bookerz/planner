package main

import (
	log "github.com/golang/glog"
	"html/template"
	"net/http"
)

var indexTmpl *template.Template

var clientListTmpl *template.Template

func initWeb(config *Config) {
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
