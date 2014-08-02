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

type LeftNavState struct {
	IsOverview  bool
	IsSchema    bool
	IsClients   bool
	IsEmployees bool
}

type NavEmployeeList struct {
	LeftNavState
	Employees []*Employee
}

type NavCustomerList struct {
	LeftNavState
	Customers []*Customer
}

func index(w http.ResponseWriter, r *http.Request, tx Transaction, vars map[string]string) error {

	e := &EmployeeList{}

	if err := e.Load(tx, 0, 40); err != nil {
		log.Warningf("[MAIN]: Unable to load data from database, error: '%v'", err)
		return err
	}

	lNav := LeftNavState{
		IsOverview: true,
	}

	nav := NavEmployeeList{
		LeftNavState: lNav,
		Employees:    e.Employees,
	}

	err := indexTmpl.Execute(w, nav)

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

	lNav := LeftNavState{
		IsClients: true,
	}

	nav := NavCustomerList{
		LeftNavState: lNav,
		Customers:    e.Customers,
	}

	err := clientListTmpl.Execute(w, nav)

	if err != nil {
		log.Warningf("[MAIN]: Unable to execute template for page 'clients', error: '%v'", err)
		return err
	}

	return nil
}
